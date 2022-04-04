package main

import (
	"fmt"
	"os"
	"smbserver/config"
	"smbserver/logger"
	"smbserver/model"
	"smbserver/router"
	"smbserver/serialmanager"
	"time"
)

func watchConnectivity() {

	enableInternet := true
	if len(config.MyIP) == 0 {
		enableInternet = false
	}

	for {
		dt := time.Now()
		for i, alarmTime := range model.Alarm {
			if alarmTime.Sub(dt).Seconds() < 0 {
				model.Status["alarm"] = 1
				go func() {
					time.Sleep(time.Second * 5)
					model.Status["alarm"] = 0
				}()
				model.Alarm = append(model.Alarm[:i], model.Alarm[i+1:]...)
			}
		}
		ip := config.GetIP()
		if len(ip) == 0 && enableInternet {
			model.Status["internetState"] = 0
			enableInternet = false
		} else if len(ip) != 0 {
			model.Status["internetState"] = 1
			model.Status["ip"] = ip
			enableInternet = true
		}
		model.Status["time"] = fmt.Sprintf("%2d.%2d", dt.Hour(), dt.Minute())
		time.Sleep(time.Second * 2)
	}
}

func main() {

	logger.Start()

	exitCh := make(chan interface{})
	go func() {
		time.Sleep(10 * time.Hour)
		exitCh <- true
	}()

	model.Status["alarm"] = 0
	model.Status["time"] = "00.00"
	serialmanager.AddDiscoverHandleFunc(func(e serialmanager.Event) {
		fmt.Println(e.Params())
	})

	serialmanager.AddRecvHandleFunc(func(e serialmanager.Event) {
		fmt.Println("param: ", e.Params())
		if model.Status["internetmodel.Status"] != model.Status["internetmodel.Status"] ||
			model.Status["time"] != e.Params()["t"] ||
			model.Status["alarm"] != e.Params()["alarm"] {
			var param = map[string]interface{}{
				"code": 200,
				"al":   model.Status["alarm"],
				"is":   model.Status["internetState"],
				"ti":   model.Status["time"],
				"ip":   model.Status["ip"],
			}
			fmt.Println(param)
			serialmanager.Write(param)
		}
	})
	go serialmanager.Run()

	// router.NewRouter().Run(config.Params["bind"].(string))

	go watchConnectivity()
	// time.Sleep(10 * time.Second)
	go router.NewRouter().Run(config.Params["bind"].(string))

	<-exitCh
	os.Exit(logger.Stop())
}

package main

import (
	"fmt"
	"github.com/HikoQiu/go-eureka-client/eureka"
	"log"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func main() {
	config := eureka.GetDefaultEurekaClientConfig()
	config.UseDnsForFetchingServiceUrls = false
	config.Region = "region-cn-hd-1"
	config.AvailabilityZones = map[string]string{
		"region-cn-hd-1": "zone-cn-hz-1",
	}
	config.ServiceUrl = map[string]string{
		"zone-cn-hz-1": "http://127.0.0.1:8761/eureka",
	}

	//custom logger
	eureka.SetLogger(func(level int, format string, a ...interface{}) {
		if level == eureka.LevelError {
			fmt.Println("[custom logger error] "+format, a)
		} else {
			fmt.Println("[custom logger] "+format, a)
		}
	})

	// run eureka client async
	client := eureka.DefaultClient.Config(config)
	client.
		Register("APP_ID_CLIENT_FROM_CONFIG", 8081).
		Run()

	api, err := client.Api()
	if err != nil {
		log.Fatalln("Failed to pick EurekaServerApi instance, err=", err.Error())
	}
	instances, err := api.QueryAllInstances()
	if err != nil {
		log.Fatalln("Failed to query all instances, err=", err.Error())
	}

	for i := range instances {
		instance := instances[i]
		log.Println(instance.Name)
		for i2 := range instance.Instances {
			actualInst := instance.Instances[i2]
			log.Println(actualInst.HealthCheckUrl)
		}
	}

	http.HandleFunc("/", handler)
	fmt.Println("Listening on port 8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}

package main

import (
	"fmt"
	"time"
)

func main() {

	for{
		apollo:=NewApollo()
		for _,serviceGroup:=range apollo.ServiceGroups{
			fmt.Println(serviceGroup.Name)
			for _,service:=range serviceGroup.Services {
				fmt.Println(service.Address,service.Online)
			}
		}
		time.Sleep(time.Second*5)
	}

}

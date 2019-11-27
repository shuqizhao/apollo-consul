package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	fmt.Println("Starting apollo-consul")
	ch := make(chan int)
	major := 0
	var curApollo *Apollo
	go func() {
		for {
			apolloEntity := NewApollo()
			time.Sleep(5 * time.Second)
			if apolloEntity.ConsulUrl != "" {
				os.Setenv("CONSUL_HTTP_ADDR", apolloEntity.ConsulUrl)
				Register(&apolloEntity)
			}
			newApolloEntity := &apolloEntity
			if apolloEntity.ConsulUrl != "" {
				newApolloEntity = Check(&apolloEntity)
			}
			if apolloEntity.MinorVersion != major || IsChange(curApollo, newApolloEntity) {
				major = apolloEntity.MinorVersion
				curApollo = newApolloEntity
				Build(newApolloEntity)
			}
		}
	}()
	<-ch
}

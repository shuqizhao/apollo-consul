package main

import (
	"fmt"
	"time"
	"os"
)

func main() {
	fmt.Println("Starting apollo-consul")
	ch := make(chan int)
	major := 0
	var curApollo *Apollo
	go func() {
		for {
			apolloEntity := NewApollo()
			Register(&apolloEntity)
			os.Setenv("CONSUL_HTTP_ADDR", apolloEntity.ConsulUrl)
			time.Sleep(5 * time.Second)
			newApolloEntity := Check(&apolloEntity)
			if apolloEntity.MinorVersion != major || IsChange(curApollo, newApolloEntity) {
				major = apolloEntity.MinorVersion
				curApollo = newApolloEntity
				Build(newApolloEntity)
			}
		}
	}()
	<-ch
}

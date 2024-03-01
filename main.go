package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("Starting apollo")
	ch := make(chan int)
	major := 0
	var curApollo *Apollo
	go func() {
		for {
			apolloEntity := NewApollo()

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

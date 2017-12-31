package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("Starting apollo-consul")

	ch := make(chan int)
	go func() {
		for  {
			Register()
			time.Sleep(5*time.Second)
			apollo:=Check()
			Build(apollo)
		}
	}()
	<-ch
}

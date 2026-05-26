package main

import (
	"fmt"
	"sync"
)

func process(ch chan int, x int){
	ch<- x
}

func receive(ch chan int, wg *sync.WaitGroup){
	defer wg.Done()
	msg:=<-ch
	fmt.Println("recieved", msg)
}

func main(){
	var wg sync.WaitGroup
	ch :=make(chan int, 5)

	for i:=0; i<5; i++ {
		go process(ch,i)
	}

	for i:=0; i<5;i++ {
		wg.Add(1)
		go receive(ch, &wg)
	}
	wg.Wait()
}
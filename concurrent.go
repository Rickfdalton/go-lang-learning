package main

import "fmt"

func sums(s []int, c chan int){
	sum:=0
	for _,v := range s {
		sum+=v
	}
	c <- sum
}

func main(){
	s:= []int{1,2,3,4,5,6,7,8,9}
	c:= make(chan int)
	go sums(s[:len(s)/2],c)
	go sums(s[len(s)/2 :],c)
	x,y := <-c, <-c

	fmt.Println(x,y,x+y)
}
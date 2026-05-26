package main

import (
	"fmt"
	"time"
	"strconv"
)

func worker(id int, jobs <-chan string, results chan<- string){
	// each worker process all jobs
	for job:=range jobs{
		fmt.Println("worker ", id,"processing ", job)
		time.Sleep(time.Second)
		results <- (job +" done!")
	}

}

func main(){
	jobs:=make(chan string,10)
	results:=make(chan string,10)

	for j:=1; j<5; j++ {
		jobs<- "job:"+strconv.Itoa(j)
	}

	for w:=1; w<4; w++ {
		go worker(w,jobs,results)
	}

	for i:=1; i<5; i++ {
		fmt.Println("result: ", <-results)
	}
}
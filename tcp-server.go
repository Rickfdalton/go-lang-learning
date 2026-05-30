package main

import (
	"os"
	"fmt"
	"net"
	"io"
	"bufio"
	"time"
)

const WORKERS int = 10

func main(){
	if len(os.Args) <2 {
		fmt.Println("Usage: go run tcp-server.go <port>")
		return
	}
	port := fmt.Sprintf(":%s", os.Args[1])

	listener, err := net.Listen("tcp", port)
	if err!= nil {
		fmt.Println("failed to create listener:", err)
		return
	}
	var ch = make(chan string, 100)
	var jobs= make(chan net.Conn, WORKERS)

	defer listener.Close()
	fmt.Printf("server listening on %s\n", listener.Addr())
	go logger(ch)
	for i:=0;i<WORKERS;i++{
		go worker(jobs,ch)
	}
	for{
		conn,err:= listener.Accept()
		if err!= nil {
			fmt.Println("failed to accept connection: ", err)
			continue
		}
		//handle connection
		jobs <- conn
	}

}

func handleConnection(ch chan string,conn net.Conn){
	defer conn.Close()
	reader := bufio.NewReader(conn)
	for{
		bytes, err:=reader.ReadBytes(byte('\n'))
		if err!= nil{
			if err!= io.EOF {
				fmt.Println("Error reading data: ", err)
				return
			}
			return
		}
		// non blocking logging.
		select{
		case ch<-fmt.Sprintf("[%s] %s", conn.RemoteAddr(), string(bytes)):
		case <-time.After(100 * time.Millisecond):
			fmt.Println("channel full")
		}
		line:= "From Backend: "+string(bytes)
		fmt.Printf("response:%s", line)
		
		_,err = conn.Write([]byte(line))
		if err!= nil{
			fmt.Println("Failed to write the data ",err )
			return
		}
	}

}

func logger(ch chan string){
	for msg:= range ch{
		time.Sleep(2 * time.Second) 
		fmt.Printf("[%s] REQUEST RECEIVED: %s", time.Now().Format("15:04:05"), msg)
	}
	return
}

func worker(jobs chan net.Conn, ch chan string){
	for conn:= range jobs{
		handleConnection(ch,conn)
	}
}
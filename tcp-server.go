package main

import (
	"os"
	"fmt"
	"net"
	"io"
	"bufio"
)

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

	defer listener.Close()
	fmt.Printf("server listening on %s\n", listener.Addr())

	for{
		conn,err:= listener.Accept()
		if err!= nil {
			fmt.Println("failed to accept connection: ", err)
			continue
		}
		//handle connection
		go handleConnection(conn)
	}

}


func handleConnection(conn net.Conn){
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
		fmt.Printf("request:%s", bytes)
		line:= "From Backend: "+fmt.Sprintf("Echo: %s",bytes)
		fmt.Printf("response:%s", line)
		
		_,err = conn.Write([]byte(line))
		if err!= nil{
			fmt.Println("Failed to write the data ",err )
			return
		}
	}

}
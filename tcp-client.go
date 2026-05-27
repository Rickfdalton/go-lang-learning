package main

import (
	"fmt"
	"net"
	"io"
	"bufio"	
)

func main(){

	conn, err := net.Dial("tcp", "localhost:9090")
	if err!= nil {
		fmt.Println("cannot dial server", err)
		return
	}
	defer conn.Close()
	_,err = conn.Write([]byte("Hi from client\n"))
	if err!= nil{
			fmt.Println("Failed to write the data ",err )
			return
	}
	reader := bufio.NewReader(conn)
	
	bytes, err:=reader.ReadBytes(byte('\n'))
	if err!= nil{
		if err!= io.EOF {
			fmt.Println("Error reading data: ", err)
			return
		}
		return
	}
	line:= "From Backend: "+string(bytes)
	fmt.Println(line)

}
package main

import (
	"fmt"
	"net"
	"io"
	"bufio"	
	"time"
)

func sendMessage(conn net.Conn,msg string ) (string, error) {	
	_,err := conn.Write([]byte(msg+"\n"))
	if err!= nil{
			fmt.Println("Failed to write the data ",err )
			return "",err

	}
	reader := bufio.NewReader(conn)
	
	bytes, err:=reader.ReadBytes(byte('\n'))
	if err!= nil{
		if err!= io.EOF {
			fmt.Println("Error reading data: ", err)
		}
			return "",err
	}
	line:= "From Backend: "+string(bytes)
	return line, nil
}

func main(){

	conn, err := net.Dial("tcp", "localhost:9090")
	if err!= nil {
		fmt.Println("cannot dial server", err)
			return 
	}
	defer conn.Close()

	for {
		fmt.Print("Enter the message to send: ") 
		var msg string
		_,err := fmt.Scanf("%s", &msg)
		if err!=nil{
			fmt.Print("sorry cannot read your message")
			continue
		}
		response, err :=sendMessage(conn,msg)
		if err!= nil {
			fmt.Println("Connection expired")
		}else{
			fmt.Println(response)
		}
		conn.SetDeadline(time.Now().Add(5 * time.Second))
	
	}
	
}
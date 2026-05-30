package main

import (
	"fmt"
	"net"
	"io"
	"bufio"	
	"time"
)


func dialServer() (net.Conn, error){
	conn, err := net.Dial("tcp", "localhost:9090")
	if err!= nil {
		fmt.Println("cannot dial server", err)
			return nil, err
	}
	return conn, nil
}

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
	var conn net.Conn
	var err error

	conn, err = dialServer()
	if err!=nil{
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
			var err2,err3 error
			conn.Close()
			fmt.Println("Reconnecting..")
			conn, err2 =  dialServer()
			if err2!= nil {
				return
			}
			response, err3 =sendMessage(conn,msg)
			if err3!=nil{
				return
			}						
		}
		fmt.Println(response)
		conn.SetDeadline(time.Now().Add(5 * time.Second))
	
	}
	
}
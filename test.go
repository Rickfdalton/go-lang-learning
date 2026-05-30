package main

import (
    "fmt"
    "net"
    "sync"
	"bufio"
	"io"
)

func main() {
    var wg sync.WaitGroup
    
    for i := 0; i < 15; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            conn, err := net.Dial("tcp", "localhost:9090")
			if err!= nil {
				fmt.Println("cannot dial server", err)
					return
			}
			response, err :=sendMessage(conn,"message"+fmt.Sprintf("%d", id) )
			if err!= nil {
				fmt.Println("error")
				return
			}
			fmt.Println("Response received:",response)
			conn.Close()

        }(i)
    }
    
    wg.Wait()
    fmt.Println("all done")
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
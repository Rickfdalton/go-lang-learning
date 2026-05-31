package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

const WORKERS int = 10

var shuttingDown bool

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run tcp-server.go <port>")
		return
	}
	port := fmt.Sprintf(":%s", os.Args[1])

	listener, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Println("failed to create listener:", err)
		return
	}
	var ch = make(chan string, 100)
	var jobs = make(chan net.Conn, WORKERS)
	var sigCh = make(chan os.Signal, 1)
	var quit = make(chan struct{})

	var wg sync.WaitGroup

	defer listener.Close()
	fmt.Printf("server listening on %s\n", listener.Addr())
	go logger(ch)
	for i := 0; i < WORKERS; i++ {
		wg.Add(1)
		go worker(jobs, ch, quit, &wg)
	}

	go checkInterrupt(sigCh, &wg, jobs, quit, listener)
	for {
		conn, err := listener.Accept()
		if err != nil {
			if shuttingDown {
				return
			}
			fmt.Println("failed to accept connection: ", err)
			continue
		}
		//handle connection
		fmt.Println("sending conn to jobs...")
		jobs <- conn
		fmt.Println("conn sent!")
	}

}

func handleConnection(ch chan string, quit chan struct{}, conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	for {
		bytes, err := reader.ReadBytes(byte('\n'))
		fmt.Println("ReadBytes returned, err:", err, "bytes:", string(bytes))

		if err != nil {
			select {
			case <-quit:
				fmt.Println("shutdown, closing connection")
				return
			default:
				if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
					continue // just a timeout
				}
				if err != io.EOF {
					fmt.Println("Error reading:", err)
				}
				return
			}
		}
		// non blocking logging.
		select {
		case ch <- fmt.Sprintf("[%s] %s", conn.RemoteAddr(), string(bytes)):
		case <-time.After(100 * time.Millisecond):
			fmt.Println("channel full")
		}
		line := "From Backend: " + string(bytes)
		fmt.Printf("response:%s", line)
		time.Sleep(5 * time.Second)

		_, err = conn.Write([]byte(line))
		if err != nil {
			fmt.Println("Failed to write the data ", err)
			return
		}

		if shuttingDown { 
			return
		}
	}
	fmt.Println("handle conneciton done")

}

func logger(ch chan string) {
	for msg := range ch {
		fmt.Printf("[%s] REQUEST RECEIVED: %s", time.Now().Format("15:04:05"), msg)
	}
	return
}

func worker(jobs chan net.Conn, ch chan string, quit chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()
	for conn := range jobs {
		handleConnection(ch, quit, conn)
	}
	fmt.Println("worker finished")
}

func checkInterrupt(sigCh chan os.Signal, wg *sync.WaitGroup, jobs chan net.Conn, quit chan struct{}, listener net.Listener) {
	signal.Notify(sigCh, syscall.SIGINT)
	<-sigCh
	fmt.Println("Ungraceful shutdown")
	shuttingDown = true
	listener.Close()
	close(jobs)
	wg.Wait()
	close(quit)
	fmt.Println("all workers done, bye!")
	os.Exit(0)
}

package main

import "fmt"

type ConnectionTracker struct{
	MaxConns   int
	connections []string
}

func (c *ConnectionTracker) Add(ip string) bool{
	if c.IsFull() {
		fmt.Println("Unable to add",ip, ": Max Connections Reached!")
		return false
	}
	c.connections = append(c.connections, ip)
	fmt.Println("Successfully added connection ", ip)
	return true
}

func (c *ConnectionTracker) Remove(ip string){
	for i,v := range c.connections {
		if v == ip {
			first_set := c.connections[:i]
			last_set := c.connections[i+1:]
			c.connections = append(first_set, last_set...)
			fmt.Println("Successfully removed connection ", ip)
			return
		}
	}
	fmt.Println ("Connection", ip, "not exists")
}

func (c *ConnectionTracker) Count() int {
	return len(c.connections)
}

func (c *ConnectionTracker) IsFull() bool {
	return c.Count() >= c.MaxConns 
}

func main() {
	tracker := ConnectionTracker{MaxConns: 3}
	tracker.Add("192.168.1.1")
	tracker.Add("192.168.1.2")
	tracker.Add("192.168.1.3")
	tracker.Add("192.168.1.4") // this should fail
	tracker.Remove("192.168.1.4")
	fmt.Println("current connections:", tracker.Count())
	fmt.Println("is maximum connections reached? ", tracker.IsFull())
}
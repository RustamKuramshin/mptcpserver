package main

import (
	"bufio"
	"flag"
	"fmt"
	"mptcpserver/lib"
	"net"
	"os"
	"time"
)

type MpTCPServer struct {
	Addr string
	Port string
}

var ttlmap = lib.NewTTLMap(1, 10, 1)

func (server *MpTCPServer) ListenAndServe(tcpaddr, tcpport string) error {

	server.Addr = tcpaddr
	server.Port = tcpport

	fmt.Printf("MpTCPServer :: Start on host %v and port %v \n", server.Addr, server.Port)
	listener, err := net.Listen("tcp", server.Addr+":"+server.Port)
	if err != nil {
		return err
	}

	defer listener.Close()

	for {

		tcpconn, err := listener.Accept()
		if err != nil {
			fmt.Printf("MpTCPServer :: Error connection :: %v \n", err)
			continue
		}

		go HandleFunc(tcpconn)
	}
}

func HandleFunc(tcpconn net.Conn) error {

	tcpReader := bufio.NewReader(tcpconn)
	tcpScan := bufio.NewScanner(tcpReader)

	defer tcpconn.Close()

	for {
		scanned := tcpScan.Scan()
		if !scanned {
			if err := tcpScan.Err(); err != nil {
				fmt.Printf("MpTCPServer :: Address:: %v :: "+
					"Error reading from socket :: %v \n", tcpconn.RemoteAddr(), err)
				return err
			}
			break
		}

		err := tcpconn.SetDeadline(time.Now().Add(time.Duration(5)))
		if err != nil {
			return err
		}

		message, err := lib.GetMpMessage(tcpScan.Bytes())
		if err != nil {
			return err
		}

		ttlmap.Put(message.Domain, message.Ip)
	}
	return nil
}

func main() {

	var addr, port string

	if len(os.Args) > 3 {
		panic("Usage go run tcp-server.go [-tcpaddr=<host>] [-tcpport=<port>]")
	}

	flag.StringVar(&addr, "tcpaddr", "0.0.0.0", "Set host address")
	flag.StringVar(&port, "tcpport", "6000", "Set tcp port")
	flag.Parse()

	server := new(MpTCPServer)
	err := server.ListenAndServe(addr, port)
	if err != nil {
		fmt.Println(err)
	}
}

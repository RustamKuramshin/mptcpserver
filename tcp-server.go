package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
)

type MpTCPServer struct {
	Addr string
	Port string
}

var ttlmap = NewTTLMap(1, 10)

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

		fmt.Printf("MpTCPServer :: Accepted connection from %v \n", tcpconn.RemoteAddr())
		HandleFunc(tcpconn)
	}
}

func HandleFunc(tcpconn net.Conn) error {

	tcpReader := bufio.NewReader(tcpconn)
	tcpScan := bufio.NewScanner(tcpReader)

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

		message, err := NewMpMessage(tcpScan.Bytes())
		if err != nil {
			return err
		}

		ttlmap.Put(message.Domain, message.Ip)

		fmt.Printf("MpTCPServer :: Accepted from %v :: %v \n", tcpconn.RemoteAddr(), message)

		//sentData := receivedData
		//tcpWriter.WriteByte(sentData)
		//tcpWriter.Flush()
		//fmt.Printf("MpTCPServer :: Sent to %v :: %v \n", tcpconn.RemoteAddr(), sentData)

	}
	return nil
}

func main() {

	var addr, port string

	if len(os.Args) > 3 {
		panic("Usage go run tcp-server.go [-tcpaddr=<host>] [-tcpport=<port>]")
	}

	flag.StringVar(&addr, "tcpaddr", "0.0.0.0", "Set host address")
	flag.StringVar(&port, "tcpport", "8080", "Set tcp port")
	flag.Parse()

	server := new(MpTCPServer)
	server.ListenAndServe(addr, port)
}

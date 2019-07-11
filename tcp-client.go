package main

import (
	"flag"
	"fmt"
	"math/rand"
	"mptcpserver/lib"
	"net"
	"os"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

func StringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func getRandomString(length int) string {
	return StringWithCharset(length, charset)
}

func main() {

	var addr, port string

	if len(os.Args) > 3 {
		panic("Usage go run tcp-client.go [-tcpaddr=<host>] [-tcpport=<port>]")
	}

	flag.StringVar(&addr, "tcpaddr", "127.0.0.1", "Set host address")
	flag.StringVar(&port, "tcpport", "6000", "Set tcp port")
	flag.Parse()

	for range time.Tick(time.Duration(1) * time.Second) {
		go func() {

			mpmessage := new(lib.MpMessage)
			mpmessage.Domain = getRandomString(8)
			mpmessage.Ip = rand.Uint32()

			conn, err := net.Dial("tcp", addr+":"+port)
			if err != nil {
				fmt.Println(err)
				return
			}
			defer conn.Close()

			fmt.Printf("Send to server :: %v \n", mpmessage)

			if n, err := conn.Write(mpmessage.Serialize()); n == 0 || err != nil {
				fmt.Println(err)
				return
			}
		}()
	}
}

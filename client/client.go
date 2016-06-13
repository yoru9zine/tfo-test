package main

import (
	"flag"
	"log"
	"net"
	"syscall"
)

var (
	tfo = flag.Bool("tfo", true, "tcp fast open")
)

func main() {
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	if err != nil {
		log.Fatalf("failed to create socket: %s", err)
	}

	ip := net.ParseIP("127.0.0.1")
	var addr [4]byte
	copy(addr[:], ip[12:16])
	sa := &syscall.SockaddrInet4{Addr: addr, Port: 3333}

	opts := 0
	if *tfo {
		opts = syscall.MSG_FASTOPEN
	}

	err = syscall.Sendto(fd, []byte("hi from client"), opts, sa)
	if err != nil {
		if err == syscall.EOPNOTSUPP {
			log.Fatalf("not supported")
		}
		log.Fatalf("unknown error: %s", err)
	}
	buf := make([]byte, 1024)
	n, err := syscall.Read(fd, buf)
	if err != nil {
		log.Fatalf("failed to read: %s", err)
	}
	log.Println(string(buf[:n]))
}

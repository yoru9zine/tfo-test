package main

import (
	"flag"
	"log"
	"net"
	"syscall"
)

const TCP_FASTOPEN int = 23

var (
	tfo = flag.Bool("tfo", true, "tcp fast open")
)

func main() {
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	if err != nil {
		log.Fatalf("failed to create socket: %s", err)
	}

	opts := syscall.SO_REUSEADDR
	if *tfo {
		opts |= TCP_FASTOPEN
	}
	if err = syscall.SetsockoptInt(fd, syscall.SOL_TCP, opts, 1); err != nil {
		log.Fatalf("failed to setsockopt: %s", err)
	}

	ip := net.ParseIP("127.0.0.1")
	var addr [4]byte
	copy(addr[:], ip[12:16])
	sa := &syscall.SockaddrInet4{Addr: addr, Port: 3333}

	if err = syscall.Bind(fd, sa); err != nil {
		log.Fatalf("failed to bind: %s", err)
	}

	if err = syscall.Listen(fd, 32); err != nil {
		log.Fatalf("failed to listen: %s", err)
	}

	for {
		cfd, _, err := syscall.Accept(fd)
		if err != nil {
			log.Fatalf("failed to accept: %s", err)
		}

		if _, err = syscall.Write(cfd, []byte("hi from server")); err != nil {
			if err == syscall.EOPNOTSUPP {
				log.Fatalf("not supported")
			}
			log.Fatalf("unknown error: %s", err)
		}

		buf := make([]byte, 1024)
		n, err := syscall.Read(cfd, buf)
		if err != nil {
			log.Fatalf("failed to read: %s", err)
		}
		log.Println(string(buf[:n]))
		if err := syscall.Close(cfd); err != nil {
			log.Fatalf("failed to close: %s", err)
		}
	}
}

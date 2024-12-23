package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

/*
=== Утилита telnet ===

Реализовать примитивный telnet клиент:
Примеры вызовов:
go-telnet --timeout=10s host port go-telnet mysite.ru 8080 go-telnet --timeout=3s 1.1.1.1 123

Программа должна подключаться к указанному хосту (ip или доменное имя) и порту по протоколу TCP.
После подключения STDIN программы должен записываться в сокет, а данные полученные и сокета должны выводиться в STDOUT
Опционально в программу можно передать таймаут на подключение к серверу (через аргумент --timeout, по умолчанию 10s).

При нажатии Ctrl+D программа должна закрывать сокет и завершаться. Если сокет закрывается со стороны сервера, программа должна также завершаться.
При подключении к несуществующему сервер, программа должна завершаться через timeout.
*/
type Config struct {
	address string
	timeout time.Duration
}

func ReadConfig() (Config, error) {
	cfg := Config{}
	timeout := flag.Duration("timeout", 10*time.Second, "connection timeout")
	flag.Parse()

	cfg.timeout = *timeout

	if flag.NArg() < 2 {
		return cfg, fmt.Errorf("Usage: go-telnet [--timeout=10s] host port")
	}

	host := flag.Arg(0)
	port := flag.Arg(1)

	cfg.address = net.JoinHostPort(host, port)

	return cfg, nil
}

func ConnectTCP(ctx context.Context, cfg Config) (net.Conn, error) {
	dialer := net.Dialer{}
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(1 * time.Second):
			conn, err := dialer.DialContext(ctx, "tcp", cfg.address)
			if err == nil {
				return conn, nil
			}
			fmt.Printf("Connection error:%s try to reconnect\n", err)
		}
	}
}
func main() {

	cfg, err := ReadConfig()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), cfg.timeout)

	conn, err := ConnectTCP(ctx, cfg)

	if err != nil {
		fmt.Printf("Failed to connect to %s: %v\n", cfg.address, err)
		return
	}
	defer conn.Close()

	fmt.Printf("Connected to %s\n", cfg.address)

	go func() {
		defer cancel()
		io.Copy(os.Stdout, conn)
	}()

	go func() {
		reader := bufio.NewReader(os.Stdin)
		defer cancel()
		for {
			input, err := reader.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					fmt.Println("Detected EOF, closing connection...")
					conn.Close()
					return
				}
				fmt.Printf("Error reading from stdin: %v\n", err)
				return
			}
			_, err = conn.Write([]byte(input))
			if err != nil {
				fmt.Printf("Error writing to connection: %v\n", err)
				return
			}
		}
	}()

	<-ctx.Done()
	fmt.Printf("\nConnection closed: %s", ctx.Err())
}

package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
)

type NetCat struct {
	protocol    string
	host        string
	port        string
	isServer    bool
	connections sync.WaitGroup
	listener    net.Listener
	packetConn  net.PacketConn
	ctx         context.Context
	cancel      context.CancelFunc
}

func NewNetCat() *NetCat {
	return &NetCat{}
}

func (nc *NetCat) parseFlags() error {
	tcp := flag.Bool("t", true, "Use TCP (default)")
	udp := flag.Bool("u", false, "Use UDP")
	port := flag.String("p", "0", "Local port")
	listen := flag.Bool("l", false, "Listen mode")
	flag.Parse()

	if flag.NArg() < 1 && !*listen {
		return fmt.Errorf("Usage: nc [-t|-u] [-l] [-p port] [hostname] [port]")
	}

	if *udp {
		nc.protocol = "udp"
	} else if *tcp {
		nc.protocol = "tcp"
	}

	nc.isServer = *listen
	nc.port = *port

	if !nc.isServer {
		if flag.NArg() < 2 {
			return fmt.Errorf("Missing hostname or port for client mode")
		}
		nc.host = flag.Arg(0)
		nc.port = flag.Arg(1)
	}

	return nil
}

func (nc *NetCat) setupSignalHandling() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		fmt.Println("\nReceived shutdown signal")
		nc.shutdown()
	}()
}

func (nc *NetCat) shutdown() {
	nc.cancel()

	if nc.listener != nil {
		nc.listener.Close()
	}
	if nc.packetConn != nil {
		nc.packetConn.Close()
	}

	// Ждем завершения всех соединений
	done := make(chan struct{})
	go func() {
		nc.connections.Wait()
		close(done)
	}()

	select {
	case <-done:
		fmt.Println("All connections closed gracefully")
	case <-time.After(5 * time.Second):
		fmt.Println("Shutdown timed out")
	}
}

func (nc *NetCat) run() error {
	if err := nc.parseFlags(); err != nil {
		return err
	}

	nc.ctx, nc.cancel = context.WithCancel(context.Background())
	nc.setupSignalHandling()

	if nc.isServer {
		return nc.runServer()
	}
	return nc.runClient()
}

func (nc *NetCat) runServer() error {
	addr := fmt.Sprintf(":%s", nc.port)

	if nc.protocol == "tcp" {
		listener, err := net.Listen(nc.protocol, addr)
		if err != nil {
			return fmt.Errorf("Error listening: %v", err)
		}
		nc.listener = listener
		defer nc.listener.Close()

		fmt.Printf("Listening on %s...\n", listener.Addr())

		return nc.handleTCPServer()
	}

	// UDP Server
	conn, err := net.ListenPacket(nc.protocol, addr)
	if err != nil {
		return fmt.Errorf("Error listening: %v", err)
	}
	nc.packetConn = conn
	defer nc.packetConn.Close()

	fmt.Printf("Listening on %s...\n", conn.LocalAddr())
	return nc.handleUDPServer()
}

func (nc *NetCat) handleTCPServer() error {
	for {
		select {
		case <-nc.ctx.Done():
			return nil
		default:
			conn, err := nc.listener.Accept()
			if err != nil {
				if !isClosedError(err) {
					fmt.Printf("Error accepting: %v\n", err)
				}
				continue
			}

			nc.connections.Add(1)
			go nc.handleTCPConnection(conn)
		}
	}
}

func (nc *NetCat) handleTCPConnection(conn net.Conn) {
	defer nc.connections.Done()
	defer conn.Close()

	fmt.Printf("New connection from %s\n", conn.RemoteAddr())

	ctx, cancel := context.WithCancel(nc.ctx)
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(2)

	// Чтение из соединения
	go func() {
		defer wg.Done()
		io.Copy(os.Stdout, conn)
	}()

	// Запись в соединение
	go func() {
		defer wg.Done()
		select {
		case <-ctx.Done():
			return
		default:
			io.Copy(conn, os.Stdin)
		}
	}()

	wg.Wait()
}

func (nc *NetCat) handleUDPServer() error {
	buffer := make([]byte, 1024)

	for {
		select {
		case <-nc.ctx.Done():
			return nil
		default:
			nc.packetConn.SetReadDeadline(time.Now().Add(time.Second))
			n, remoteAddr, err := nc.packetConn.ReadFrom(buffer)
			if err != nil {
				if !isTimeoutError(err) && !isClosedError(err) {
					fmt.Printf("Error reading UDP: %v\n", err)
				}
				continue
			}

			fmt.Printf("Received %d bytes from %s\n", n, remoteAddr)
			os.Stdout.Write(buffer[:n])

			nc.connections.Add(1)
			go func() {
				defer nc.connections.Done()
				_, err := nc.packetConn.WriteTo(buffer[:n], remoteAddr)
				if err != nil && !isClosedError(err) {
					fmt.Printf("Error writing UDP: %v\n", err)
				}
			}()
		}
	}
}

func (nc *NetCat) runClient() error {
	addr := fmt.Sprintf("%s:%s", nc.host, nc.port)

	if nc.protocol == "tcp" {
		return nc.handleTCPClient(addr)
	}
	return nc.handleUDPClient(addr)
}

func (nc *NetCat) handleTCPClient(addr string) error {
	conn, err := net.Dial(nc.protocol, addr)
	if err != nil {
		return fmt.Errorf("Error connecting: %v", err)
	}
	defer conn.Close()

	ctx, cancel := context.WithCancel(nc.ctx)
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(2)

	// Чтение из соединения
	go func() {
		defer wg.Done()
		io.Copy(os.Stdout, conn)
	}()

	// Запись в соединение
	go func() {
		defer wg.Done()
		select {
		case <-ctx.Done():
			return
		default:
			io.Copy(conn, os.Stdin)
		}
	}()

	wg.Wait()
	return nil
}

func (nc *NetCat) handleUDPClient(addr string) error {
	conn, err := net.Dial(nc.protocol, addr)
	if err != nil {
		return fmt.Errorf("Error connecting: %v", err)
	}
	defer conn.Close()

	ctx, cancel := context.WithCancel(nc.ctx)
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(2)

	// Чтение из соединения
	go func() {
		defer wg.Done()
		buffer := make([]byte, 1024)
		for {
			select {
			case <-ctx.Done():
				return
			default:
				n, err := conn.Read(buffer)
				if err != nil {
					if !isClosedError(err) {
						fmt.Printf("Error reading UDP: %v\n", err)
					}
					return
				}
				os.Stdout.Write(buffer[:n])
			}
		}
	}()

	// Запись в соединение
	go func() {
		defer wg.Done()
		buffer := make([]byte, 1024)
		for {
			select {
			case <-ctx.Done():
				return
			default:
				n, err := os.Stdin.Read(buffer)
				if err != nil {
					if err != io.EOF {
						fmt.Printf("Error reading stdin: %v\n", err)
					}
					return
				}
				_, err = conn.Write(buffer[:n])
				if err != nil {
					fmt.Printf("Error writing UDP: %v\n", err)
					return
				}
			}
		}
	}()

	wg.Wait()
	return nil
}

func isClosedError(err error) bool {
	if err == nil {
		return false
	}
	return err == io.EOF || err == io.ErrClosedPipe ||
		strings.Contains(err.Error(), "use of closed network connection")
}

func isTimeoutError(err error) bool {
	if err == nil {
		return false
	}
	netErr, ok := err.(net.Error)
	return ok && netErr.Timeout()
}

func main() {
	nc := NewNetCat()
	if err := nc.run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

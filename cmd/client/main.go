package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"osi/internal/layers"
	"osi/internal/models"
	"time"
)

const (
	ServerPort = 9000
	Protocol   = "CHAT"
	MaxRetries = 3
	AckTimeout = 2 * time.Second
	HandshakeTimeout = 5 * time.Second
	BufferSize = 1024

	AsusIP  = "192.168.1.101"
	AsusMAC = "aa:bb:cc:dd:ee:01"

	MacIP  = "192.168.1.100"
	MacMAC = "aa:bb:cc:dd:ee:02"
)

var (
	sender = models.Machine{
		IPAddress:  MacIP,
		MACAddress: MacMAC,
		Port:       0,
	}

	receiver = models.Machine{
		IPAddress:  MacIP,
		MACAddress: MacMAC,
		Port:       ServerPort,
	}
)

func main() {
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d",
		receiver.IPAddress, receiver.Port))
	if err != nil {
		log.Fatal("ResolveUDPAddr:", err)
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Fatal("DialUDP:", err)
	}
	defer conn.Close()

	establishHandshake(conn)

	sessionID := layers.NewSessionID()
	fmt.Println("Session ID:", sessionID)
	fmt.Println("Type a message and press Enter ('exit' to quit):")

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		msg := scanner.Text()

		if msg == "" {
			continue
		}

		if msg == "exit" {
			teardown(conn)
			return
		}

		sendMessage(conn, []byte(msg), sessionID)
	}
}

func sendMessage(conn *net.UDPConn, msg []byte, sessionID string) {
	data := layers.ApplicationEncode(msg, Protocol)
	data = layers.PresentationEncode(data)
	data = layers.SessionEncode(data, sessionID)

	chunks := layers.TransportChunk(data, sender, receiver)
	fmt.Printf("Split into %d chunk(s)\n", len(chunks))

	ackBuffer := make([]byte, BufferSize)

	for i, chunk := range chunks {
		packet := layers.NetworkEncode(chunk, sender, receiver)
		packet = layers.DatalinkEncode(packet, sender, receiver)

		acked := false

		for attempt := range MaxRetries {
			if _, err := conn.Write(packet); err != nil {
				fmt.Println("Write error:", err)
				break
			}

			conn.SetReadDeadline(time.Now().Add(AckTimeout))
			n, _, err := conn.ReadFromUDP(ackBuffer)
			if err == nil {
				fmt.Printf("Chunk %d/%d — %s\n", i+1, len(chunks), string(ackBuffer[:n]))
				acked = true
				break
			}
			fmt.Printf("Chunk %d — No ACK, retrying... (attempt %d/%d)\n", i+1, attempt+1, MaxRetries)
		}

		if !acked {
			fmt.Printf("Failed to send chunk %d\n", i+1)
		}
	}
}

func establishHandshake(conn *net.UDPConn) {
	fmt.Println("Initiating handshake...")

	if _, err := conn.Write([]byte("SYN")); err != nil {
		log.Fatal("Handshake: failed to send SYN:", err)
	}
	fmt.Println("Sent SYN")

	buffer := make([]byte, BufferSize)
	conn.SetDeadline(time.Now().Add(HandshakeTimeout))

	n, _, err := conn.ReadFromUDP(buffer)
	if err != nil {
		log.Fatal("Handshake: no response:", err)
	}
	if string(buffer[:n]) != "SYN-ACK" {
		log.Fatal("Handshake: expected SYN-ACK, got:", string(buffer[:n]))
	}
	fmt.Println("Got SYN-ACK")

	if _, err := conn.Write([]byte("ACK")); err != nil {
		log.Fatal("Handshake: failed to send ACK:", err)
	}
	conn.SetDeadline(time.Time{})
	fmt.Println("Connection ESTABLISHED")
}

func teardown(conn *net.UDPConn) {
	buffer := make([]byte, BufferSize)

	if _, err := conn.Write([]byte("FIN")); err != nil {
		log.Fatal("Teardown: failed to send FIN:", err)
	}
	fmt.Println("Sent FIN")

	conn.SetReadDeadline(time.Now().Add(HandshakeTimeout))

	n, _, err := conn.ReadFromUDP(buffer)
	if err != nil || string(buffer[:n]) != "FIN-ACK" {
		log.Fatal("Teardown: expected FIN-ACK")
	}
	fmt.Println("Got FIN-ACK")

	n, _, err = conn.ReadFromUDP(buffer)
	if err != nil || string(buffer[:n]) != "FIN" {
		log.Fatal("Teardown: expected FIN")
	}
	fmt.Println("Got FIN from server")

	conn.Write([]byte("ACK"))
	fmt.Println("Sent ACK — Connection closed")
}

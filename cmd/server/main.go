package main

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"net"
	"osi/internal/layers"
)

const (
	ServerPort = ":9000"
	BufferSize = 1024

	// Toggle these to simulate network issues
	EnablePacketDrop    = false
	PacketDropPercent   = 30

	EnableCorruption    = false
	CorruptionPercent   = 20
	CorruptedBytesCount = 5
)

func main() {
	conn := startServer()
	defer conn.Close()

	establishHandshake(conn)

	buffer := make([]byte, BufferSize)
	chunkStore := make(map[string][][]byte)

	for {
		n, clientAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("ReadFromUDP:", err)
			continue
		}

		data := make([]byte, n)
		copy(data, buffer[:n])

		if string(data) == "FIN" {
			teardown(conn, clientAddr)
			return
		}

		if shouldDrop() {
			fmt.Println("DROPPED packet from", clientAddr)
			continue
		}

		if shouldCorrupt() {
			corrupt(data)
			fmt.Println("CORRUPTED packet from", clientAddr)
		}

		fmt.Println("\n--- Received packet ---")

		data = layers.DatalinkDecode(data)
		data = layers.NetworkDecode(data)

		if _, err = conn.WriteToUDP([]byte("ACK"), clientAddr); err != nil {
			fmt.Println("WriteToUDP:", err)
		}

		var total int
		parts := bytes.SplitN(data, []byte("|"), 5)
		fmt.Sscanf(string(parts[3]), "TOTAL:%d", &total)

		key := clientAddr.String()
		chunkStore[key] = append(chunkStore[key], data)

		if len(chunkStore[key]) < total {
			fmt.Printf("Chunk received (%d/%d), waiting for more...\n",
				len(chunkStore[key]), total)
			continue
		}

		fmt.Println("\n--- All chunks received, reassembling ---")
		reassembled := layers.TransportDechunk(chunkStore[key])
		delete(chunkStore, key)

		reassembled = layers.SessionDecode(reassembled)
		reassembled = layers.PresentationDecode(reassembled)
		reassembled = layers.ApplicationDecode(reassembled)
		fmt.Println("Original message:", string(reassembled))
	}
}

func startServer() *net.UDPConn {
	addr, err := net.ResolveUDPAddr("udp", ServerPort)
	if err != nil {
		log.Fatal("ResolveUDPAddr:", err)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Fatal("ListenUDP:", err)
	}

	fmt.Println("UDP server running on", ServerPort)
	return conn
}

func establishHandshake(conn *net.UDPConn) {
	fmt.Println("Waiting for connection...")
	buffer := make([]byte, BufferSize)

	n, clientAddr, err := conn.ReadFromUDP(buffer)
	if err != nil {
		log.Fatal("Handshake: ReadFromUDP:", err)
	}
	if string(buffer[:n]) != "SYN" {
		log.Fatal("Handshake: expected SYN, got:", string(buffer[:n]))
	}
	fmt.Println("Got SYN from", clientAddr)

	if _, err = conn.WriteToUDP([]byte("SYN-ACK"), clientAddr); err != nil {
		log.Fatal("Handshake: failed to send SYN-ACK:", err)
	}
	fmt.Println("Sent SYN-ACK")

	n, _, err = conn.ReadFromUDP(buffer)
	if err != nil {
		log.Fatal("Handshake: ReadFromUDP:", err)
	}
	if string(buffer[:n]) != "ACK" {
		log.Fatal("Handshake: expected ACK, got:", string(buffer[:n]))
	}
	fmt.Println("Connection ESTABLISHED")
}

func teardown(conn *net.UDPConn, clientAddr *net.UDPAddr) {
	buffer := make([]byte, BufferSize)
	conn.WriteToUDP([]byte("FIN-ACK"), clientAddr)
	fmt.Println("Got FIN, sent FIN-ACK")

	conn.WriteToUDP([]byte("FIN"), clientAddr)
	fmt.Println("Sent FIN")

	n, _, err := conn.ReadFromUDP(buffer)
	if err != nil || string(buffer[:n]) != "ACK" {
		log.Fatal("Teardown: expected ACK")
	}
	fmt.Println("Got ACK — Connection closed")
}

func shouldDrop() bool {
	return EnablePacketDrop && rand.Intn(100) < PacketDropPercent
}

func shouldCorrupt() bool {
	return EnableCorruption && rand.Intn(100) < CorruptionPercent
}

func corrupt(data []byte) {
	for range CorruptedBytesCount {
		data[rand.Intn(len(data))] = byte(rand.Intn(256))
	}
}

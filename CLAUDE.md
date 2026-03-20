# OSI Networking Project

## Goal
Learn networking fundamentals by building things from scratch in Go — not reading textbooks. Started from a Rust networking repo (tcp-over-udp, dns-resolver, gossip-membership) and decided to learn the underlying concepts hands-on before attempting those projects.

## Learning Path
1. ~~OSI model visualizer~~ — done
2. ~~Build TCP features on top of UDP~~ — done
3. DNS resolver
4. SWIM gossip protocol

## What's Been Built
- Two real UDP servers (client + server) that communicate over the local network
- 7 OSI layer functions (encode + decode) that wrap/unwrap data as it moves down/up the stack
- Tested across two real machines on the same WiFi (Mac @ 192.168.1.100 + Asus laptop @ 192.168.1.101)
- TCP features built on top of UDP:
  - **Checksums** (CRC32) in Data Link layer — detects corruption
  - **ACK + Retry** — client retries 3 times with 2s timeout if no acknowledgment
  - **Sequence numbers** — chunks numbered with SEQ/TOTAL headers
  - **Chunking** — large messages split into 100-byte pieces, reassembled on server
- Interactive client — takes terminal input, sends messages without restart
- Session persistence — same session ID reused across all messages in a conversation
- Server simulation toggles — configurable packet drop and corruption via constants
- Test suite covering all layers, roundtrip encode/decode, chunking, out-of-order reassembly, and checksum corruption detection

## What's Been Learned
- **OSI model**: data flows top-down (Layer 7 to 1) on sender, bottom-up (1 to 7) on receiver. Each layer adds its own header (encapsulation)
- **UDP**: connectionless, no guarantees — just fire bytes at an address
- **Sockets**: how programs bind to ports and send/receive bytes
- **ARP**: how devices discover each other on a local network (IP + MAC mapping), arp cache vs live state
- **Packet loss**: UDP packets can vanish silently — client has no idea
- **Data corruption**: flipped bytes crash the decoder — no built-in integrity check
- **Large messages**: UDP has size limits — no built-in chunking
- **Why TCP exists**: solves all three problems above with checksums, ACKs, retries, sequencing
- **TCP/IP model**: real internet uses 4 layers (Application, Transport, Internet, Network Access), not 7
- **Session layer**: tracks conversations via session IDs — how a server distinguishes multiple clients
- **Checksums**: CRC32 hash detects corruption — sender calculates, receiver verifies
- **ACK/Retry**: sender waits for acknowledgment, retransmits on timeout
- **Chunking + reassembly**: split large data into numbered segments, sort and combine on receiver
- **3-way handshake**: SYN → SYN-ACK → ACK — how TCP establishes connections before data flows
- **4-way teardown**: FIN → FIN-ACK → FIN → ACK — how TCP closes connections cleanly
- **Connection state**: both sides must agree on connection status before sending data

## What's Next
- DNS resolver — how domain names become IP addresses
- SWIM gossip protocol — how nodes discover and monitor each other

## Project Structure
```
osi/
├── cmd/
│   ├── client/main.go          -- UDP client (sender, encapsulates, interactive input)
│   └── server/main.go          -- UDP server (receiver, decapsulates, ACK, simulation toggles)
├── internal/
│   ├── layers/
│   │   ├── 01_physical.go      -- bytes to binary string
│   │   ├── 02_datalink.go      -- MAC headers + CRC32 checksum
│   │   ├── 03_network.go       -- IP headers
│   │   ├── 04_transport.go     -- port headers + chunking/dechunking
│   │   ├── 05_session.go       -- session ID management
│   │   ├── 06_presentation.go  -- Base64 encode/decode
│   │   ├── 07_application.go   -- protocol tag
│   │   └── layers_test.go      -- test suite
│   └── models/
│       └── server.go           -- Machine struct (IP, MAC, Port)
└── go.mod
```

## Approach
- User writes the code, Claude guides step by step
- No code handouts — explain the concept, user implements
- Use real devices on local WiFi for demonstrations
- Go language, production project layout (cmd/ + internal/)

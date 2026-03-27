# UDP Networking from Scratch

Learning networking fundamentals by building them in Go — not reading textbooks. Each concept is implemented from scratch over real UDP sockets, tested across two physical machines on the same WiFi.

## What's Built

A full OSI-stack message pipeline: raw text enters at Layer 7, gets wrapped by each layer on the way down, travels over UDP, then gets unwrapped layer by layer on the way back up.

On top of that: the reliability features that make TCP what it is — checksums, ACKs, retries, chunking, sequence numbers, 3-way handshake, and 4-way teardown — all built manually over UDP.

## OSI Layer Pipeline

| Layer | File | What it does |
|-------|------|--------------|
| 7 — Application | `01_physical.go` → `07_application.go` | Tags data with protocol identifier (`CHAT`) |
| 6 — Presentation | `06_presentation.go` | Base64 encode/decode |
| 5 — Session | `05_session.go` | Attaches session ID to track conversations |
| 4 — Transport | `04_transport.go` | Splits data into 100-byte chunks, numbers them with SEQ/TOTAL headers |
| 3 — Network | `03_network.go` | Wraps IP source/destination headers |
| 2 — Data Link | `02_datalink.go` | Adds MAC addresses + CRC32 checksum |
| 1 — Physical | `01_physical.go` | Encodes bytes to binary string |

Data flows top-down on the sender (Layer 7 → 1), bottom-up on the receiver (Layer 1 → 7). Each layer adds its own header — that's encapsulation.

## TCP Features Built Over UDP

| Feature | How it works |
|---------|-------------|
| Checksums | CRC32 hash in the Data Link layer — sender calculates, receiver verifies; mismatch = corrupted packet |
| ACK + Retry | Client waits 2s for acknowledgment per chunk, retries up to 3 times on timeout |
| Sequence numbers | Each chunk carries `SEQ:n` and `TOTAL:n` headers |
| Chunking | Messages split into 100-byte pieces, reassembled in order on the server |
| 3-way handshake | `SYN → SYN-ACK → ACK` before any data flows |
| 4-way teardown | `FIN → FIN-ACK → FIN → ACK` on `exit` |
| Session persistence | Same session ID reused across all messages in a conversation |

## Project Structure

```
osi/
├── cmd/
│   ├── client/main.go      — interactive UDP client (encapsulates, sends, handles ACKs)
│   └── server/main.go      — UDP server (receives, decapsulates, sends ACKs)
├── internal/
│   ├── layers/
│   │   ├── 01_physical.go
│   │   ├── 02_datalink.go
│   │   ├── 03_network.go
│   │   ├── 04_transport.go
│   │   ├── 05_session.go
│   │   ├── 06_presentation.go
│   │   ├── 07_application.go
│   │   └── layers_test.go
│   └── models/
│       └── server.go       — Machine struct (IP, MAC, Port)
└── go.mod
```

## Running It

Start the server:
```bash
go run ./cmd/server
```

Start the client (in a separate terminal):
```bash
go run ./cmd/client
```

Type a message and press Enter. Type `exit` to close the connection cleanly via 4-way teardown.

### Simulating Network Issues

Toggle these constants in `cmd/server/main.go` to test reliability features:

```go
EnablePacketDrop  = true   // randomly drops packets
PacketDropPercent = 30     // drop 30% of packets

EnableCorruption    = true   // randomly corrupts packets
CorruptionPercent   = 20     // corrupt 20% of packets
CorruptedBytesCount = 5      // flip 5 bytes per corrupted packet
```

With drops enabled, you'll see the client retry and eventually give up. With corruption enabled, the Data Link layer will catch the CRC32 mismatch and reject the packet.

## Running Tests

```bash
go test ./internal/layers/...
```

Tests cover: all layer encode/decode roundtrips, chunking, out-of-order reassembly, and checksum corruption detection.

## Key Concepts Learned

- **UDP**: connectionless, no guarantees — just fire bytes at an address
- **Why TCP exists**: UDP loses packets silently, delivers them out of order, and has no integrity checks — TCP solves all three
- **Sockets**: how programs bind to ports and send/receive raw bytes
- **ARP**: how devices on a local network discover each other's MAC addresses
- **Encapsulation**: each OSI layer wraps data with its own header; the receiver peels them off in reverse
- **CRC32 checksums**: detect corruption by comparing sender-calculated vs receiver-calculated hash
- **3-way handshake / 4-way teardown**: how TCP establishes and cleanly closes connections
- **Session layer**: session IDs let a server distinguish multiple simultaneous clients

## What's Next

- DNS resolver — how domain names become IP addresses
- SWIM gossip protocol — how nodes discover and monitor each other in a distributed system

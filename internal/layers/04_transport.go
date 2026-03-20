package layers

import (
	"bytes"
	"fmt"
	"osi/internal/models"
	"sort"
)

const ChunkSize = 100

// TransportEncode adds port headers to data
func TransportEncode(data []byte, source models.Machine,
	destination models.Machine) []byte {
	header := fmt.Sprintf("SRC:%d|DST:%d|", source.Port,
		destination.Port)
	return append([]byte(header), data...)
}

// TransportChunk splits data into numbered chunks with port headers
func TransportChunk(data []byte, source models.Machine,
	destination models.Machine) [][]byte {
	total := (len(data) + ChunkSize - 1) / ChunkSize
	chunks := make([][]byte, 0, total)

	for i := 0; i < len(data); i += ChunkSize {
		end := min(i+ChunkSize, len(data))
		seq := i / ChunkSize
		header := fmt.Sprintf("SRC:%d|DST:%d|SEQ:%d|TOTAL:%d|",
			source.Port, destination.Port, seq, total)
		chunk := append([]byte(header), data[i:end]...)
		chunks = append(chunks, chunk)
	}

	return chunks
}

// TransportDecode strips port headers from a single packet
func TransportDecode(data []byte) []byte {
	parts := bytes.SplitN(data, []byte("|"), 3)
	fmt.Println("Layer 4 (Transport):", string(parts[0]), "|", string(parts[1]))
	return parts[2]
}

// TransportDechunk reassembles chunks into original data
func TransportDechunk(chunks [][]byte) []byte {
	type numbered struct {
		seq  int
		data []byte
	}

	var ordered []numbered

	for _, chunk := range chunks {
		// strip SRC port
		parts := bytes.SplitN(chunk, []byte("|"), 5)
		fmt.Println("Layer 4 (Transport):", string(parts[0]), "|", string(parts[1]),
			"|", string(parts[2]), "|", string(parts[3]))

		var seq int
		fmt.Sscanf(string(parts[2]), "SEQ:%d", &seq)

		ordered = append(ordered, numbered{seq: seq, data: parts[4]})
	}

	sort.Slice(ordered, func(i, j int) bool {
		return ordered[i].seq < ordered[j].seq
	})

	var result []byte
	for _, item := range ordered {
		result = append(result, item.data...)
	}

	return result
}

package layers

import (
	"bytes"
	"fmt"
)

func ApplicationEncode(data []byte, protocol string) []byte {
	return append([]byte(protocol+"|"), data...)
}

func ApplicationDecode(data []byte) []byte {
	parts := bytes.SplitN(data, []byte("|"), 2)
	fmt.Println("Layer 7 (Application): protocol =", string(parts[0]))
	return parts[1]
}

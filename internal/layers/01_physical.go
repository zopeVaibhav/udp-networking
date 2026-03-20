package layers

import (
	"fmt"
	"strings"
)

func BytesToBinary(data []byte) string {
	var sb strings.Builder
	for _, b := range data {
		fmt.Fprintf(&sb, "%08b ", b)
	}
	return sb.String()
}

func PhysicalEncode(data []byte) string {
	return BytesToBinary(data)
}

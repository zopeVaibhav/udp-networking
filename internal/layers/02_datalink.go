package layers

import (
	"bytes"
	"fmt"
	"hash/crc32"
	"osi/internal/models"
)

func DatalinkEncode(data []byte, source models.Machine,
	destination models.Machine) []byte {
	header := []byte("SRC:" + source.MACAddress + "|DST:" + destination.MACAddress + "|")
	finalData := append([]byte(header), data...)

	checksum := crc32.ChecksumIEEE(finalData)
	return append(fmt.Appendf(nil, "CRC:%d|", checksum), finalData...)
}

func DatalinkDecode(data []byte) []byte {
	parts := bytes.SplitN(data, []byte("|"), 2)
	receivedCRC := string(parts[0])
	rest := parts[1]
	actual := fmt.Sprintf("CRC:%d", crc32.ChecksumIEEE(rest))
	if receivedCRC != actual {
		fmt.Println(`Layer 2 (Data Link): CHECKSUM MISMATCH - packet corrupted!`)
		return nil
	}
	macParts := bytes.SplitN(rest, []byte("|"), 3)
	fmt.Println("Layer 2 (Data Link):", string(macParts[0]), "|", string(macParts[1]))
	return macParts[2]
}

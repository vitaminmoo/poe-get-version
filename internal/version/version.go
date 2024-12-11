package version

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"strings"
)

const (
	mock          = false
	poePatchHost  = "patch.pathofexile.com:12995"
	poe2PatchHost = "patch.pathofexile2.com:13060"
)

var poe1example = []byte{
	0x2, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
	0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
	0x0, 0x0, 0x0, 0x0, 0x22, 0x68, 0x0, 0x74, 0x0, 0x74, 0x0, 0x70, 0x0, 0x73,
	0x0, 0x3a, 0x0, 0x2f, 0x0, 0x2f, 0x0, 0x70, 0x0, 0x61, 0x0, 0x74, 0x0, 0x63,
	0x0, 0x68, 0x0, 0x2e, 0x0, 0x70, 0x0, 0x6f, 0x0, 0x65, 0x0, 0x63, 0x0, 0x64,
	0x0, 0x6e, 0x0, 0x2e, 0x0, 0x63, 0x0, 0x6f, 0x0, 0x6d, 0x0, 0x2f, 0x0, 0x33,
	0x0, 0x2e, 0x0, 0x32, 0x0, 0x35, 0x0, 0x2e, 0x0, 0x33, 0x0, 0x2e, 0x0, 0x34,
	0x0, 0x2f, 0x0, 0x0, 0x22, 0x68, 0x0, 0x74, 0x0, 0x74, 0x0, 0x70, 0x0, 0x73,
	0x0, 0x3a, 0x0, 0x2f, 0x0, 0x2f, 0x0, 0x70, 0x0, 0x61, 0x0, 0x74, 0x0, 0x63,
	0x0, 0x68, 0x0, 0x2e, 0x0, 0x70, 0x0, 0x6f, 0x0, 0x65, 0x0, 0x63, 0x0, 0x64,
	0x0, 0x6e, 0x0, 0x2e, 0x0, 0x63, 0x0, 0x6f, 0x0, 0x6d, 0x0, 0x2f, 0x0, 0x33,
	0x0, 0x2e, 0x0, 0x32, 0x0, 0x35, 0x0, 0x2e, 0x0, 0x33, 0x0, 0x2e, 0x0, 0x34,
	0x0, 0x2f, 0x0,
}

var poe2example = []byte{
	0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x28, 0x68,
	0x00, 0x74, 0x00, 0x74, 0x00, 0x70, 0x00, 0x73, 0x00, 0x3a, 0x00, 0x2f,
	0x00, 0x2f, 0x00, 0x70, 0x00, 0x61, 0x00, 0x74, 0x00, 0x63, 0x00, 0x68,
	0x00, 0x2d, 0x00, 0x70, 0x00, 0x6f, 0x00, 0x65, 0x00, 0x32, 0x00, 0x2e,
	0x00, 0x70, 0x00, 0x6f, 0x00, 0x65, 0x00, 0x63, 0x00, 0x64, 0x00, 0x6e,
	0x00, 0x2e, 0x00, 0x63, 0x00, 0x6f, 0x00, 0x6d, 0x00, 0x2f, 0x00, 0x34,
	0x00, 0x2e, 0x00, 0x31, 0x00, 0x2e, 0x00, 0x30, 0x00, 0x2e, 0x00, 0x36,
	0x00, 0x2e, 0x00, 0x32, 0x00, 0x2f, 0x00, 0x00, 0x28, 0x68, 0x00, 0x74,
	0x00, 0x74, 0x00, 0x70, 0x00, 0x73, 0x00, 0x3a, 0x00, 0x2f, 0x00, 0x2f,
	0x00, 0x70, 0x00, 0x61, 0x00, 0x74, 0x00, 0x63, 0x00, 0x68, 0x00, 0x2d,
	0x00, 0x70, 0x00, 0x6f, 0x00, 0x65, 0x00, 0x32, 0x00, 0x2e, 0x00, 0x70,
	0x00, 0x6f, 0x00, 0x65, 0x00, 0x63, 0x00, 0x64, 0x00, 0x6e, 0x00, 0x2e,
	0x00, 0x63, 0x00, 0x6f, 0x00, 0x6d, 0x00, 0x2f, 0x00, 0x34, 0x00, 0x2e,
	0x00, 0x31, 0x00, 0x2e, 0x00, 0x30, 0x00, 0x2e, 0x00, 0x36, 0x00, 0x2e,
	0x00, 0x32, 0x00, 0x2f, 0x00,
}

func Poe() (string, error) {
	var result []byte
	var err error
	if mock {
		result = poe1example
	} else {
		result, err = get(poePatchHost, []byte{1, 6})
		if err != nil {
			return "", fmt.Errorf("sending: %w", err)
		}
	}
	version, err := parse(result)
	if err != nil {
		return "", fmt.Errorf("parsing poe: %w", err)
	}
	if !strings.HasPrefix(version, "3.") {
		return "", fmt.Errorf("unexpected version: %s", version)
	}
	return version, nil
}

func Poe2() (string, error) {
	var result []byte
	var err error
	if mock {
		result = poe2example
	} else {
		result, err = get(poe2PatchHost, []byte{1, 7})
		if err != nil {
			return "", fmt.Errorf("sending: %w", err)
		}
	}
	version, err := parse(result)
	if err != nil {
		return "", fmt.Errorf("parsing poe2: %w", err)
	}
	if !strings.HasPrefix(version, "4.") {
		return "", fmt.Errorf("unexpected version: %s", version)
	}
	return version, nil
}

func header(r *bytes.Reader) error {
	var protoVer uint8
	var unknown [4]uint64
	if err := binary.Read(r, binary.BigEndian, &protoVer); err != nil {
		return fmt.Errorf("reading protocol version: %w", err)
	}
	if err := binary.Read(r, binary.BigEndian, &unknown); err != nil {
		return fmt.Errorf("reading unknown: %w", err)
	}
	if protoVer != 2 {
		return fmt.Errorf("unexpected protocol version: %d", protoVer)
	}
	return nil
}

func lenString(r *bytes.Reader) (string, error) {
	var len uint16
	if err := binary.Read(r, binary.BigEndian, &len); err != nil {
		return "", fmt.Errorf("reading length: %w", err)
	}
	len *= 2 // utf16 :/
	buf := make([]byte, len)
	if _, err := r.Read(buf); err != nil {
		return "", fmt.Errorf("reading string: %w", err)
	}
	return utf16ToString(buf), nil
}

func parse(data []byte) (string, error) {
	// hexdump(data)
	r := bytes.NewReader(data)
	err := header(r)
	if err != nil {
		return "", fmt.Errorf("reading header: %w", err)
	}
	str, err := lenString(r)
	if err != nil {
		return "", fmt.Errorf("reading length-string: %w", err)
	}
	str, err = lenString(r)
	if err != nil {
		return "", fmt.Errorf("reading length-string: %w", err)
	}
	str = strings.TrimPrefix(str, "https://patch")
	str = strings.TrimPrefix(str, "-poe2")
	str = strings.TrimPrefix(str, ".poecdn.com/")
	str = strings.TrimSuffix(str, "/")
	return str, nil
}

func get(host string, data []byte) ([]byte, error) {
	conn, err := net.Dial("tcp", host)
	if err != nil {
		return nil, fmt.Errorf("dialing: %w", err)
	}
	defer conn.Close()
	_, err = conn.Write(data)
	if err != nil {
		return nil, fmt.Errorf("writing: %w", err)
	}
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		return nil, fmt.Errorf("reading: %w", err)
	}
	return buf[:n], nil
}

func utf16ToString(data []byte) string {
	runes := make([]rune, len(data)/2)
	for i := 0; i < len(data); i += 2 {
		runes[i/2] = rune(data[i]) | (rune(data[i+1]) << 8)
	}
	return string(runes)
}

func hexdump(data []byte) {
	for i := 0; i < len(data); i += 16 {
		fmt.Printf("%08x  ", i)

		// Hex dump
		for j := 0; j < 16; j++ {
			if i+j < len(data) {
				fmt.Printf("%02x ", data[i+j])
			} else {
				fmt.Print("   ")
			}
			if j == 7 {
				fmt.Print(" ")
			}
		}

		// ASCII dump
		fmt.Print(" |")
		for j := 0; j < 16; j++ {
			if i+j < len(data) {
				if data[i+j] >= 32 && data[i+j] <= 126 {
					fmt.Printf("%c", data[i+j])
				} else {
					fmt.Print(".")
				}
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Println("|")
	}
}

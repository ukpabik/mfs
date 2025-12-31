package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
)

const (
	CREATE byte = 0x00
	DELETE byte = 0x01
	READ   byte = 0x02
	WRITE  byte = 0x03
)

// Internal client func for testing
func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:3000")
	if err != nil {
		log.Fatalf("dial: %v", err)
	}
	defer conn.Close()

	if err := sendCreate(conn, "testfile.txt"); err != nil {
		log.Fatalf("create: %v", err)
	}

	if err := sendDelete(conn, "testfile.txt"); err != nil {
		log.Fatalf("delete: %v", err)
	}
}

func sendCreate(conn net.Conn, filePath string) error {
	var payload bytes.Buffer
	payload.WriteByte(CREATE)
	binary.Write(&payload, binary.BigEndian, uint64(len(filePath)))
	payload.WriteString(filePath)

	fmt.Printf("sending CREATE request for: %s\n", filePath)

	if _, err := conn.Write(payload.Bytes()); err != nil {
		return fmt.Errorf("write: %w", err)
	}

	respBuf := make([]byte, 1024)
	n, err := conn.Read(respBuf)
	if err != nil && err != io.EOF {
		return fmt.Errorf("read response: %w", err)
	}

	if n > 0 {
		fmt.Printf("CREATE response: %s\n", string(respBuf[:n]))
	} else {
		fmt.Println("CREATE: no response")
	}

	return nil
}

func sendDelete(conn net.Conn, filePath string) error {
	var payload bytes.Buffer
	payload.WriteByte(DELETE)
	binary.Write(&payload, binary.BigEndian, uint64(len(filePath)))
	payload.WriteString(filePath)

	fmt.Printf("sending DELETE request for: %s\n", filePath)

	if _, err := conn.Write(payload.Bytes()); err != nil {
		return fmt.Errorf("write: %w", err)
	}

	respBuf := make([]byte, 1024)
	n, err := conn.Read(respBuf)
	if err != nil && err != io.EOF {
		return fmt.Errorf("read response: %w", err)
	}

	if n > 0 {
		fmt.Printf("DELETE response: %s\n", string(respBuf[:n]))
	} else {
		fmt.Println("DELETE: no response")
	}

	return nil
}

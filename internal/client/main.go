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

	if err := sendWrite(conn, "testfile.txt", []byte("hello world")); err != nil {
		log.Fatalf("write: %v", err)
	}

	if err := sendRead(conn, "testfile.txt", 0); err != nil {
		log.Fatalf("read: %v", err)
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
	}

	return nil
}

// sendRead sends a READ request and prints the returned file data.
func sendRead(conn net.Conn, filePath string, size uint64) error {
	var payload bytes.Buffer
	payload.WriteByte(READ)
	binary.Write(&payload, binary.BigEndian, uint64(len(filePath)))
	payload.WriteString(filePath)
	binary.Write(&payload, binary.BigEndian, size)

	fmt.Printf("sending READ request for: %s (size=%d)\n", filePath, size)

	if _, err := conn.Write(payload.Bytes()); err != nil {
		return fmt.Errorf("write: %w", err)
	}

	respBuf := make([]byte, 1024*1024)
	n, err := conn.Read(respBuf)
	if err != nil && err != io.EOF {
		return fmt.Errorf("read response: %w", err)
	}

	if n > 0 {
		fmt.Printf("READ response: %s\n", string(respBuf[:n]))
	}

	return nil
}

// sendWrite sends a WRITE request with data to be written to the file.
func sendWrite(conn net.Conn, filePath string, data []byte) error {
	var payload bytes.Buffer
	payload.WriteByte(WRITE)
	binary.Write(&payload, binary.BigEndian, uint64(len(filePath)))
	payload.WriteString(filePath)
	binary.Write(&payload, binary.BigEndian, uint64(len(data)))
	payload.Write(data)

	fmt.Printf("sending WRITE request for: %s (size=%d bytes)\n", filePath, len(data))

	if _, err := conn.Write(payload.Bytes()); err != nil {
		return fmt.Errorf("write: %w", err)
	}

	respBuf := make([]byte, 1024)
	n, err := conn.Read(respBuf)
	if err != nil && err != io.EOF {
		return fmt.Errorf("read response: %w", err)
	}

	if n > 0 {
		fmt.Printf("WRITE response: %s\n", string(respBuf[:n]))
	}

	return nil
}

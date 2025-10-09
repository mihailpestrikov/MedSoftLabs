package hl7

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net"
)

const (
	MLLP_START = 0x0B
	MLLP_END1  = 0x1C
	MLLP_END2  = 0x0D
)

type MLLPListener struct {
	listener net.Listener
	handler  func([]byte) []byte
}

func NewMLLPListener(port string, certPath string, keyPath string, handler func([]byte) []byte) (*MLLPListener, error) {
	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load TLS certificate: %w", err)
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS12,
	}

	listener, err := tls.Listen("tcp", ":"+port, tlsConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to start TLS listener: %w", err)
	}

	log.Printf("MLLP/TLS listener started on port %s", port)

	return &MLLPListener{
		listener: listener,
		handler:  handler,
	}, nil
}

func (ml *MLLPListener) Start() {
	for {
		conn, err := ml.listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			continue
		}

		go ml.handleConnection(conn)
	}
}

func (ml *MLLPListener) handleConnection(conn net.Conn) {
	defer conn.Close()
	log.Printf("New MLLP/TLS connection from %s", conn.RemoteAddr())

	reader := bufio.NewReader(conn)

	for {
		message, err := ml.readMLLPMessage(reader)
		if err != nil {
			if err != io.EOF {
				log.Printf("Error reading MLLP message: %v", err)
			}
			return
		}

		ack := ml.handler(message)

		if err := ml.writeMLLPMessage(conn, ack); err != nil {
			log.Printf("Error sending ACK: %v", err)
			return
		}
	}
}

func (ml *MLLPListener) readMLLPMessage(reader *bufio.Reader) ([]byte, error) {
	startByte, err := reader.ReadByte()
	if err != nil {
		return nil, err
	}

	if startByte != MLLP_START {
		return nil, fmt.Errorf("invalid start byte: expected 0x%02X, got 0x%02X", MLLP_START, startByte)
	}

	var message []byte
	for {
		b, err := reader.ReadByte()
		if err != nil {
			return nil, err
		}

		if b == MLLP_END1 {
			end2, err := reader.ReadByte()
			if err != nil {
				return nil, err
			}
			if end2 == MLLP_END2 {
				break
			}
			message = append(message, b, end2)
		} else {
			message = append(message, b)
		}
	}

	return message, nil
}

func (ml *MLLPListener) writeMLLPMessage(conn net.Conn, message []byte) error {
	frame := make([]byte, 0, len(message)+3)
	frame = append(frame, MLLP_START)
	frame = append(frame, message...)
	frame = append(frame, MLLP_END1, MLLP_END2)

	_, err := conn.Write(frame)
	return err
}

func (ml *MLLPListener) Close() error {
	return ml.listener.Close()
}

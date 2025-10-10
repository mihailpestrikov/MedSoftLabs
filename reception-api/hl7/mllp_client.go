package hl7

import (
	"bufio"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"os"
	"time"
)

const (
	MLLP_START = 0x0B
	MLLP_END1  = 0x1C
	MLLP_END2  = 0x0D
)

type MLLPClient struct {
	address   string
	tlsConfig *tls.Config
}

func NewMLLPClient(address string, certPath string) (*MLLPClient, error) {
	cert, err := os.ReadFile(certPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read certificate: %w", err)
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(cert) {
		return nil, fmt.Errorf("failed to parse certificate")
	}

	tlsConfig := &tls.Config{
		RootCAs: certPool,
	}

	return &MLLPClient{
		address:   address,
		tlsConfig: tlsConfig,
	}, nil
}

func (mc *MLLPClient) SendMessage(message []byte) ([]byte, error) {
	conn, err := tls.Dial("tcp", mc.address, mc.tlsConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}
	defer conn.Close()

	conn.SetDeadline(time.Now().Add(10 * time.Second))

	if err := mc.writeMLLPMessage(conn, message); err != nil {
		return nil, fmt.Errorf("failed to send message: %w", err)
	}

	ack, err := mc.readMLLPMessage(bufio.NewReader(conn))
	if err != nil {
		return nil, fmt.Errorf("failed to read ACK: %w", err)
	}

	return ack, nil
}

func (mc *MLLPClient) writeMLLPMessage(conn net.Conn, message []byte) error {
	frame := make([]byte, 0, len(message)+3)
	frame = append(frame, MLLP_START)
	frame = append(frame, message...)
	frame = append(frame, MLLP_END1, MLLP_END2)

	_, err := conn.Write(frame)
	return err
}

func (mc *MLLPClient) readMLLPMessage(reader *bufio.Reader) ([]byte, error) {
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

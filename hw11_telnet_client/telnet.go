package main

import (
	"fmt"
	"io"
	"net"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) *SimpleTelnetClient {
	return &SimpleTelnetClient{
		address:        address,
		connectTimeout: timeout,
		in:             in,
		out:            out,
		conn:           nil,
	}
}

type SimpleTelnetClient struct {
	address        string
	connectTimeout time.Duration
	in             io.ReadCloser
	out            io.Writer
	conn           net.Conn
}

func (s *SimpleTelnetClient) Connect() error {
	conn, err := net.DialTimeout("tcp", s.address, s.connectTimeout)
	if err != nil {
		return fmt.Errorf("failed to connect to %s: %w", s.address, err)
	}
	s.conn = conn
	return nil
}

func (s *SimpleTelnetClient) Close() error {
	if err := s.conn.Close(); err != nil {
		return fmt.Errorf("failed to close connection: %w", err)
	}
	return nil
}

func (s *SimpleTelnetClient) Send() error {
	_, err := io.Copy(s.conn, s.in)
	return err
}

func (s *SimpleTelnetClient) Receive() error {
	_, err := io.Copy(s.out, s.conn)
	return err
}

package after_db

import (
	"fmt"
	"net"
)

type Client struct {
	conn net.Conn
}

func NewClient(serverAddr string) (*Client, error) {
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		return nil, err
	}

	return &Client{conn: conn}, nil
}

func (c *Client) Close() {
	c.conn.Close()
}

func (c *Client) Get(key string) (string, error) {
	command := fmt.Sprintf("GET %s", key)
	_, err := c.conn.Write([]byte(command + "\n"))
	if err != nil {
		return "", err
	}

	// Leer la respuesta del servidor
	response := make([]byte, 512)
	n, err := c.conn.Read(response)
	if err != nil {
		return "", err
	}

	return string(response[:n]), nil
}

func (c *Client) Save(key, value string) error {
	command := fmt.Sprintf("SAVE %s %s", key, value)
	_, err := c.conn.Write([]byte(command + "\n"))
	if err != nil {
		return err
	}

	// Leer la respuesta del servidor
	response := make([]byte, 512)
	n, err := c.conn.Read(response)
	if err != nil {
		return err
	}

	if string(response[:n]) != "SAVE exitoso\n" {
		return fmt.Errorf("Error en SAVE: %s", string(response[:n]))
	}

	return nil
}

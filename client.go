package after_db

import (
	"errors"
	"net"
	"sync"
)

var connections sync.Map

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

func (c *Client) GetByKey(key string, value any) (any, error) {
	value, found := connections.Load(key)
	if found {
		return value, errors.New("error")
	}
	return nil, errors.New("key not found")
}

func (c *Client) GetKeyByValue(valueToFind any) (any, error) {
	var keyValue any

	connections.Range(func(key, value any) bool {
		if value == valueToFind {
			keyValue = key
		}
		return true
	})

	if keyValue != nil {
		return keyValue, nil
	}
	return nil, errors.New("value not found finding by value")
}

func (c *Client) Save(key, value string) error {
	connections.Store(key, value)
	return nil
}

func (c *Client) Delete(key string) error {
	connections.Delete(key)
	return nil
}

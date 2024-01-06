package after_db

import (
	"encoding/json"
	"fmt"
	"net"
)

type Client struct {
	conn net.Conn
	dec  *json.Decoder
	enc  *json.Encoder
}

func NewClient(addr string, port string) (*Client, error) {
	c := new(Client)

	var err error
	c.conn, err = net.Dial("tcp", fmt.Sprintf("%s:%s", addr, port))
	if err != nil {
		return nil, err
	}

	c.dec = json.NewDecoder(c.conn)
	c.enc = json.NewEncoder(c.conn)
	return c, nil
}

func (c *Client) Insert(key string, value interface{}) (int, error) {
	msg := map[string]interface{}{
		"op":    "insert",
		"key":   key,
		"value": value,
	}

	err := c.enc.Encode(msg)
	if err != nil {
		return 0, err
	}

	var resp map[string]interface{}
	err = c.dec.Decode(&resp)
	if err != nil {
		return 0, err
	}

	code, ok := resp["code"].(float64)
	if !ok {
		return 0, fmt.Errorf("invalid code")
	}
	return int(code), nil
}

func (c *Client) SearchByKey(key string) (interface{}, int, error) {
	msg := map[string]interface{}{
		"op":  "search_by_key",
		"key": key,
	}

	err := c.enc.Encode(msg)
	if err != nil {
		return nil, 0, err
	}

	var resp map[string]interface{}
	err = c.dec.Decode(&resp)
	if err != nil {
		return nil, 0, err
	}

	code, ok := resp["code"].(float64)
	if !ok {
		return nil, 0, fmt.Errorf("invalid code")
	}
	value, ok := resp["value"]
	if !ok {
		return nil, int(code), nil
	}
	return value, int(code), nil
}

func (c *Client) SearchByValue(value string) (interface{}, int, error) {
	msg := map[string]interface{}{
		"op":    "search_by_value",
		"value": value,
	}

	err := c.enc.Encode(msg)
	if err != nil {
		return nil, 0, err
	}

	var resp map[string]interface{}
	err = c.dec.Decode(&resp)
	if err != nil {
		return nil, 0, err
	}

	code, ok := resp["code"].(float64)
	if !ok {
		return nil, 0, fmt.Errorf("invalid code")
	}
	key, ok := resp["key"]
	if !ok {
		return nil, int(code), nil
	}
	return key, int(code), nil
}

// Método no utilizado por ahora
func (c *Client) Update(key string, value interface{}) (int, error) {
	msg := map[string]interface{}{
		"op":    "update",
		"key":   key,
		"value": value,
	}

	err := c.enc.Encode(msg)
	if err != nil {
		return 0, err
	}

	var resp map[string]interface{}
	err = c.dec.Decode(&resp)
	if err != nil {
		return 0, err
	}

	code, ok := resp["code"].(float64)
	if !ok {
		return 0, fmt.Errorf("invalid code")
	}
	return int(code), nil
}

func (c *Client) Delete(key string) (int, error) {
	msg := map[string]interface{}{
		"op":  "delete",
		"key": key,
	}

	err := c.enc.Encode(msg)
	if err != nil {
		return 0, err
	}

	var resp map[string]interface{}
	err = c.dec.Decode(&resp)
	if err != nil {
		return 0, err
	}

	code, ok := resp["code"].(float64)
	if !ok {
		return 0, fmt.Errorf("invalid code")
	}
	return int(code), nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

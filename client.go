package after_db

import (
	"encoding/json"
	"fmt"
	"net"
)

// Estructura que representa al cliente del servicio
type Client struct {
	conn net.Conn
	dec  *json.Decoder
	enc  *json.Encoder
}

// Función que crea un nuevo cliente, dado la dirección y el puerto del servicio
func NewClient(addr string, port string) (*Client, error) {
	c := new(Client)
	var err error
	// Crear una conexión TCP con el servicio, usando la función net.Dial
	c.conn, err = net.Dial("tcp", fmt.Sprintf("%s:%s", addr, port))
	if err != nil {
		return nil, err
	}
	// Crear un decodificador y un codificador JSON para leer y escribir datos en la conexión
	c.dec = json.NewDecoder(c.conn)
	c.enc = json.NewEncoder(c.conn)
	return c, nil
}

func (c *Client) Insert(key string, value interface{}) (int, error) {
	// Crear un mensaje con la operación, la clave y el valor
	msg := map[string]interface{}{
		"op":    "insert",
		"key":   key,
		"value": value,
	}
	// Codificar el mensaje en formato JSON y escribirlo en la conexión
	err := c.enc.Encode(msg)
	if err != nil {
		return 0, err
	}
	// Leer una respuesta del servicio
	var resp map[string]interface{}
	err = c.dec.Decode(&resp)
	if err != nil {
		return 0, err
	}
	// Obtener el código de la respuesta
	code, ok := resp["code"].(float64)
	if !ok {
		return 0, fmt.Errorf("Código inválido")
	}
	return int(code), nil
}

// Método que busca un valor en el servicio, dado una clave
func (c *Client) Search(key string) (interface{}, int, error) {
	// Crear un mensaje con la operación y la clave
	msg := map[string]interface{}{
		"op":  "search",
		"key": key,
	}
	// Codificar el mensaje en formato JSON y escribirlo en la conexión
	err := c.enc.Encode(msg)
	if err != nil {
		return nil, 0, err
	}
	// Leer una respuesta del servicio
	var resp map[string]interface{}
	err = c.dec.Decode(&resp)
	if err != nil {
		return nil, 0, err
	}
	// Obtener el código y el valor de la respuesta
	code, ok := resp["code"].(float64)
	if !ok {
		return nil, 0, fmt.Errorf("Código inválido")
	}
	value, ok := resp["value"]
	if !ok {
		return nil, int(code), nil
	}
	return value, int(code), nil
}

// Método que actualiza un valor en el servicio, dado una clave y un valor
func (c *Client) Update(key string, value interface{}) (int, error) {
	// Crear un mensaje con la operación, la clave y el valor
	msg := map[string]interface{}{
		"op":    "update",
		"key":   key,
		"value": value,
	}
	// Codificar el mensaje en formato JSON y escribirlo en la conexión
	err := c.enc.Encode(msg)
	if err != nil {
		return 0, err
	}
	// Leer una respuesta del servicio
	var resp map[string]interface{}
	err = c.dec.Decode(&resp)
	if err != nil {
		return 0, err
	}
	// Obtener el código de la respuesta
	code, ok := resp["code"].(float64)
	if !ok {
		return 0, fmt.Errorf("Código inválido")
	}
	return int(code), nil
}

// Método que elimina un valor del servicio, dado una clave
func (c *Client) Delete(key string) (int, error) {
	// Crear un mensaje con la operación y la clave
	msg := map[string]interface{}{
		"op":  "delete",
		"key": key,
	}
	// Codificar el mensaje en formato JSON y escribirlo en la conexión
	err := c.enc.Encode(msg)
	if err != nil {
		return 0, err
	}
	// Leer una respuesta del servicio
	var resp map[string]interface{}
	err = c.dec.Decode(&resp)
	if err != nil {
		return 0, err
	}
	// Obtener el código de la respuesta
	code, ok := resp["code"].(float64)
	if !ok {
		return 0, fmt.Errorf("Código inválido")
	}
	return int(code), nil
}

// Método que cierra la conexión con el servicio
func (c *Client) Close() error {
	return c.conn.Close()
}

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"sync"
)

type Database struct {
	dbserver sync.Map
}

func NewDatabase() *Database {
	return &Database{
		dbserver: sync.Map{},
	}
}

func (db *Database) Insert(key string, value interface{}) error {
	db.dbserver.Store(key, value)
	return nil
}

// Método que busca un valor en la base de datos, dado una clave
func (db *Database) SearchByKey(key string) (interface{}, bool) {
	value, found := db.dbserver.Load(key)
	if found {
		return value, true
	}
	return nil, false
}

// Método que actualiza un valor en la base de datos, dado una clave y un nuevo valor
func (db *Database) SearchByValue(keyToSearch string, value interface{}) (interface{}, error) {
	var keyValue interface{}
	db.dbserver.Range(func(key, value any) bool {
		if value == keyToSearch {
			keyValue = key
		}
		return true
	})

	if keyValue != nil {
		return keyValue, nil
	}
	return nil, errors.New("key not found syncmap")
}

// Método que elimina un valor de la base de datos, dado una clave
func (db *Database) Delete(key string) error {
	db.dbserver.Delete(key)
	return nil
}

func main() {
	StartService(NewDatabase(), "127.0.0.1", "12345")
}

func StartService(db *Database, ip string, port string) {
	// Crear un servidor TCP en la dirección y puerto especificados
	addr := fmt.Sprintf("%s:%s", ip, port)
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Iniciando el servicio en %s\n", addr)
	// Aceptar conexiones entrantes y manejarlas
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go handleConnection(conn, db)
	}
}

func handleConnection(conn net.Conn, db *Database) {
	defer conn.Close()
	log.Printf("Conectado con: %s\n", conn.RemoteAddr())
	// Crear un decodificador y un codificador JSON para leer y escribir datos en la conexión
	dec := json.NewDecoder(conn)
	enc := json.NewEncoder(conn)
	// Leer un mensaje del cliente
	var msg map[string]interface{}
	err := dec.Decode(&msg)
	if err != nil {
		log.Println(err)
		return
	}
	// Procesar el mensaje según el tipo de operación
	op, ok := msg["op"].(string)
	if !ok {
		log.Println("Operación inválida")
		return
	}
	switch op {
	case "insert":
		// Obtener la clave y el valor del mensaje
		key, ok := msg["key"].(string)
		if !ok {
			log.Println("Clave inválida")
			return
		}
		value, ok := msg["value"]
		if !ok {
			log.Println("Valor inválido")
			return
		}
		// Insertar el valor en la base de datos
		db.Insert(key, value)
		// Enviar una respuesta con el código 201 (Created)
		resp := map[string]interface{}{
			"code": 201,
		}
		err := enc.Encode(resp)
		if err != nil {
			log.Println(err)
			return
		}
	case "search":
		// Obtener la clave del mensaje
		key, ok := msg["key"].(string)
		if !ok {
			log.Println("Clave inválida")
			return
		}
		// Buscar el valor en la base de datos
		value, ok := db.SearchByKey(key)
		if !ok {
			// Enviar una respuesta con el código 404 (Not Found)
			resp := map[string]interface{}{
				"code": 404,
			}
			err := enc.Encode(resp)
			if err != nil {
				log.Println(err)
				return
			}
		} else {
			// Enviar una respuesta con el código 200 (OK) y el valor
			resp := map[string]interface{}{
				"code":  200,
				"value": value,
			}
			err := enc.Encode(resp)
			if err != nil {
				log.Println(err)
				return
			}
		}
	case "delete":
		// Obtener la clave del mensaje
		key, ok := msg["key"].(string)
		if !ok {
			log.Println("Clave inválida")
			return
		}
		// Eliminar el valor de la base de datos
		err = db.Delete(key)
		if err != nil {
			// Enviar una respuesta con el código 404 (Not Found)
			resp := map[string]interface{}{
				"code": 404,
			}
			err := enc.Encode(resp)
			if err != nil {
				log.Println(err)
				return
			}
		} else {
			// Enviar una respuesta con el código 200 (OK)
			resp := map[string]interface{}{
				"code": 200,
			}
			err := enc.Encode(resp)
			if err != nil {
				log.Println(err)
				return
			}
		}
	default:
		// Enviar una respuesta con el código 400 (Bad Request)
		resp := map[string]interface{}{
			"code": 400,
		}
		err := enc.Encode(resp)
		if err != nil {
			log.Println(err)
			return
		}
	}
}

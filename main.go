package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

const (
	CleanupInterval = 10 * time.Second
	Port            = "6379"
)

type Cache struct {
	Data   map[string]string
	Expiry map[string]time.Time
	mu     sync.RWMutex
}

func NewCache() *Cache {
	return &Cache{
		Data:   make(map[string]string),
		Expiry: make(map[string]time.Time),
	}
}

func (c *Cache) set(key string, value string, expiry time.Time) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.Data[key] = value
	if !expiry.IsZero() {
		c.Expiry[key] = expiry
	}
}

func (c *Cache) get(key string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	value, exists := c.Data[key]
	if !exists {
		return "", false
	}
	expiry, expExists := c.Expiry[key]
	if expExists && time.Now().After(expiry) {
		delete(c.Data, key)
		delete(c.Expiry, key)
		return "", false
	}
	return value, true
}

func (c *Cache) StartCleanup() {
	for range time.Tick(CleanupInterval) {
		now := time.Now()
		c.mu.Lock()
		for k, v := range c.Expiry {
			if now.After(v) {
				delete(c.Data, k)
				delete(c.Expiry, k)
			}
		}
		c.mu.Unlock()
	}
}

func main() {
	ln, err := net.Listen("tcp", ":"+Port)
	if err != nil {
		log.Fatalf("Could not listen to Port %s: %s\n", Port, err)
	}
	log.Printf("Listening to TCP connections on Port %s ...\n", Port)

	cache := NewCache()
	go cache.StartCleanup()

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go handleConnection(conn, cache)
	}
}

func handleConnection(conn net.Conn, cache *Cache) {
	defer func() {
		conn.Close()
		log.Printf("%s has disconnected.", conn.RemoteAddr())
	}()
	log.Printf("%s has connected.", conn.RemoteAddr())

	reader := bufio.NewReader(conn)
	for {
		fmt.Fprint(conn, "gedis-cli> ")
		command, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		command = strings.TrimSuffix(command, "\n")
		args := strings.Fields(command)

		if len(args) < 1 {
			fmt.Fprintln(conn, "Usage: command \"argument\"")
			continue
		}

		err = runCommand(args, cache, conn)
		if err != nil {
			fmt.Fprintln(conn, err)
		}
	}
}

func runCommand(args []string, cache *Cache, conn net.Conn) error {
	switch args[0] {
	case "PING":
		fmt.Fprintln(conn, "PONG")
	case "ECHO":
		echoMessage(args[1:], conn)
	case "SET":
		res, err := setKey(args, cache)
		if err != nil {
			return err
		}
		fmt.Fprintln(conn, res)
	case "GET":
		res, err := getKey(args, cache)
		if err != nil {
			return err
		}
		fmt.Fprintln(conn, res)
	default:
		return fmt.Errorf("unknown command: %s", args[0])
	}

	return nil
}

func echoMessage(args []string, conn net.Conn) {
	message := strings.Join(args, " ")
	fmt.Fprintln(conn, message)
}

func setKey(args []string, cache *Cache) (string, error) {
	if len(args) < 3 {
		return "", fmt.Errorf("SET requires at least a 'key' and a 'value'")
	}
	key, value := args[1], args[2]
	expiry := time.Time{}
	if len(args) > 3 {
		if len(args) != 5 {
			return "", fmt.Errorf("requires a second expiry parameter")
		}
		switch flag := args[3]; flag {
		case "EX":
			duration, err := time.ParseDuration(args[4] + "s")
			if err != nil {
				return "", fmt.Errorf("invalid seconds value for 'EX'")
			}
			expiry = time.Now().Add(duration)
		case "PX":
			duration, err := time.ParseDuration(args[4] + "ms")
			if err != nil {
				return "", fmt.Errorf("invalid seconds value for 'PX'")
			}
			expiry = time.Now().Add(duration)
		default:
			return "", fmt.Errorf("invalid argument for 'SET'")
		}
	}
	cache.set(key, value, expiry)
	return "OK", nil
}

func getKey(args []string, cache *Cache) (string, error) {
	if len(args) < 2 {
		return "", fmt.Errorf("GET requires a 'value' parameter")
	}
	value, ok := cache.get(args[1])
	if !ok {
		return "(nil)", nil
	} else {
		return value, nil
	}
}

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

const cleanupInterval = 10 * time.Second

type Cache struct {
	Data   map[string]string
	Expiry map[string]time.Time
}

func NewCache() *Cache {
	return &Cache{
		Data:   make(map[string]string),
		Expiry: make(map[string]time.Time),
	}
}

func (c *Cache) StartCleanup() {
	for range time.Tick(cleanupInterval) {
		now := time.Now()
		for k, v := range c.Expiry {
			if now.After(v) {
				delete(c.Data, k)
				delete(c.Expiry, k)
			}
		}
	}
}

func main() {
	cache := NewCache()
	go cache.StartCleanup()

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("% gedis-cli ")
		command, _ := reader.ReadString('\n')
		command = strings.TrimSuffix(command, "\n")
		args := strings.Fields(command)
		if len(args) < 2 {
			fmt.Println("Usage: command \"argument\"")
			continue
		}
		err := runCommand(args, cache)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func runCommand(args []string, cache *Cache) error {
	switch args[0] {
	case "PING":
		fmt.Println("PONG")
	case "ECHO":
		echoMessage(args[1:])
	case "SET":
		res, err := setKey(args, cache)
		if err != nil {
			return err
		}
		fmt.Println(res)
	case "GET":
		res := getKey(args, cache)
		fmt.Println(res)
	default:
		return fmt.Errorf("unknown command: %s", args[0])
	}

	return nil
}

func echoMessage(args []string) {
	message := strings.Join(args, " ")
	fmt.Println(message)
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
	cache.Data[key] = value
	if !expiry.IsZero() {
		cache.Expiry[key] = expiry
	}
	return "OK", nil
}

func getKey(args []string, cache *Cache) string {
	value, ok := cache.Data[args[1]]
	expiry, expExists := cache.Expiry[args[1]]
	if !ok {
		return "(nil)"
	} else if expExists && time.Now().After(expiry) {
		delete(cache.Data, args[1])
		delete(cache.Expiry, args[1])
		return "(nil)"
	} else {
		return value
	}
}

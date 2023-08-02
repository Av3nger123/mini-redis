package core

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type CacheRecord struct {
	value  string
	expiry time.Time
}

type Cache struct {
	data        map[string]*CacheRecord
	mutex       sync.RWMutex
	cleanup     time.Duration
	stopCleanup chan int
	file        string
}

func NewCache(filepath string, expiry time.Duration) (*Cache, error) {
	cache := &Cache{
		data:        make(map[string]*CacheRecord),
		file:        filepath,
		stopCleanup: make(chan int),
		cleanup:     expiry,
	}
	err := cache.loadFromDisk()
	if err != nil {
		return nil, err
	}
	go cache.persist()
	go cache.clear()
	return cache, nil
}

func (cache *Cache) StopClearingRecords() {
	cache.stopCleanup <- 1
}

// Clearing expired records form cache
func (cache *Cache) clear() {
	ticker := time.NewTicker(cache.cleanup)
	for {
		select {
		case <-ticker.C:
			cache.clearExpiredRecords()
		case <-cache.stopCleanup:
			ticker.Stop()
			return
		}
	}
}

func (cache *Cache) clearExpiredRecords() {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()
	now := time.Now()
	for key, record := range cache.data {
		if record.expiry.Before(now) {
			delete(cache.data, key)
		}
	}
}

// Persist Data into the disk from time to time
func (cache *Cache) persist() {
	for {
		time.Sleep(1 * time.Minute)
		err := cache.loadToDisk()
		if err != nil {
			fmt.Println(err)
		}
	}
}

func (cache *Cache) loadToDisk() error {
	file, err := os.Create(cache.file)
	if err != nil {
		return err
	}
	writer := bufio.NewWriter(file)
	for key, value := range cache.data {
		line := fmt.Sprintf("%v:%v:%v\n", key, value.value, value.expiry.Format(time.RFC3339))
		_, err := writer.WriteString(line)
		if err != nil {
			return err
		}
	}
	writer.Flush()
	return nil
}

// Load the data when the server comes up
func (cache *Cache) loadFromDisk() error {
	file, err := os.Open(cache.file)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("File doesn't exist")
			return nil
		}
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ":")
		expiry, err := time.Parse(time.RFC3339, parts[2])
		if err != nil {
			fmt.Println("Error in converting string to time.Time for key", parts[0])
		}
		cache.data[parts[0]] = &CacheRecord{
			value:  parts[1],
			expiry: expiry,
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

func (cache *Cache) HandleConnection(connection net.Conn) {
	defer connection.Close()

	reader := bufio.NewReader(connection)
	writer := bufio.NewWriter(connection)

	for {
		command, err := reader.ReadString('\n')
		if err != nil {
			if err.Error() == "EOF" {
				fmt.Println("Client closed the connection.")
				break
			}
			fmt.Println("Error reading command:", err)
		}
		response := cache.handleCommand(command, connection)
		bytesWritten, err := writer.WriteString(response + "\n")
		if err != nil {
			fmt.Println("Error writing response:", err)
		}
		fmt.Println(bytesWritten)

		writer.Flush()

	}
}

func (cache *Cache) handleCommand(command string, connection net.Conn) string {
	parts := strings.Split(strings.TrimSpace(command), " ")
	if len(parts) < 2 {
		return "Invalid command"
	}
	cmd := strings.ToUpper(parts[0])

	switch cmd {
	case "GET":
		return cache.get(parts[1])
	case "SET":
		return cache.put(parts)
	case "DELETE":
		return cache.delete(parts[1])
	default:
		return "Invalid Command"

	}
}

func (cache *Cache) get(key string) string {
	cache.mutex.RLock()
	defer cache.mutex.RUnlock()
	record, exists := cache.data[key]
	if !exists {
		return "NULL"
	}
	if record.expiry.Before(time.Now()) {
		delete(cache.data, key)
		return "NULL"
	}
	return record.value
}

func (cache *Cache) put(command []string) string {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()
	key := command[1]
	value := command[2]
	var ttl string
	if len(command) > 3 {
		ttl = command[3]
	}

	expiryTime, err := strconv.Atoi(ttl)
	if err != nil {
		return "Invalid TTL"
	}
	cache.data[key] = &CacheRecord{
		value:  value,
		expiry: time.Now().Add(time.Duration(expiryTime) * time.Second),
	}
	record, exists := cache.data[key]
	if !exists {
		return "NULL"
	}
	return record.value
}

func (cache *Cache) delete(key string) string {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	_, exists := cache.data[key]
	if !exists {
		return "Record doesn't exist"
	}

	delete(cache.data, key)
	return "Record removed"
}

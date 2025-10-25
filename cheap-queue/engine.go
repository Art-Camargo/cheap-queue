package cheapqueue

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func (c *CheapQueueEngine) Init(projectId string) {
	if projectId == "" {
		panic("projectId cannot be empty - use a unique identifier for your project")
	}
	
	c.projectId = projectId
	c.queues = make(map[string]chan Message)

	c.recoverMessagesFromTemp()
}

func (c *CheapQueueEngine) recoverMessagesFromTemp() {
	tmpDir := c.getQueueCacheDir()
	
	files, err := os.ReadDir(tmpDir)
	if err != nil {
		return 
	}
	
	messagesByQueue := make(map[string][]Message)

	prefix := c.projectId + "_"
	
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		
		filename := file.Name()

		if !strings.HasPrefix(filename, prefix) {
			continue
		}

		withoutPrefix := strings.TrimPrefix(filename, prefix)

		parts := strings.Split(withoutPrefix, "_")
		if len(parts) < 2 {
			continue
		}

		queueId := parts[0]

		filePath := filepath.Join(tmpDir, filename)
		data, err := os.ReadFile(filePath)
		if err != nil {
			continue 
		}
		
		msg := Message{
			Data:           data,
			MessageId:      filename,
			PathToTempFile: filePath,
		}
		
		messagesByQueue[queueId] = append(messagesByQueue[queueId], msg)
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	
	for queueId, messages := range messagesByQueue {
		bufferSize := max(len(messages), 10)
		
		c.queues[queueId] = make(chan Message, bufferSize)

		for _, msg := range messages {
			c.queues[queueId] <- msg
		}
	}
}

func (c *CheapQueueEngine) Bind(queueId string, bufferSize int) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if queue, exists := c.queues[queueId]; exists {
		if bufferSize > cap(queue) {
			newQueue := make(chan Message, bufferSize)
			close(queue)
			for msg := range queue {
				newQueue <- msg
			}
			c.queues[queueId] = newQueue
		}
		return nil 
	}

	c.queues[queueId] = make(chan Message, bufferSize)
	return nil
}

func (c *CheapQueueEngine) generateUniqueFilePath(queueId string) string {
	timestamp := strconv.FormatInt(time.Now().UnixNano(), 10)
	// Formato: projectId_queueId_timestamp
	baseFilename := fmt.Sprintf("%s_%s_%s", c.projectId, queueId, timestamp)
	filePath := filepath.Join(c.getQueueCacheDir(), baseFilename)

	for {
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			break
		}
		randomLetter := string(rune('a' + rand.Intn(26)))
		// Formato com letra: projectId_queueId_timestamp_letra
		filePath = filepath.Join(c.getQueueCacheDir(), fmt.Sprintf("%s_%s", baseFilename, randomLetter))
	}

	return filePath
}

func (c *CheapQueueEngine) saveMessageToTempFile(queueId string, data []byte) (string, error) {
	filePath := c.generateUniqueFilePath(queueId)
	err := os.WriteFile(filePath, data, 0644)
	if err != nil {
		return "", fmt.Errorf("failed to save message to temp file: %w", err)
	}
	return filePath, nil
}


func (c *CheapQueueEngine) Publish(queueId string, data []byte) error {
	c.mu.Lock()
	queue, exists := c.queues[queueId]
	c.mu.Unlock()

	if !exists {
		return fmt.Errorf("queue '%s' does not exists", queueId)
	}

	filePath, err := c.saveMessageToTempFile(queueId, data)
	if err != nil {
		return err
	}

	msg := Message{
		Data:           data,
		MessageId:      filepath.Base(filePath),
		PathToTempFile: filePath,
	}

	queue <- msg
	return nil
}

func (q *CheapQueueEngine) QueueLen(id string) int {
	q.mu.Lock()
	defer q.mu.Unlock()
	if ch, ok := q.queues[id]; ok {
		return len(ch)
	}
	return -1
}


func (c *CheapQueueEngine) Consume(queueId string, handler func([]byte)) error {
	c.mu.Lock()
	queue, exists := c.queues[queueId]
	c.mu.Unlock()

	if !exists {
		return fmt.Errorf("queue '%s' does not exists", queueId)
	}

	go func() {
		for msg := range queue {
			if handler != nil {
				handler(msg.Data)
			}
			os.Remove(msg.PathToTempFile)
		}
	}()
	return nil
}

func (c *CheapQueueEngine) ConsumeOnce(queueId string) ([]byte, error) {
	c.mu.Lock()
	queue, exists := c.queues[queueId]
	c.mu.Unlock()

	if !exists {
		return nil, fmt.Errorf("queue '%s' does not exists", queueId)
	}

	msg := <-queue

	defer os.Remove(msg.PathToTempFile)
	
	return msg.Data, nil
}

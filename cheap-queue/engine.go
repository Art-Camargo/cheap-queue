package cheapqueue

import "fmt"

func (c *CheapQueueEngine) Init() {
	c.queues = make(map[string]chan []byte)
}

func (c *CheapQueueEngine) Bind(queueId string, bufferSize int) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, exists := c.queues[queueId]; exists {
		return fmt.Errorf("queue '%s' already exists", queueId)
	}

	c.queues[queueId] = make(chan []byte, bufferSize)
	return nil
}


func (c *CheapQueueEngine) Publish(queueId string, data []byte) error {
	c.mu.Lock()
	queue, exists := c.queues[queueId]
	c.mu.Unlock()

	if !exists {
		return fmt.Errorf("queue '%s' does not exists", queueId)
	}

	queue <- data
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
				handler(msg)
			}
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
	return msg, nil
}

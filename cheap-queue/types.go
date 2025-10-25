package cheapqueue

import (
	"sync"
)

type Queue struct {
	Id   string
	Data []byte
}

type CheapQueueEngine struct {
	mu     sync.Mutex
	queues map[string]chan []byte
}

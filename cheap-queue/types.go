package cheapqueue

import (
	"sync"
)

type Message struct {
	Data           []byte
	MessageId      string
	PathToTempFile string
}

type Queue struct {
	Id   string
	Messages []Message
}

type CheapQueueEngine struct {
	mu        sync.Mutex
	queues    map[string]chan Message
	projectId string 
}

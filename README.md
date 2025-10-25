# 🚀 Cheap Queue

A lightweight, persistent, file-backed message queue for Go with automatic crash recovery.

## ✨ Features

- **Persistent**: Messages are saved to disk, surviving crashes and restarts
- **Automatic Recovery**: Reloads pending messages on startup
- **Multi-Project Support**: Isolated namespaces for different applications
- **Simple API**: Easy to use with minimal setup
- **Cross-Platform**: Works on Linux, macOS, and Windows

## 📦 Installation

```bash
go get github.com/Art-Camargo/cheap-queue
```

## 🎯 Quick Start

```go
package main

import (
    "fmt"
    cheapqueue "github.com/Art-Camargo/cheap-queue"
    "time"
)

func main() {
    // Initialize with a unique project ID
    var q cheapqueue.CheapQueueEngine
    q.Init("my-app-v1")

    // Create a queue with buffer size
    q.Bind("tasks", 100)

    // Publish messages
    q.Publish("tasks", []byte("Process this task"))

    // Consume messages
    q.Consume("tasks", func(data []byte) {
        fmt.Println("Processing:", string(data))
    })

    time.Sleep(1 * time.Second)
}
```

## 📖 API Reference

### `Init(projectId string)`

Initializes the queue engine with a unique project identifier. This ID is used to namespace temporary files.

**Important**: Must be called before any other operations.

### `Bind(queueId string, bufferSize int) error`

Creates or resizes a queue with the specified buffer size.

### `Publish(queueId string, data []byte) error`

Publishes a message to the queue. Message is persisted to disk before being enqueued.

### `Consume(queueId string, handler func([]byte)) error`

Starts consuming messages asynchronously. Handler is called for each message.

### `ConsumeOnce(queueId string) ([]byte, error)`

Consumes a single message synchronously.

### `QueueLen(queueId string) int`

Returns the number of messages currently in the queue.

## 🔄 Crash Recovery

Messages are automatically saved to temporary files. On restart:

1. Call `Init()` with the same `projectId`
2. Messages are automatically recovered
3. Call `Bind()` to recreate queues
4. Resume consumption

## 💡 Best Practices

- Use a unique, stable `projectId` for each application
- Set appropriate buffer sizes based on expected load
- Handle errors from `Publish()` and `Consume()`
- Clean shutdown: ensure all messages are consumed before exit

## 📝 License

MIT License - see LICENSE file for details

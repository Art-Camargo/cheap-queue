package cheapqueue

import "testing"

func TestQueueBasicOperations(t *testing.T) {
	var q CheapQueueEngine
	q.Init("test-project")

	// Test Bind
	err := q.Bind("test-queue", 10)
	if err != nil {
		t.Fatalf("Failed to bind queue: %v", err)
	}

	// Test Publish
	testData := []byte("test message")
	err = q.Publish("test-queue", testData)
	if err != nil {
		t.Fatalf("Failed to publish: %v", err)
	}

	// Test QueueLen
	length := q.QueueLen("test-queue")
	if length != 1 {
		t.Fatalf("Expected queue length 1, got %d", length)
	}

	// Test ConsumeOnce
	data, err := q.ConsumeOnce("test-queue")
	if err != nil {
		t.Fatalf("Failed to consume: %v", err)
	}

	if string(data) != string(testData) {
		t.Fatalf("Expected '%s', got '%s'", string(testData), string(data))
	}

	// Queue should be empty now
	length = q.QueueLen("test-queue")
	if length != 0 {
		t.Fatalf("Expected queue length 0, got %d", length)
	}
}

func TestQueueDoesNotExist(t *testing.T) {
	var q CheapQueueEngine
	q.Init("test-project-2")

	err := q.Publish("nonexistent", []byte("test"))
	if err == nil {
		t.Fatal("Expected error when publishing to non-existent queue")
	}
}

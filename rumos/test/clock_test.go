package test

import (
	"rumos/internal/server"
	"testing"

	"github.com/google/uuid"
)

func resetVectorClock() {
	system.VectorClock = make(map[string]int)
}

var system = server.GetInstance()

const numIDs = 1000

func TestEqual(t *testing.T) {
	defer resetVectorClock()

	TestClock := make(map[string]int)
	for i := 0; i < numIDs; i++ {
		id, _ := uuid.NewUUID()
		system.VectorClock[id.String()] = i
		TestClock[id.String()] = i
	}

	result := system.CompareClock(system.VectorClock, TestClock)
	if result != server.EQUAL {
		t.Errorf("Expected %d, got %d", server.EQUAL, result)
	}
}

func TestNewer(t *testing.T) {
	defer resetVectorClock()

	TestClock := make(map[string]int)
	for i := 0; i < numIDs; i++ {
		id, _ := uuid.NewUUID()
		system.VectorClock[id.String()] = i
		TestClock[id.String()] = i + 1
	}

	result := system.CompareClock(system.VectorClock, TestClock)
	if result != server.NEWER {
		t.Errorf("Expected %d, got %d", server.NEWER, result)
	}
}

func TestOlder(t *testing.T) {
	defer resetVectorClock()

	TestClock := make(map[string]int)
	for i := 0; i < numIDs; i++ {
		id, _ := uuid.NewUUID()
		system.VectorClock[id.String()] = i + 1
		TestClock[id.String()] = i
	}

	result := system.CompareClock(system.VectorClock, TestClock)
	if result != server.OLDER {
		t.Errorf("Expected %d, got %d", server.OLDER, result)
	}
}

func TestConcurrent(t *testing.T) {
	defer resetVectorClock()

	TestClock := make(map[string]int)
	for i := 0; i < numIDs; i++ {
		id, _ := uuid.NewUUID()
		if i%2 == 0 {
			system.VectorClock[id.String()] = i + 1
			TestClock[id.String()] = i
		} else {
			system.VectorClock[id.String()] = i
			TestClock[id.String()] = i + 1
		}
	}

	result := system.CompareClock(system.VectorClock, TestClock)
	if result != server.CONCURRENT {
		t.Errorf("Expected %d, got %d", server.CONCURRENT, result)
	}
}

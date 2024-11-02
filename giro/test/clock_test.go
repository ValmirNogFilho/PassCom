package test

import (
	"giro/internal/server"
	"testing"

	"github.com/google/uuid"
)

func resetVectorClock() {
	system.VectorClock = make(map[uuid.UUID]int)
}

var system = server.GetInstance()

const numIDs = 1000

func TestEqual(t *testing.T) {
	defer resetVectorClock()

	TestClock := make(map[uuid.UUID]int)
	for i := 0; i < numIDs; i++ {
		id, _ := uuid.NewUUID()
		system.VectorClock[id] = i
		TestClock[id] = i
	}

	result := system.CompareClock(system.VectorClock, TestClock)
	if result != server.EQUAL {
		t.Errorf("Expected %d, got %d", server.EQUAL, result)
	}
}

func TestNewer(t *testing.T) {
	defer resetVectorClock()

	TestClock := make(map[uuid.UUID]int)
	for i := 0; i < numIDs; i++ {
		id, _ := uuid.NewUUID()
		system.VectorClock[id] = i
		TestClock[id] = i + 1
	}

	result := system.CompareClock(system.VectorClock, TestClock)
	if result != server.NEWER {
		t.Errorf("Expected %d, got %d", server.NEWER, result)
	}
}

func TestOlder(t *testing.T) {
	defer resetVectorClock()

	TestClock := make(map[uuid.UUID]int)
	for i := 0; i < numIDs; i++ {
		id, _ := uuid.NewUUID()
		system.VectorClock[id] = i + 1
		TestClock[id] = i
	}

	result := system.CompareClock(system.VectorClock, TestClock)
	if result != server.OLDER {
		t.Errorf("Expected %d, got %d", server.OLDER, result)
	}
}

func TestConcurrent(t *testing.T) {
	defer resetVectorClock()

	TestClock := make(map[uuid.UUID]int)
	for i := 0; i < numIDs; i++ {
		id, _ := uuid.NewUUID()
		if i%2 == 0 {
			system.VectorClock[id] = i + 1
			TestClock[id] = i
		} else {
			system.VectorClock[id] = i
			TestClock[id] = i + 1
		}
	}

	result := system.CompareClock(system.VectorClock, TestClock)
	if result != server.CONCURRENT {
		t.Errorf("Expected %d, got %d", server.CONCURRENT, result)
	}
}

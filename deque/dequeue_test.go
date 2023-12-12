// Testing for dequeue.go
package deque

import (
	"sync"
	"sync/atomic"
	"testing"
)

func TestNewUnboundedDEQue(t *testing.T) {
	NewUnboundedDEQue[*int]()
}

func TestUnboundedDEQue_PushBottom(t *testing.T) {
	deque := NewUnboundedDEQue[int]()
	value1 := 1
	value2 := 2
	deque.PushBottom(&value1)
	deque.PushBottom(&value2)
	if deque.IsEmpty() {
		t.Errorf("Expected false, got %v", deque.IsEmpty())
	}
}

func TestUnboundedDEQue_IsEmpty(t *testing.T) {
	deque := NewUnboundedDEQue[int]()
	if !deque.IsEmpty() {
		t.Errorf("Expected true, got %v", deque.IsEmpty())
	}
}

func TestUnboundedDEQue_PopTop(t *testing.T) {
	deque := NewUnboundedDEQue[int]()
	value1 := 1
	value2 := 2
	deque.PushBottom(&value1)
	deque.PushBottom(&value2)
	deque.PopTop()
	if deque.IsEmpty() {
		t.Errorf("Expected false, got %v", deque.IsEmpty())
	}
}

func TestUnboundedDEQue_PopBottom(t *testing.T) {
	deque := NewUnboundedDEQue[int]()
	value1 := 1
	value2 := 2
	deque.PushBottom(&value1)
	deque.PushBottom(&value2)
	deque.PopBottom()
	if deque.IsEmpty() {
		t.Errorf("Expected false, got %v", deque.IsEmpty())
	}
}

func ParallelDequeTest(t *testing.T, deque DEQue[int64], numThreads int, numOps int) {

	var wg sync.WaitGroup
	var sum int64
	var parallel_sum atomic.Int64

	for thread := 0; thread < numThreads; thread++ {
		for i := 0; i < numOps/numThreads; i++ {
			i64 := int64(i)
			sum += i64
			deque.PushBottom(&i64)
		}
	}

	wg.Add(numThreads)
	for i := 0; i < numThreads; i++ {
		go func() {
			var value *int64
			defer wg.Done()
			for j := 0; j < numOps/numThreads; j++ {
				value = nil
				for value == nil {
					value = deque.PopTop()
				}
				parallel_sum.Add(*value)
			}
		}()
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < numOps/numThreads; i++ {
			i64 := int64(i)
			sum += i64
			deque.PushBottom(&i64)
			value := deque.PopBottom()
			parallel_sum.Add(*value)
		}
	}()

	wg.Wait()

	if parallel_sum.Load() != sum {
		t.Errorf("Expected %v, got %v", sum, parallel_sum.Load())
	} else {
		t.Logf("Expected %v, got %v", sum, parallel_sum.Load())
	}
}

func TestUnboundedDEQue_Parallel(t *testing.T) {
	deque := NewUnboundedDEQue[int64]()
	ParallelDequeTest(t, deque, 50, 1000000)
}

func ParallelContentionTest(t *testing.T, deque DEQue[int64], numThreads int, numOps int) {

	var wg sync.WaitGroup
	var sum int64
	var parallel_sum atomic.Int64

	wg.Add(numThreads)
	for i := 0; i < numThreads; i++ {
		go func() {
			var value *int64
			defer wg.Done()
			for j := 0; j < numOps/numThreads; j++ {
				value = nil
				for value == nil {
					value = deque.PopTop()
				}
				parallel_sum.Add(*value)
			}
		}()
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < numOps; i++ {
			i64 := int64(i)
			sum += i64
			deque.PushBottom(&i64)
			if i%5 == 0 {
				value := deque.PopBottom()
				if value != nil {
					parallel_sum.Add(*value)
					deque.PushBottom(&i64)
					sum += i64
				}
			}
		}
	}()

	wg.Wait()

	if parallel_sum.Load() != sum {
		t.Errorf("Expected %v, got %v", sum, parallel_sum.Load())
	} else {
		t.Logf("Expected %v, got %v", sum, parallel_sum.Load())
	}
}

func TestUnboundedDEQue_ParallelContention(t *testing.T) {
	deque := NewUnboundedDEQue[int64]()
	ParallelContentionTest(t, deque, 50, 10000000)
}

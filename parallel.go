package radixsort

import (
	"runtime"
	"sync"
)

// helpers to coordinate parallel sorts

type sortFunc func(interface{}, sortTask, func(sortTask))

// MaxProcs controls how many goroutines to start for large sorts. If 0,
// GOMAXPROCS will be used; if 1, all sorts will be serial.
var MaxProcs = 2

// MinParallel is the size of the smallest collection we will try to sort in
// parallel.
var MinParallel = 10000

// MinOffload is the size of the smallest range that can be offloaded to
// another goroutine.
var MinOffload = 127

// BufferRatio is how many sorting tasks to queue (buffer) up per
// worker goroutine.
var BufferRatio float32 = 1

// parallelSort calls the sorters with an asyncSort function that will hand
// the task off to another goroutine when possible.
func parallelSort(data interface{}, sorter sortFunc, initialTask sortTask) {
	max := runtime.GOMAXPROCS(0)
	if MaxProcs > 0 && MaxProcs < max {
		max = MaxProcs
	}
	l := data.(interface {
		Len() int
	}).Len()
	if l < MinParallel {
		max = 1
	}

	var syncSort func(t sortTask)
	syncSort = func(t sortTask) {
		sorter(data, t, syncSort)
	}
	if max == 1 {
		syncSort(initialTask)
		return
	}

	wg := new(sync.WaitGroup)
	// buffer up one extra task to keep each cpu busy
	sorts := make(chan sortTask, int(float32(max)*BufferRatio))
	var asyncSort func(t sortTask)
	asyncSort = func(t sortTask) {
		if t.end-t.pos < MinOffload {
			sorter(data, t, syncSort)
			return
		}
		wg.Add(1)
		select {
		case sorts <- t:
		default:
			sorter(data, t, asyncSort)
			wg.Done()
		}
	}
	doSortWork := func() {
		for task := range sorts {
			sorter(data, task, asyncSort)
			wg.Done()
		}
	}
	for i := 0; i < max; i++ {
		go doSortWork()
	}

	asyncSort(initialTask)

	wg.Wait()
	close(sorts)
}

package pool

import (
	"sync"
)

// WorkerPool представляет пул воркеров для выполнения задач.
type WorkerPool struct {
	taskQueue  chan func()
	workerPool int
	wg         sync.WaitGroup
	mu         sync.Mutex
	done       chan struct{}
}

// NewWorkerPool создает новый пул воркеров с заданным количеством воркеров.
func NewWorkerPool(workerCount int) *WorkerPool {
	wp := &WorkerPool{
		taskQueue:  make(chan func(), 100),
		workerPool: workerCount,
		done:       make(chan struct{}),
	}

	wp.StartWorkers(workerCount)
	return wp
}

// StartWorkers запускает указанное количество воркеров.
func (wp *WorkerPool) StartWorkers(workerCount int) {
	for i := 0; i < workerCount; i++ {
		go wp.worker()
	}
}

// SetWorkerCount обновляет количество воркеров.
func (wp *WorkerPool) SetWorkerCount(newWorkerCount int) {
	wp.mu.Lock()
	defer wp.mu.Unlock()

	close(wp.done)

	wp.done = make(chan struct{})

	wp.StartWorkers(newWorkerCount)
	wp.workerPool = newWorkerCount

}

// worker запускает задачи из очереди.
func (wp *WorkerPool) worker() {
	for {
		select {
		case task := <-wp.taskQueue:
			task()
			wp.wg.Done()
		case <-wp.done:
			return
		}
	}
}

// SubmitTask добавляет задачу в очередь
func (wp *WorkerPool) SubmitTask(task func()) {
	wp.wg.Add(1)
	wp.taskQueue <- task
}

// Wait завершает выполнение всех задач.
func (wp *WorkerPool) Wait() {
	wp.wg.Wait()
}

// Close завершает работу пула воркеров.
func (wp *WorkerPool) Close() {
	wp.mu.Lock()
	defer wp.mu.Unlock()
	close(wp.done)
	wp.Wait()
}

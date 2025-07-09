package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
)

func worker(ctx context.Context, id int, jobs <-chan int, wg *sync.WaitGroup, errChan chan<- error) {
	defer wg.Done()
	for {
		select {
		case job, ok := <-jobs:
			if !ok {
				return
			}

			if job%2 != 0 {
				errChan <- fmt.Errorf("Worker %d: ошибка при обработке %d", id, job)
				continue
			}

			fmt.Printf("Worker %d отработал задачу %d\n", id, job)
			time.Sleep(500 * time.Millisecond)
		case <-ctx.Done():
			fmt.Printf("Worker %d завершил работу из-за отмены\n", id)
			return
		}
	}
}

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	jobs := make(chan int)
	var wg sync.WaitGroup
	errChan := make(chan error)
	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go worker(ctx, i, jobs, &wg, errChan)
	}
	go errorsFinder(errChan)
	for i := 1; i <= 10; i++ {
		jobs <- i
	}
	time.Sleep(1 * time.Second)
	cancel()
	close(jobs)

	wg.Wait()
	close(errChan)
	fmt.Println("finish")
}

func errorsFinder(errChan <-chan error) {
	for err := range errChan {
		log.Println("Ошибка:", err)
	}
}

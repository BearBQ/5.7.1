package main

import (
	"fmt"
	"log"
	"sync"
)

func worker(id int, jobs <-chan int, wg *sync.WaitGroup, errChan chan<- error) {
	defer wg.Done()
	for {
		job, ok := <-jobs
		if !ok {
			return
		}

		if job%2 != 0 {
			errChan <- fmt.Errorf("Worker %d: ошибка при обработке %d", id, job)
			continue
		}

		fmt.Printf("Worker %d отработал задачу %d\n", id, job)
	}
}

func main() {
	jobs := make(chan int)
	var wg sync.WaitGroup
	errChan := make(chan error)
	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go worker(i, jobs, &wg, errChan)
	}
	go errorsFinder(errChan)
	for i := 1; i <= 10; i++ {
		jobs <- i
	}
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

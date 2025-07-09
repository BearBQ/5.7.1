package main

import (
	"fmt"
	"sync"
)

func worker(id int, jobs <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		val, ok := <-jobs
		if !ok {
			return
		}

		fmt.Printf("Worker %d отработал задачу %d\n", id, val)
	}
}

func main() {
	jobs := make(chan int, 10)
	var wg sync.WaitGroup

	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go worker(i, jobs, &wg)
	}
	for i := 1; i <= 10; i++ {
		jobs <- i
	}
	close(jobs)
	wg.Wait()
	fmt.Println("finish")
}

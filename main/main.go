// package main

// import "fmt"

// func partSum(a int, sum chan int, cur chan int, prev chan int) {
// 	cur <- a
// 	for i := <-prev; i < <-cur; i++ {
// 		sum <- <-sum + i
// 	}
// 	prev <- a
// }

// func sumFromRange(a int, b int) int {
// 	cur := make(chan int)
// 	sum := make(chan int)
// 	prev := make(chan int)
// 	prev <- a
// 	for i := a; i < b; i += 1 {
// 		go partSum(i, sum, cur, prev)
// 	}
// 	return <-sum
// }

// func rm(a int, ch chan int) {
// 	total := 0
// 	for i := 0; i < a; i++ {
// 		total += i
// 		fmt.Println(i)
// 	}
// 	ch <- total
// }

// func main() {
// 	fmt.Println(sumFromRange(0, 10))
// 	ch1 := make(chan int)
// 	ch2 := make(chan int)
// 	go rm(9, ch1)
// 	go rm(123, ch2)
// 	fmt.Println(<-ch1)
// }


package main

import (
	"fmt"
	"net/http"
)

type job struct {
	url  string
	resp chan *http.Response
}

type worker struct {
	jobs  chan *job
	count int
}

func (w *worker) getter(done chan *worker) {
	for {
		j := <-w.jobs
		resp, _ := http.Get(j.url)
		j.resp <- resp
		done <- w
	}
}

func get(jobs chan *job, url string, answer chan string) {
	resp := make(chan *http.Response)
	jobs <- &job{url, resp}
	r := <-resp
	answer <- r.Request.URL.String()
}

func balancer(count int, depth int) chan *job {
	jobs := make(chan *job)
	done := make(chan *worker)
	workers := make([]*worker, count)

	for i := 0; i < count; i++ {
		workers[i] = &worker{make(chan *job, depth), 0}
		go workers[i].getter(done)
	}

	go func() {
		for {
			var free *worker
			min := depth
			for _, w := range workers {
				if w.count < min {
					free = w
					min = w.count
				}
			}
			var jobsource chan *job
			if free != nil {
				jobsource = jobs
			}
			select {
			case j := <-jobsource:
				free.jobs <- j
				free.count++

			case w := <-done:
				w.count--
			}
		}

	}()

	return jobs
}

func main() {
	jobs := balancer(10, 10)
	answer := make(chan string)
	for {
		var url string
		if _, err := fmt.Scanln(&url); err != nil {
			break
		}
		go get(jobs, url, answer)
	}
	for u := range answer {
		fmt.Printf("%s\n", u)
	}
}

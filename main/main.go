package main

import "fmt"

func partSum(a int, sum chan int, cur chan int, prev chan int) {
	cur <- a
	for i := <-prev; i < <-cur; i++ {
		sum <- <-sum + i
	}
	prev <- a
}

func sumFromRange(a int, b int) int {
	cur := make(chan int)
	sum := make(chan int)
	prev := make(chan int)
	prev <- a
	for i := a; i < b; i += 1 {
		go partSum(i, sum, cur, prev)
	}
	return <-sum
}

func rm(a int, ch chan int) {
	total := 0
	for i := 0; i < a; i++ {
		total += i
		fmt.Println(i)
	}
	ch <- total
}

func main() {
	fmt.Println(sumFromRange(0, 10))
	ch1 := make(chan int)
	ch2 := make(chan int)
	go rm(9, ch1)
	go rm(123, ch2)
	fmt.Println(<-ch1)
}

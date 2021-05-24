package main

import "fmt"

func rm(a int,ch chan int){
	total:=0
	for i := 0; i < a; i++ {
		total+=i
		fmt.Println(i)
	}
	ch<-total
}

func main() {
	ch1:=make(chan int)
	ch2:=make(chan int)
	go rm(9,ch1)
	go rm(123,ch2)
	fmt.Println(<-ch1)
}

package main

import (
	"fmt"
	"math/rand"
	"time"
)

func sendToChannel(nums []int, ch chan<- int) {
	defer close(ch)

	for _, v := range nums {
		time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
		ch <- v
	}
}

func main() {
	ch1 := make(chan int)
	ch2 := make(chan int)
	result := make(chan int)

	go sendToChannel([]int{1, 2, 3, 4}, ch1)
	go sendToChannel([]int{5, 6, 7, 8}, ch2)

	go func() {
		defer close(result)
		for ch1 != nil || ch2 != nil {
			select {
			case v, ok := <-ch1:
				if !ok {
					ch1 = nil
					continue
				}
				result <- v
			case v, ok := <-ch2:
				if !ok {
					ch2 = nil
					continue
				}
				result <- v
			}
		}
	}()

	for v := range result {
		fmt.Println(v)
	}
}

package main

import (
	"fmt"
)

func main() {
	nums := []int{2, 4, 6, 8, 10}
	
	results := make(chan int, len(nums))

	for _, n := range nums {
		go func(val int) {
			square := val * val
			results <- square
		}(n)
	}

	for i := 0; i < len(nums); i++ {
		fmt.Printf("Result %d: %d\n", i+1, <-results)
	}
}
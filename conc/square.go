package main
import "fmt"
func main() {
	nums := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	results := make(chan int, len(nums))
	for _, i := range nums {
		go func(val int) {
			results <- val * val
		}(i)
	}
	for i := 0; i < len(nums); i++ {
		fmt.Printf("%d\n", <- results)
	}
}
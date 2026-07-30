package main

import (
	"fmt"
	"math/rand"
)

func main() {
	nums := make(chan int)
	results := make(chan int)

	go makeRandNums(nums)
	go powerNums(nums, results)

	for num := range results {
		fmt.Print(num, " ")
	}
	fmt.Println()

}

func makeRandNums(ch chan int) {
	nums := make([]int, 10)
	for i := 0; i < 10; i++ {
		nums[i] = rand.Intn(10)
	}
	fmt.Println(nums)
	for _, num := range nums {
		ch <- num
	}
	close(ch)
}

func powerNums(in chan int, out chan int) {
	for num := range in {
		out <- num * num
	}

	close(out)
}

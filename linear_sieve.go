package main

import (
	"fmt"
	"os"
	"strconv"
)

func LinearSieve(n int) []int {
	arr := make([]int, n+1)
	for i := 2; i <= n; i++ {
		arr[i] = i
	}
	prime := []int{}
	for i := 2; i <= n; i++ {
		if arr[i] == i {
			prime = append(prime, i)
		}
		for _, j := range prime {
			m := i * j
			if j <= arr[i] && m <= n {
				arr[m] = j
			} else {
				break
			}
		}
	}
	return prime
}

func main() {
	if len(os.Args) <= 1 {
		fmt.Println("No arguments.")
		return
	}
	num, _ := strconv.Atoi(os.Args[1])
	if num < 2 {
		return
	}
	fmt.Println(LinearSieve(num))
}

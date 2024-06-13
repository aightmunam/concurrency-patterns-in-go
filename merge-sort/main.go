package main

import (
	"fmt"
	"math/rand"
	"sync"
)

// Naive implementation of mergeSort where we try to parallelize
// everything without thinking is even slower than doing it
// sequentially, because the overhead of creating a goroutine
// for sorting very small arrays is too costly. For example:
// sorting an array of size 1 or 2, we create a goroutine, which will
// be much more costly than doing it sequentially.

// A better approach is to define a threshold and if the array size
// becomes smaller than that threshold, we do not use a goroutine to
// sort it but do it sequentially.

// An optimal threshold is 2048.
const min_size_threshold = 2048

func mergeSort(items []int) []int {
	low := 0
	high := len(items)
	if high == 1 {
		return items
	}
	middle := len(items) / 2

	var items_left, items_right []int
	if high < min_size_threshold {
		items_left = mergeSort(items[low:middle])
		items_right = mergeSort(items[middle:high])
	} else {
		var wg sync.WaitGroup
		wg.Add(2)
		go func() {
			defer wg.Done()
			items_left = mergeSort(items[low:middle])
		}()

		go func() {
			defer wg.Done()
			items_right = mergeSort(items[middle:high])
		}()
		wg.Wait()
	}

	return merge(items_left, items_right)
}

func merge(items_left, items_right []int) []int {
	var i, j, k int
	total_left := len(items_left)
	total_right := len(items_right)

	output := make([]int, total_right+total_left)
	for i < total_left && j < total_right {
		if items_left[i] < items_right[j] {
			output[k] = items_left[i]
			i++
		} else {
			output[k] = items_right[j]
			j++
		}
		k++
	}

	for ; i < total_left; i++ {
		output[k] = items_left[i]
		k++
	}

	for ; j < total_right; j++ {
		output[k] = items_right[j]
		k++
	}
	return output
}

func main() {
	const array_size = 15000
	unsorted := make([]int, array_size)

	for i := 0; i < array_size; i++ {
		unsorted[i] = rand.Intn(array_size)
	}
	sorted := mergeSort(unsorted)
	fmt.Println(unsorted)
	fmt.Println(sorted)
}

package main

import "fmt"

func mergeSort(items []int) []int {
	low := 0
	high := len(items)
	if high == 1 {
		return items
	}
	middle := len(items) / 2

	items_left := mergeSort(items[low:middle])
	items_right := mergeSort(items[middle:high])

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
	unsorted := []int{2, 1, -4, 5, -12, 9, 6, 0, 2, 12, 4, -3}
	sorted := mergeSort(unsorted)
	fmt.Println(unsorted)
	fmt.Println(sorted)

	// fmt.Println(
	// 	merge([]int{-12, -4, 1, 2, 5, 9}, []int{-3, 0, 2, 4, 6, 12}),
	// )
}

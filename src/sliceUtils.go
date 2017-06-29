package main

/*
* This file contains utilitary functions that are reusable for slices
 */

/*
* Functions
 */

func intersection(l1 []int, l2 []int) []int {
	var result []int
	tmpMap := make(map[int]int)

	for _, value := range l1 {
		tmpMap[value]++
	}

	for _, value := range l2 {
		tmpMap[value]++
	}

	for _, v := range tmpMap {
		if v == 2 {
			result = append(result, v)
		}
	}

	return result
}

func union(l1 []int, l2 []int) []int {
	return append(l1, l2...)
}

func contains(x int, list []int) bool {
	for _, value := range list {
		if x == value {
			return true
		}
	}
	return false
}

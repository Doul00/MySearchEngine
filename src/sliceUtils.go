package main

/*
* This file contains utilitary functions that are reusable for slices
 */

/*
* Functions
 */

/*
* @l1 first slice
* @l2 second slice
* Returns the intersection of the two slices
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

	for key, v := range tmpMap {
		if v == 2 {
			result = append(result, key)
		}
	}

	return result
}

/*
* @l1 first slice
* @l2 second slice
* Returns the union of the two slices, without duplicates
 */
func union(l1 []int, l2 []int) []int {
	var result []int
	tmpMap := make(map[int]int)

	for _, value := range l1 {
		tmpMap[value]++
	}

	for _, value := range l2 {
		tmpMap[value]++
	}

	for key := range tmpMap {
		result = append(result, key)
	}

	return result
}

/*
* @x the value to find
* @list the list of values
* Returns true if the value was found
 */
func contains(x int, list []int) bool {
	for _, value := range list {
		if x == value {
			return true
		}
	}
	return false
}

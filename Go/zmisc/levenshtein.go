package main

import "fmt"

func Levenshtein(a, b string) int {
	la := len(a)
	lb := len(b)

	// Create a 2D slice to hold the distances
	distance := make([][]int, la+1)
	for i := range distance {
		distance[i] = make([]int, lb+1)
	}

	// Initialize the distance matrix
	for i := 0; i <= la; i++ {
		distance[i][0] = i
	}
	for j := 0; j <= lb; j++ {
		distance[0][j] = j
	}

	// Calculate the distances
	for i := 1; i <= la; i++ {
		for j := 1; j <= lb; j++ {
			cost := 0
			if a[i-1] != b[j-1] {
				cost = 1
			}
			distance[i][j] = min(
				distance[i-1][j]+1,      // Deletion
				distance[i][j-1]+1,      // Insertion
				distance[i-1][j-1]+cost, // Substitution
			)
		}
	}

	return distance[la][lb]
}

func main() {
	got := "{1:2,2:3,3:4}[okay]"
	want := "{1:2,2:3,3:4}[okay]"
	fmt.Println(Levenshtein(got, want))
}

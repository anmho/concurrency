package main

import (
	"testing"
)

// func Test_DoWork_Nondeterministic(t *testing.T) {

// 	tests := []struct {
// 		name     string
// 		intSlice []int
// 	}{
// 		{
// 			name: "happy path: genearates all numbers",
// 			intSlice: []int{
// 				0, 1, 2, 3, 5,
// 			},
// 		},
// 	}

// 	// check all values in the output chan are equal to the
// 	for _, tc := range tests {
// 		t.Run(tc.name, func(t *testing.T) {
// 			done := make(chan interface{})
// 			_, results := DoWork(done, tc.intSlice...)
// 			for i, expectedInt := range tc.intSlice {
// 				select {
// 				case num := <-results:
// 					if expectedInt != num {
// 						t.Errorf("index %d: expected %d but got %d", i, expectedInt, num)
// 					}

// 				case <-time.After(1 * time.Second):
// 					t.Fatal("test timed out")
// 				}
// 			}
// 		})
// 	}
// }

func Test_DoWork_Deterministic(t *testing.T) {

	done := make(chan interface{})
	defer close(done)

	intSlice := []int{0, 1, 2, 3, 5}

	heartbeat, results := DoWork(done, intSlice...)

	i := 0
	for r := range results {
		expected := intSlice[i]

		<-heartbeat

		// wait for heartbeat before testing
		if r != expected {
			t.Errorf("index %v: expected %v, but received %v,", i, expected, r)
		}
		i++
	}

}

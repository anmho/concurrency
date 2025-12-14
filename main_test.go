package main

import (
	"testing"
	"time"
)

func Test_DoWork(t *testing.T) {

	tests := []struct {
		name     string
		intSlice []int
	}{
		{
			name: "happy path",
			intSlice: []int{
				0, 1, 2, 3, 5,
			},
		},
	}

	// check all values in the output chan are equal to the
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			done := make(chan interface{})
			_, results := DoWork(done, tc.intSlice...)
			for i, expectedInt := range tc.intSlice {
				select {
				case num := <-results:
					if expectedInt != num {
						t.Errorf("index %d: expected %d but got %d", i, expectedInt, num)
					}

				case <-time.After(1 * time.Second):
					t.Fatal("test timed out")
				}
			}
		})
	}
}

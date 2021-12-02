package main

import (
	"fmt"
	"testing"

	"uk.ac.bris.cs/gameoflife/gol"
)

// TestGol tests 16x16, 64x64 and 512x512 images on 0, 1 and 100 turns using 1-16 worker threads.
func Benchmark(b *testing.B) {
	tests := []gol.Params{
		{ImageWidth: 512, ImageHeight: 512},
	}
	for _, p := range tests {
		for _, turns := range []int{100} {
			p.Turns = turns
			for threads := 1; threads <= 16; threads++ {
				p.Threads = threads
				testName := fmt.Sprintf("%dx%dx%d-%d", p.ImageWidth, p.ImageHeight, p.Turns, p.Threads)
				b.Run(testName, func(b *testing.B) {
					events := make(chan gol.Event)
					go gol.Run(p, events, nil)
				LOOP:
					for event := range events {
						switch e := event.(type) {
						case gol.FinalTurnComplete:
							fmt.Println(e)
							break LOOP
						}
					}
				})
			}
		}
	}
}
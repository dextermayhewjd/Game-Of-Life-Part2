package main

import (
	//"errors"
	"flag"
	"fmt"
	"math/rand"
	"net"
	"net/rpc"
	"os"
	"time"
	"uk.ac.bris.cs/gameoflife/stubs"
)

type Params struct {
	Turns       int
	Threads     int
	ImageWidth  int
	ImageHeight int
}

func makeEmptyWorld(width int, height int) [][]uint8 {
	result := make([][]uint8, height)
	for i := range result {
		result[i] = make([]uint8, width)
	}
	return result
}

func worker(StartX int, StartY int, EndX int, EndY int, world [][]uint8, out chan<- [][]uint8) { // use this worker in Game of Life func

	result := calculateNextState(StartX, StartY, EndX, EndY, world)
	out <- result
}
func calculateNextState(StartX int, StartY int, EndX int, EndY int, world [][]uint8) [][]uint8 { //could be run correct at thread 1

	width := EndX - StartX
	height := EndY - StartY

	nextWorld := makeEmptyWorld(width, height)

	for h := StartY; h < EndY; h++ {
		for w := StartX; w < EndX; w++ {
			nei := CheckCellAround(world, h, w)
			if nei > 3 || nei < 2 {
				nextWorld[h-StartY][w-StartX] = 0
				//if world[h][w] != 0 {
				//	reportChan <- gol.CellFlipped{CompletedTurns: turn,
				//		Cell: util.Cell{X: w, Y: h}}
				//}
			} else if nei == 3 {
				nextWorld[h-StartY][w-StartX] = 255
				//if world[h][w] != 255 {
				//	reportChan <- gol.CellFlipped{CompletedTurns: turn,
				//		Cell: util.Cell{X: w, Y: h}}
				//}

			} else if world[h][w] == 255 {
				nextWorld[h-StartY][w-StartX] = 255
			}
		}
	}
	return nextWorld
}

//func CheckCellAround(world [][]uint8, x int, y int) int {
//	sum := 0
//	for i := -1; i < 2; i++ {
//		for j := -1; j < 2; j++ {
//			if i == 0 && j == 0 { //ignore the centre cell
//				continue
//			}
//			x1 := (x + i + len(world)) % len(world)
//			y1 := (y + j + len(world[0])) % len(world[0])
//
//			if world[x1][y1] == 255 {
//				sum += 1
//			}
//		}
//	}
//	return sum
//}

func CheckCellAround(world [][]uint8, x int, y int) int {
	sum := 0
	for i := -1; i < 2; i++ {
		for j := -1; j < 2; j++ {
			if i == 0 && j == 0 { //ignore the centre cell
				continue
			}
			x1 := x + i
			y1 := y + j
			if x+i < 0 {
				x1 = len(world) - 1
			}
			if x+i >= len(world) {
				x1 = x + i - len(world)
			}
			if y+j < 0 {
				y1 = len(world[x1]) - 1
			}
			if y+j >= len(world[x1]) {
				y1 = y + j - len(world[x1])
			}
			if world[x1][y1] == 255 {
				sum += 1
			}
		}
	}
	return sum
}

func calculateAliveCells(world [][]uint8) int {
	var numberOfLivingCell = 0
	for wn, width := range world {
		for hn := range width {
			if world[wn][hn] == 255 {
				numberOfLivingCell += 1
			}
		}
	}

	return numberOfLivingCell
}

func aliveCellsCountTicker(world [][]uint8) int {
	ticker := time.NewTicker(2 * time.Second)

	for {
		select {
		case <-ticker.C:
			calculateAliveCells(world)
		}
	}
}

func gameOfLife(p Params, world [][]uint8) [][]uint8 {
	//turnCompleted1 = 0
	finalWorld := world
	thread := p.Threads
	width := p.ImageWidth
	height := p.ImageHeight


	if thread == 1 {

		//for i := 0; i < turn; i++ {
		finalWorld = calculateNextState(0, 0, width, height, finalWorld)
		//}
	} else {
		out := make([]chan [][]uint8, p.Threads)

		for j := 0; j < thread; j++ {
			out[j] = make(chan [][]uint8)
		}

		//for i := 0; i < turn; i++ {
		var partFinalWorld [][]uint8

		//fmt.Println("mutex locked")

		for k := 0; k < thread; k++ {

			go worker(0, height*k/thread, width, height*(k+1)/thread, finalWorld, out[k]) //worker(StartX int, StartY int, EndX int, EndY int ,world [][]uint8,out chan<- [][]uint8)

		}

		for l := 0; l < thread; l++ {

			partFinalWorld = append(partFinalWorld, <-out[l]...)
		}

		finalWorld = partFinalWorld
		//worldChan <- finalWorld
		//completedTurn <- i + 1

		//fmt.Println("mutex unlocked")

		//}

	}
	return finalWorld
}

func shutDown() {
	os.Exit(10)
}

type ServerOperations struct{}

func (s *ServerOperations) ShutDown(req stubs.Kill, res *stubs.Response) (err error) {
	fmt.Println("ShutDown Operation processing")
	if req.DeathMessage == "shutdown" {
		fmt.Println("shutting down")
		shutDown()
	}
	return
}

type GameOfLifeOperations struct{}

func (s *GameOfLifeOperations) GameOfLife(req stubs.Request, res *stubs.Response) (err error) {
	fmt.Println("Request received")
	p := Params{
		Threads:     req.Threads,
		ImageWidth:  req.ImageWidth,
		ImageHeight: req.ImageHeight,
	}

	res.World = gameOfLife(p, req.CurrentWorld)
	return
}

func (s *GameOfLifeOperations) DistributedWorld(req stubs.Request2, res *stubs.Response2) (err error) {
	fmt.Println("DistributedWorld receives the request")
	res.PartWorld = calculateNextState(req.StartX,req.StartY,req.EndX,req.EndY,req.CurrentWorld)
	fmt.Println("DistributedWorld finished")
	return
}

func main() {
	pAddr := flag.String("port", "8030", "Port to listen on")
	flag.Parse()
	rand.Seed(time.Now().UnixNano())
	rpc.Register(&GameOfLifeOperations{})
	rpc.Register(&ServerOperations{})
	listener, _ := net.Listen("tcp", ":"+*pAddr)
	defer listener.Close()
	rpc.Accept(listener)
}

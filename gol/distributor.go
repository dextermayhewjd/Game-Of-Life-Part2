package gol

import (
	"fmt"
	"net/rpc"
	"strconv"
	"sync"
	"time"
	"uk.ac.bris.cs/gameoflife/stubs"
	"uk.ac.bris.cs/gameoflife/util"
	//"uk.ac.bris.cs/gameoflife/util"
)

type distributorChannels struct {
	events     chan<- Event
	ioCommand  chan<- ioCommand
	ioIdle     <-chan bool
	ioFilename chan<- string
	ioOutput   chan<- uint8
	ioInput    <-chan uint8
}

func makeEmptyWorld(width int, height int) [][]uint8 {
	result := make([][]uint8, height)
	for i := range result {
		result[i] = make([]uint8, width)
	}
	return result
}

//func gameOfLife(p Params,world [][]uint8, turn int, mutex *sync.Mutex)  [][]uint8{
//func gameOfLife(p Params, world [][]uint8, turn int, completedTurn chan int, reportChan chan<- Event, worldChan chan [][]uint8) [][]uint8 {
//	//turnCompleted1 = 0
//	finalWorld := world
//	thread := p.Threads
//	width := p.ImageWidth
//	height := p.ImageHeight
//	if turn == 0 {
//		return world
//	}
//
//	if thread == 1 {
//
//		for i := 0; i < turn; i++ {
//			finalWorld = calculateNextState(0, 0, width, height, i+1, finalWorld, reportChan)
//		}
//
//	} else {
//		out := make([]chan [][]uint8, p.Threads)
//
//		for j := 0; j < thread; j++ {
//			out[j] = make(chan [][]uint8)
//		}
//
//		for i := 0; i < turn; i++ {
//			//semaphore.Wait()
//			var partFinalWorld [][]uint8
//
//			//fmt.Println("mutex locked")
//
//			for k := 0; k < thread; k++ {
//
//				go worker(0, height*k/thread, width, height*(k+1)/thread, i+1, finalWorld, out[k], reportChan) //worker(StartX int, StartY int, EndX int, EndY int ,world [][]uint8,out chan<- [][]uint8)
//
//			}
//
//			for l := 0; l < thread; l++ {
//
//				partFinalWorld = append(partFinalWorld, <-out[l]...)
//			}
//
//			finalWorld = partFinalWorld
//			worldChan <- finalWorld
//			completedTurn <- i + 1
//
//			//fmt.Println("mutex unlocked")
//
//			//semaphore.Post()
//
//		}
//
//	}
//	return finalWorld
//}

//func worker(StartX int, StartY int, EndX int, EndY int, turn int, world [][]uint8, out chan<- [][]uint8, reportChan chan<- Event) { // use this worker in Game of Life func
//
//	result := calculateNextState(StartX, StartY, EndX, EndY, turn, world, reportChan)
//	out <- result
//}
//todo to make cellflipped event isolated from the calcaulatenext stage
//func calculateNextState(StartX int, StartY int, EndX int, EndY int, turn int, world [][]uint8, reportChan chan<- Event) [][]uint8 { //could be run correct at thread 1
//
//	width := EndX - StartX
//	height := EndY - StartY
//
//	nextWorld := makeEmptyWorld(width, height)
//
//	for h := StartY; h < EndY; h++ {
//		for w := StartX; w < EndX; w++ {
//			nei := CheckCellAround(world, h, w)
//			if nei > 3 || nei < 2 {
//				nextWorld[h-StartY][w-StartX] = 0
//				if world[h][w] != 0 {
//					reportChan <- CellFlipped{CompletedTurns: turn,
//						Cell: util.Cell{X: w, Y: h}}
//				}
//			} else if nei == 3 {
//				nextWorld[h-StartY][w-StartX] = 255
//				if world[h][w] != 255 {
//					reportChan <- CellFlipped{CompletedTurns: turn,
//						Cell: util.Cell{X: w, Y: h}}
//				}
//
//			} else if world[h][w] == 255 {
//				nextWorld[h-StartY][w-StartX] = 255
//			}
//		}
//	}
//	return nextWorld
//}

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
//
func calculateAliveCells(world [][]uint8) []util.Cell {
	var currentLivingCell []util.Cell
	for wn, width := range world {
		for hn := range width {
			if world[wn][hn] == 255 {
				cell1 := util.Cell{X: hn, Y: wn}
				currentLivingCell = append(currentLivingCell, cell1)
			}
		}
	}

	return currentLivingCell
}

func countCell(world [][]uint8) int {
	sum := 0
	for h, h1 := range world {
		for w := range h1 {
			if world[h][w] == 255 {
				sum++
			}
		}
	}
	return sum
}

//
//func reportTimer(c distributorChannels, timer *time.Ticker, reportTurn chan int, report chan<- Event, isClose chan bool, worldChan chan [][]uint8, keyPress <-chan rune) {
//	var CompletedTurn = 0
//	var finalWorld [][]uint8
//	var num = 0
//
//	for {
//		select {
//		case <-timer.C:
//			//fmt.Println("case <-timer.C:")
//
//			report <- AliveCellsCount{
//				CompletedTurns: CompletedTurn,
//				CellsCount:     num,
//			}
//
//			// mutex lock  at gamoflife
//		case CompletedTurn = <-reportTurn:
//
//			report <- TurnComplete{CompletedTurns: CompletedTurn}
//
//		case finalWorld = <-worldChan:
//
//			num = countCell(finalWorld)
//		//mutex unlock at gameoflife
//
//		case handle := <-keyPress:
//			filename := strconv.Itoa(len(finalWorld[0])) + "x" + strconv.Itoa(len(finalWorld)) + "x turn:" + strconv.Itoa(CompletedTurn) + "create by reportTimer"
//			switch handle {
//			case 's': // make current graph
//				makeGraph(c, filename, CompletedTurn, finalWorld)
//			case 'q': // make current graph and close it
//				makeGraph(c, filename, CompletedTurn, finalWorld)
//				c.ioCommand <- ioCheckIdle
//				<-c.ioIdle
//
//				c.events <- StateChange{CompletedTurn, Quitting}
//
//				close(c.events)
//				isClose <- true
//				os.Exit(0)
//
//			case 'p': // pause the graph
//				c.events <- StateChange{CompletedTurns: CompletedTurn, NewState: Paused}
//				for 'p' != <-keyPress {
//					handle = <-keyPress
//					switch handle {
//					case 's':
//						makeGraph(c, filename, CompletedTurn, finalWorld)
//					case 'q':
//						makeGraph(c, filename, CompletedTurn, finalWorld)
//						c.ioCommand <- ioCheckIdle
//						<-c.ioIdle
//
//						c.events <- StateChange{CompletedTurn, Quitting}
//
//						close(c.events)
//						isClose <- true
//						os.Exit(0)
//
//					}
//
//				}
//				c.events <- StateChange{CompletedTurns: CompletedTurn, NewState: Executing}
//
//			}
//
//		case <-isClose:
//			return
//
//		}
//	}
//
//}

func DistributeCall(client *rpc.Client, StartX, StartY, EndX, EndY int, currentWorld [][]uint8) *stubs.Response2 {
	request := stubs.Request2{
		StartX:       StartX,
		StartY:       StartY,
		EndX:         EndX,
		EndY:         EndY,
		CurrentWorld: currentWorld,
	}
	response := new(stubs.Response2)
	client.Call(stubs.DWHandler, request, response)
	return response
}

// distributor divides the work between workers and interacts with other goroutines.

func makeCall(client *rpc.Client, threads, imageWidth, imageHeight int, currentWorld [][]uint8) *stubs.Response {
	request := stubs.Request{
		Threads:      threads,
		ImageWidth:   imageWidth,
		ImageHeight:  imageHeight,
		CurrentWorld: currentWorld,
	}
	response := new(stubs.Response)
	client.Call(stubs.GOLHandler, request, response)
	return response
}

func goServerToShutDown(client *rpc.Client) {
	shutDownRequest := stubs.Kill{DeathMessage: "shutdown"}
	response := new(stubs.Response)
	client.Go(stubs.QuitHandler, shutDownRequest, response, nil)
	return
}

func reportTicker(done chan bool, report chan<- Event, world *[][]uint8, mutex *sync.Mutex, turn *int) {
	ticker := time.NewTicker(2 * time.Second)
	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			mutex.Lock()
			numOfCurrentLivingCell := countCell(*world)
			report <- AliveCellsCount{
				CompletedTurns: *turn,
				CellsCount:     numOfCurrentLivingCell,
			}
			mutex.Unlock()
		}
	}
}
func distributor(p Params, c distributorChannels, keyPresses <-chan rune) {
	filename := strconv.Itoa(p.ImageHeight) + "x" + strconv.Itoa(p.ImageWidth)
	c.ioCommand <- ioInput
	c.ioFilename <- filename
	eventChan := c.events
	currentWorld := makeEmptyWorld(p.ImageWidth, p.ImageWidth)

	done := make(chan bool)
	var turnCompleted1 = 0

	for i := range currentWorld {
		for j := range currentWorld[i] {
			currentWorld[i][j] = <-c.ioInput
			if currentWorld[i][j] == 255 {
				eventChan <- CellFlipped{
					CompletedTurns: 0,
					Cell:           util.Cell{X: j, Y: i},
				}
			}
		}

	}
	mutex := &sync.Mutex{}

	serverList := []string{
		"127.0.0.1:8040",
		"127.0.0.1:8030",
		"127.0.0.1:8050",
		"127.0.0.1:8060",
	}
	serverNum := len(serverList)
	//serverNum := 1
	//server := flag.String("serve
	//r","127.0.0.1:8030","IP:port string to connect to as server")
	//flag.Parse()

	clientList := make([]*rpc.Client, serverNum)
	for i := range serverList {
		clientList[i], _ = rpc.Dial("tcp", serverList[i])
	}

	//completedTurn := make(chan int)
	//wordChan := make(chan [][]uint8)
	go reportTicker(done, eventChan, &currentWorld, mutex, &turnCompleted1)
	go keyPress(c, &turnCompleted1, &currentWorld, keyPresses, clientList, mutex)
	if serverNum == 1 {

		for turnCompleted1 < p.Turns {
			mutex.Lock()
			previousWorld := currentWorld
			currentWorld = makeCall(clientList[0], p.Threads, p.ImageWidth, p.ImageHeight, currentWorld).World
			turnCompleted1 += 1
			compareTwoWorld(eventChan, previousWorld, currentWorld, turnCompleted1)
			mutex.Unlock()
		}
	} else {
		outPartWord := make([][][]uint8, serverNum)
		//out := make([]chan [][]uint8,serverNum)
		//for i:=0 ; i < serverNum ; i++ {
		//	out[i] = make(chan [][]uint8)
		//}
		Extre := p.ImageHeight % p.Threads
		for j := 0; j < p.Turns; j++ {
			mutex.Lock()
			var partFinalWorld [][]uint8
			for k := 0; k < serverNum; k++ {
				if Extre != 0 && k == serverNum-1 {
					outPartWord[k] = DistributeCall(clientList[k], 0, p.ImageHeight/serverNum*k, p.ImageWidth, p.ImageHeight/serverNum*(k+1)+Extre, currentWorld).PartWorld
				} else {
					outPartWord[k] = DistributeCall(clientList[k], 0, p.ImageHeight/serverNum*k, p.ImageWidth, p.ImageHeight/serverNum*(k+1), currentWorld).PartWorld
				}
			}

			for i := 0; i < serverNum; i++ {
				partFinalWorld = append(partFinalWorld, outPartWord[i]...)
			}
			compareTwoWorld(eventChan, currentWorld, partFinalWorld, turnCompleted1)
			currentWorld = partFinalWorld
			turnCompleted1 += 1
			mutex.Unlock()

		}

	}
	done <- true
	//isClose := make(chan bool)

	//timer := time.NewTicker(2 * time.Second)

	//go reportTimer(c, timer, completedTurn, eventChan, isClose, wordChan, keyPresses)

	//finalWorld := gameOfLife(p, initWorld, p.Turns, completedTurn, eventChan, wordChan)

	existingCells := calculateAliveCells(currentWorld)

	eventChan <- FinalTurnComplete{
		CompletedTurns: p.Turns,
		Alive:          existingCells,
	}
	filename = strconv.Itoa(p.ImageWidth) + "x" + strconv.Itoa(p.ImageHeight) + "x" + strconv.Itoa(p.Turns)
	makeGraph(c, filename, p.Turns, currentWorld)

	// Make sure that the Io has finished any output before exiting.
	c.ioCommand <- ioCheckIdle
	<-c.ioIdle

	eventChan <- StateChange{p.Turns, Quitting}

	// Close the channel to stop the SDL goroutine gracefully. Removing may cause deadlock.

	//isClose <- true
	//todo defer client.close()
	close(eventChan)

}

//
func makeGraph(c distributorChannels, filename string, turn int, finalWorld [][]uint8) {
	c.ioCommand <- ioOutput
	c.ioFilename <- filename
	for h, h1 := range finalWorld {
		for w := range h1 {
			c.ioOutput <- finalWorld[h][w]
		}
	}
	c.events <- ImageOutputComplete{CompletedTurns: turn, Filename: filename}
}
func keyPress(c distributorChannels, completedTurn *int, currentWorld *[][]uint8, keyPress <-chan rune, clientList []*rpc.Client, mutex *sync.Mutex) {
	filename := "test file name"

	for {
		select {

		//filename := strconv.Itoa(len(currentWorld[0])) + "x" + strconv.Itoa(len(*currentWorld)) //+ "x turn:" + strconv.Itoa(CompletedTurn)//

		case handle := <-keyPress:
			switch handle {
			case 's':
				makeGraph(c, filename, *completedTurn, *currentWorld)
			case 'q':
				c.ioCommand <- ioCheckIdle
				<-c.ioIdle

				c.events <- StateChange{*completedTurn, Quitting}

				c.events <- FinalTurnComplete{CompletedTurns: *completedTurn}

			case 'k':
				makeGraph(c, filename, *completedTurn, *currentWorld)
				for i := range clientList {
					goServerToShutDown(clientList[i])
				}
				c.events <- StateChange{*completedTurn, Quitting}

				c.events <- FinalTurnComplete{CompletedTurns: *completedTurn}

			case 'p': // pause the sdl
				c.events <- StateChange{NewState: Paused, CompletedTurns: *completedTurn}
				mutex.Lock()

				for 'p' != <-keyPress {

				}
				c.events <- StateChange{NewState: Executing, CompletedTurns: *completedTurn}
				fmt.Println("Continuing")

				mutex.Unlock()

			}

		}
	}

}
func compareTwoWorld(eventChan chan<- Event, previousWorld [][]uint8, currentWorld [][]uint8, turn int) {
	for h, h1 := range previousWorld {
		for w := range h1 {
			if previousWorld[h][w] != currentWorld[h][w] {
				eventChan <- CellFlipped{CompletedTurns: turn, Cell: util.Cell{X: w, Y: h}}
			}
		}
	}
	eventChan <- TurnComplete{CompletedTurns: turn}
}
package gol

import (
	"net/rpc"
	"strconv"
	"uk.ac.bris.cs/gameoflife/stubs"
	"uk.ac.bris.cs/gameoflife/util"
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
//
//func countCell(world [][]uint8) int {
//	sum := 0
//	for h, h1 := range world {
//		for w := range h1 {
//			if world[h][w] == 255 {
//				sum++
//			}
//		}
//	}
//	return sum
//}
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

// distributor divides the work between workers and interacts with other goroutines.

func makeCall(
	client *rpc.Client,
	turns int, //total turns
	threads       int,
	imageWidth    int,
	imageHeight   int,
	currentWorld  [][]uint8,
	turn          int, //target turns
	completedTurn chan int,
	worldChan     chan [][]uint8) *stubs.Response {
	request := stubs.Request{
		Turns:         turns,
		Threads:       threads,
		ImageWidth:    imageWidth,
		ImageHeight:   imageHeight,
		CurrentWorld:  currentWorld,
		Turn:          turn,
		CompletedTurn: completedTurn,
		WorldChan:     worldChan,
	}
	response := new(stubs.Response)
	client.Call(stubs.GOLHandler,request,response)
	return response
}

func distributor(p Params, c distributorChannels, keyPresses <-chan rune) {
	filename := strconv.Itoa(p.ImageHeight) + "x" + strconv.Itoa(p.ImageWidth)
	c.ioCommand <- ioInput
	c.ioFilename <- filename
	eventChan := c.events

	initWorld := makeEmptyWorld(p.ImageWidth, p.ImageWidth)

	for i := range initWorld {
		for j := range initWorld[i] {
			initWorld[i][j] = <-c.ioInput
			//if initWorld[i][j] == 255 {
			//	eventChan <- CellFlipped{CompletedTurns: 0,
			//		Cell: util.Cell{X: j, Y: i}}
			}

		}


	server :="127.0.0.1:8030"
	//server := flag.String("server","127.0.0.1:8030","IP:port string to connect to as server")
	//flag.Parse()
	client, _ := rpc.Dial("tcp", server)
	defer client.Close()
	completedTurn := make(chan int)
	wordChan := make(chan [][]uint8)
	finalWorld := makeCall(client,p.Turns,p.Threads,p.ImageWidth,p.ImageHeight,initWorld,p.Turns,completedTurn,wordChan).World
	//isClose := make(chan bool)


	//timer := time.NewTicker(2 * time.Second)

	//go reportTimer(c, timer, completedTurn, eventChan, isClose, wordChan, keyPresses)

	//finalWorld := gameOfLife(p, initWorld, p.Turns, completedTurn, eventChan, wordChan)

	existingCells := calculateAliveCells(finalWorld)


	eventChan <- FinalTurnComplete{
		CompletedTurns: p.Turns,
		Alive:          existingCells,
	}
	//filename = strconv.Itoa(p.ImageWidth) + "x" + strconv.Itoa(p.ImageHeight) + "x" + strconv.Itoa(p.Turns)
	//makeGraph(c, filename, p.Turns, finalWorld)

	// Make sure that the Io has finished any output before exiting.
	c.ioCommand <- ioCheckIdle
	<-c.ioIdle


	eventChan <- StateChange{p.Turns, Quitting}


	// Close the channel to stop the SDL goroutine gracefully. Removing may cause deadlock.

	//isClose <- true

	close(eventChan)

}
//
//func makeGraph(c distributorChannels, filename string, turn int, finalWorld [][]uint8) {
//	c.ioCommand <- ioOutput
//	c.ioFilename <- filename
//	for h, h1 := range finalWorld {
//		for w := range h1 {
//			c.ioOutput <- finalWorld[h][w]
//		}
//	}
//	c.events <- ImageOutputComplete{CompletedTurns: turn, Filename: filename}
//}

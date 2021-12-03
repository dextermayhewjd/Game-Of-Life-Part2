package gol

import (
	"fmt"
	"net/rpc"
	"os"
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


func DistributeCall(client *rpc.Client, StartX, StartY, EndX, EndY  int, currentWorld [][]uint8)  *stubs.Response2{
	request := stubs.Request2{
		StartX: StartX,
		StartY: StartY,
		EndX: EndX,
		EndY: EndY,
		CurrentWorld: currentWorld,

	}
	response := new(stubs.Response2)
	client.Call(stubs.DWHandler, request ,response)
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

func goServerToShutDown(client *rpc.Client ){
	shutDownRequest := stubs.Kill{DeathMessage: "shutdown"}
	response := new(stubs.Response)
	client.Go(stubs.QuitHandler, shutDownRequest, response, nil)
	return
}

func ServerWorker(client *rpc.Client,StartX int, StartY int, EndX int, EndY int,world [][]uint8, out chan<- [][]uint8) {
	result :=DistributeCall(client,StartX,StartY,EndX,EndY,world).PartWorld
	out<- result
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
	var turnCompleted1  = 0

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

	serverList := []string {

		"127.0.0.1:8040",
		"127.0.0.1:8030",
		"127.0.0.1:8050",
		"127.0.0.1:8060",
		//"54.86.80.225:8030",
		//"54.175.11.189:8030",
		//"18.233.6.75;8030",
		//"34.230.71.21:8030",
		//"54.227.116.197:8030",
		////"54.205.61.97:8030",
		//"3.86.60.13:8030",
		//"34.226.140.161:8030",
		//"100.27.27.121:8030",
		//"18.212.240.37:8030",

	}
	serverNum := len(serverList)

	clientList := make([]*rpc.Client,serverNum)
	for i:= range serverList{
		clientList[i], _ = rpc.Dial("tcp",serverList[i])
	}


	go reportTicker(done, eventChan, &currentWorld, mutex, &turnCompleted1)
	go keyPress(c, &turnCompleted1, &currentWorld, keyPresses, clientList,mutex)
	if serverNum == 1 {

		for turnCompleted1 < p.Turns {
			mutex.Lock()
			previousWorld := currentWorld
			currentWorld = makeCall(clientList[0], p.Threads, p.ImageWidth, p.ImageHeight, currentWorld).World
			turnCompleted1 += 1
			compareTwoWorld(eventChan, previousWorld, currentWorld, turnCompleted1)
			mutex.Unlock()
		}
	}else {


		for j:=0 ; j<p.Turns ; j++{
			mutex.Lock()

			out := make([]chan [][]uint8, serverNum)
			for j := 0; j < serverNum; j++ {
				out[j] = make(chan [][]uint8)
			}

			//
			var partFinalWorld [][]uint8
			for k := 0 ; k<serverNum ;k++ {
				//outPartWord[k]=DistributeCall(clientList[k],0,p.ImageHeight/serverNum*k,p.ImageWidth,p.ImageHeight/serverNum*(k+1),currentWorld).PartWorld
				go ServerWorker(clientList[k],0,p.ImageHeight*k/serverNum,p.ImageWidth,p.ImageHeight*(k+1)/serverNum,currentWorld,out[k])
			}

			for i :=0 ; i<serverNum ; i++ {
				partFinalWorld=append(partFinalWorld,<-out[i]...)
				//partFinalWorld=append(partFinalWorld,outPartWord[i]...)
			}
			compareTwoWorld(eventChan,currentWorld,partFinalWorld,turnCompleted1)
			currentWorld = partFinalWorld
			turnCompleted1 += 1
			mutex.Unlock()

		}

	}
	done <- true

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
func keyPress(c distributorChannels, completedTurn *int, currentWorld *[][]uint8, keyPress <-chan rune, clientList []*rpc.Client,mutex *sync.Mutex) {
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

				os.Exit(0)

			case 'k':
				makeGraph(c, filename, *completedTurn, *currentWorld)
				for i :=range clientList{
					goServerToShutDown(clientList[i])
				}


				os.Exit(10)
			case 'p':// pause the sdl
				c.events <- StateChange{NewState: Paused,CompletedTurns: *completedTurn}
				mutex.Lock()

				for 'p' != <- keyPress {

				}
				c.events <- StateChange{NewState: Executing,CompletedTurns: *completedTurn}
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

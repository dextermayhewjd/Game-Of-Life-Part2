package stubs
var GOLHandler = "GameOfLifeOperations.GameOfLife"
type Response struct {
	World [][]uint8
}

type Params struct {
}

type Request struct {
	Turns         int
	Threads       int
	ImageWidth    int
	ImageHeight   int
	CurrentWorld  [][]uint8
	Turn          int
	CompletedTurn chan int
	WorldChan     chan [][]uint8
}

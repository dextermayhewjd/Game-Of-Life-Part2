package stubs
var GOLHandler = "GameOfLifeOperations.GameOfLife"
var QuitHandler = "ServerOperations.ShutDown"
type Response struct {
	World [][]uint8
}

type Params struct {
}

type Request struct {
	Threads       int
	ImageWidth    int
	ImageHeight   int
	CurrentWorld  [][]uint8
}

type Kill struct {
	DeathMessage	string
}



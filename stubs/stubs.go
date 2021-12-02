package stubs
var GOLHandler = "GameOfLifeOperations.GameOfLife"
var DWHandler = "GameOfLifeOperations.DistributedWorld"
var QuitHandler = "ServerOperations.ShutDown"
type Response struct {
	World [][]uint8
}

type Response2 struct {
	PartWorld [][]uint8

}

type Params struct {
}

type Request struct {
	Threads       int
	ImageWidth    int
	ImageHeight   int
	CurrentWorld  [][]uint8
}

type Request2 struct {
	StartX        int
	StartY        int
	EndX          int
	EndY          int
	CurrentWorld  [][]uint8


}


type Kill struct {
	DeathMessage	string
}
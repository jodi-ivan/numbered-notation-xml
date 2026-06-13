package params

import "sync"

type VerseDiagnostic struct {
	SingleMode map[int]bool
	StackMode  map[int]bool
}

type DiagParam struct {
	Mu              *sync.RWMutex
	MapMtx          *sync.Map
	VerseSyllMatch  chan map[int]bool
	VerseDiagnostic chan VerseDiagnostic
	Finish          chan bool
}

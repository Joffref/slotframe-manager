package api

type Slots struct {
	EmittingSlots  map[int]int `json:"emittingSlots"`
	ListeningSlots map[int]int `json:"listeningSlots"`
}

type FrameVersion struct {
	Version int `json:"version"`
}

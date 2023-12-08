package api

// Slots is a struct that represents the emitting and listening slots of a node
type Slots struct {
	// EmittingSlots is a map of emitting slots indexed by their slot number
	EmittingSlots map[int]int `json:"emittingSlots"`
	// ListeningSlots is a map of listening slots indexed by their slot number
	ListeningSlots map[int]int `json:"listeningSlots"`
}

// FrameVersion is a struct that represents the version of a frame
type FrameVersion struct {
	// Version is the version of the frame
	Version int `json:"version"`
}

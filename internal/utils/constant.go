package utils

import "time"

const (
	// KeepAliveInterval is the interval at which the slotframe-manager sends keep-alive messages
	KeepAliveInterval = 10 * time.Second
	// KeepAliveTimeout is the timeout after which a node is considered dead
	KeepAliveTimeout = 6 * KeepAliveInterval
)

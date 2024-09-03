package utils

type ProcessPID struct {
	PID int `json:"pid"`
}

type ProcessState struct {
	PID   int    `json:"pid"`
	State string `json:"state"`
}

type ProcessPath struct {
	Path string `json:"path"`
}

type InterruptReason struct {
	Reason string `json:"reason"`
}

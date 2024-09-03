package internal

type ProcessPath struct {
	Path string `json:"path"`
	Pid  int    `json:"pid"`
}
type ProcessAssets struct {
	Pc  int `json:"pc"`
	Pid int `json:"pid"`
}
type ProcessDelete struct {
	Pid int `json:"pid"`
}
type Page struct {
	PageNumber int `json:"page_number"`
	Pid        int `json:"pid"`
}
type Resize struct {
	Pid    int `json:"pid"`
	Frames int `json:"frames"`
}

type MemStruct struct {
	Pid       int   `json:"pid"`
	Content   int   `json:"content"`
	Length    int   `json:"length"`
	NumFrames []int `json:"numframe"`
	Offset    int   `json:"offset"`
}
type MemCopyString struct {
	Pid           int   `json:"pid"`
	Length        int   `json:"length"`
	NumFramesRead []int `json:"numframeRead"`
	OffsetRead    int   `json:"offsetRead"`
	NumFramesCopy []int `json:"numframeCopy"`
	OffsetCopy    int   `json:"offsetCopy"`
}

type MemStdIO struct {
	Pid       int    `json:"pid"`
	Content   string `json:"content"`
	Length    int    `json:"length"`
	NumFrames []int  `json:"numframe"`
	Offset    int    `json:"offset"`
}

package tlb

type TLBEntry struct {
	PID    int
	Page   int
	Frame  int
	Access int // LRU
}

type TLB struct {
	entries       []TLBEntry
	capacity      int
	replacement   string
	accessCounter int // LRU
	fifoPointer   int // Para FIFO
}

func NewTLB(capacity int, replacement string) *TLB {
	return &TLB{
		entries:     make([]TLBEntry, 0, capacity),
		capacity:    capacity,
		replacement: replacement,
	}
}

func (tlb *TLB) Search(pid, page int) (int, bool) {
	if tlb.capacity == 0 {
		return -1, false // TLB deshabilitada
	}

	for i, entry := range tlb.entries {
		if entry.PID == pid && entry.Page == page {
			if tlb.replacement == "LRU" {
				tlb.entries[i].Access = tlb.accessCounter
				tlb.accessCounter++
			}
			return entry.Frame, true // TLB Hit
		}
	}

	return -1, false // TLB Miss
}

func (tlb *TLB) AddEntry(pid, page, frame int) {
	if tlb.capacity == 0 {
		return // TLB deshabilitada
	}

	if len(tlb.entries) >= tlb.capacity {
		tlb.replaceEntry(pid, page, frame)
	} else {
		tlb.entries = append(tlb.entries, TLBEntry{
			PID:    pid,
			Page:   page,
			Frame:  frame,
			Access: tlb.accessCounter,
		})
		tlb.accessCounter++
	}
}

func (tlb *TLB) replaceEntry(pid, page, frame int) {
	var index int
	if tlb.replacement == "FIFO" {
		index = tlb.fifoPointer
		tlb.fifoPointer = (tlb.fifoPointer + 1) % tlb.capacity // Actualiza el índice FIFO
	} else if tlb.replacement == "LRU" {
		index = tlb.findLRUIndex()
	}

	tlb.entries[index] = TLBEntry{
		PID:    pid,
		Page:   page,
		Frame:  frame,
		Access: tlb.accessCounter,
	}
	tlb.accessCounter++
}

func (tlb *TLB) findLRUIndex() int {
	lruIndex := 0
	minAccess := tlb.entries[0].Access
	for i, entry := range tlb.entries {
		if entry.Access < minAccess {
			lruIndex = i
			minAccess = entry.Access
		}
	}
	return lruIndex
}

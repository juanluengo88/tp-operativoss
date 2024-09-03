package utils

import (
	"fmt"
	"math"
	"time"

	"github.com/sisoputnfrba/tp-golang/kernel/global"
	log "github.com/sisoputnfrba/tp-golang/utils/logger"
	"github.com/sisoputnfrba/tp-golang/utils/model"
)

func TimeCalc(startTime time.Time, quantumTime time.Duration, pcb *model.PCB) int {
	elapsedTime := time.Since(startTime)
	elapsedSeconds := math.Round(elapsedTime.Seconds())
	elapsedMillisRounded := int64(elapsedSeconds * 1000)
	global.Logger.Log("estoy dentro de block", log.DEBUG)
	remainingQuantum := quantumTime - elapsedTime
	remainingMilis := remainingQuantum.Milliseconds()

	global.Logger.Log(fmt.Sprintf("PID: %d - Rounded ElapsedTime: %d ms", pcb.PID, elapsedMillisRounded), log.DEBUG)
	global.Logger.Log(fmt.Sprintf("PID: %d - Rounded RemainingTime: %d ms", pcb.PID, remainingMilis), log.DEBUG)

	return int(remainingMilis)
}

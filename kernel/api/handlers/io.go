package handlers

import (
	"fmt"
	"net/http"

	"github.com/sisoputnfrba/tp-golang/kernel/global"
	log "github.com/sisoputnfrba/tp-golang/utils/logger"
	"github.com/sisoputnfrba/tp-golang/utils/model"
	"github.com/sisoputnfrba/tp-golang/utils/serialization"
)

func ConnectIOHandler(w http.ResponseWriter, r *http.Request) {

	var Device model.IoDevice
	err := serialization.DecodeHTTPBody(r, &Device)
	if err != nil {
		global.Logger.Log("Error al decodear el body: "+err.Error(), log.ERROR)
		http.Error(w, "Error al decodear el body", http.StatusBadRequest)
		return
	}
	Device.Sem = make(chan int, 1)
	global.IoMap[Device.Name] = Device
	global.Logger.Log(fmt.Sprintf("Se conecto %+v", global.IoMap), log.DEBUG)

	w.WriteHeader(http.StatusNoContent)
}

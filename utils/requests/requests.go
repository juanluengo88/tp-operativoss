package requests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func GetHTTP[T any](ip string, port int, endpoint string) (*T, error) {
	url := fmt.Sprintf("http://%s:%d/%s", ip, port, endpoint)

	resp, err := http.Get(url)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data T
	err = json.NewDecoder(resp.Body).Decode(&data)

	if err != nil {
		return nil, err
	}
	return &data, nil
}

func PutHTTPwithBody[T any, R any](ip string, port int, endpoint string, data T) (*R, error) {
	var RespData *R
	url := fmt.Sprintf("http://%s:%d/%s", ip, port, endpoint)

	body, err := json.Marshal(data)
	if err != nil {
		return RespData, err
	}

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(body))
	if err != nil {
		return RespData, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return RespData, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent {
		return nil, nil
	}

	err = json.NewDecoder(resp.Body).Decode(&RespData)
	if err != nil {
		return RespData, err
	}

	return RespData, nil
}

// desarolle una funcion q permite hacer los deletes con plani y process con PID aplicable a la funcion de arriba
// para q funcione deberiamos pasarles como parametro en donde dice endpointwithpid en el caso de un proccess/pid
// en el caso de plani/
func DeleteHTTP[T any](ip string, port int, endpoint string, data T) (*T, error) {
	var RespData T
	url := fmt.Sprintf("http://%s:%d/%s", ip, port, endpoint)
	body, err := json.Marshal(data)
	if err != nil {
		return &RespData, err
	}
	req, err := http.NewRequest(http.MethodDelete, url, bytes.NewBuffer(body))

	if err != nil {
		return &RespData, err
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return &RespData, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent {
		return nil, nil
	}

	err = json.NewDecoder(resp.Body).Decode(&RespData)
	if err != nil {
		return &RespData, err
	}

	return &RespData, nil
}

package handler

import (
	"agent/src"
	"encoding/json"
	"io"
	"net/http"
)

func SendTask(writer http.ResponseWriter, request *http.Request) {
	action := src.Action{}
	body := json.NewDecoder(request.Body)
	if err := body.Decode(&action); err != nil {
		io.WriteString(writer, err.Error())
		return
	}
	writer.Write([]byte("ok"))
}

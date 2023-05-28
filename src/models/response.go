package models

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type MyResponse struct {
	Code         int         `json:"code"`
	Data         interface{} `json:"body"`
	Message      string      `json:"message"`
	ResponseType string      `json:"response_type"`
	WsStatusCode string      `json:"ws_status_code,omitempty"`
}

// MyJSONResponse respuesta JSON standard
type MyJSONResponse struct {
	MyResponse
	status      int
	contentType string
	writer      http.ResponseWriter
}

// GenerateResponse GenerateResponse
func GenerateResponse(w http.ResponseWriter, r MyResponse, status int) {
	MyJSONResponse := CreateJSONResponse(w, r, status)
	MyJSONResponse.SendJSONResponse()
}

// CreateJSONResponse crea una respuesta a partir de r MyResponse en w
func CreateJSONResponse(w http.ResponseWriter, r MyResponse, statusCode int) MyJSONResponse {
	jsonResponse := MyJSONResponse{status: statusCode, contentType: "application/json", writer: w}
	jsonResponse.MyResponse.Data = r.Data
	jsonResponse.MyResponse.Code = r.Code
	jsonResponse.MyResponse.Message = r.Message
	return jsonResponse
}

// SendJSONResponse escribe la respuesta JSON en el writer
func (my *MyJSONResponse) SendJSONResponse() {
	my.writer.Header().Set("Content-Type", my.contentType)
	my.writer.WriteHeader(my.status)
	output, _ := json.Marshal(&my)
	fmt.Fprint(my.writer, string(output))
}

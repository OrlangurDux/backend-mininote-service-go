package middlewares

import (
	"encoding/json"
	"fmt"
	"net/http"
	"orlangur.link/services/mini.note/models"
)

// AuthorizationResponse -> response authorize
func AuthorizationResponse(error models.Error, writer http.ResponseWriter) {
	temp := &models.UniversalDTO{Success: false, StatusCode: 401, Error: error}
	//Send header, status code and output to writer
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusUnauthorized)
	err := json.NewEncoder(writer).Encode(temp)
	if err != nil {
		error.Code = 170
		error.Message = err.Error()
		ServerErrResponse(error, writer)
	}
}

// SuccessResponse -> success formatter
func SuccessResponse(data interface{}, writer http.ResponseWriter) {
	var errors models.Error
	temp := &models.UniversalDTO{StatusCode: 200, Success: true, Data: data}
	//Send header, status code and output to writer
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	err := json.NewEncoder(writer).Encode(temp)
	if err != nil {
		errors.Code = 240
		errors.Message = err.Error()
		ServerErrResponse(errors, writer)
	}
}

// SuccessResponseJwt -> success formatter
func SuccessResponseJwt(jwt models.JWT, writer http.ResponseWriter) {
	var errors models.Error
	//Send header, status code and output to writer
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	err := json.NewEncoder(writer).Encode(jwt)
	if err != nil {
		errors.Code = 180
		errors.Message = err.Error()
		ServerErrResponse(errors, writer)
	}
}

// ErrorResponse -> error formatter
func ErrorResponse(error models.Error, writer http.ResponseWriter) {
	var errors models.Error
	temp := &models.UniversalDTO{Success: false, StatusCode: 400, Error: error}
	//Send header, status code and output to writer
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusBadRequest)
	err := json.NewEncoder(writer).Encode(temp)
	if err != nil {
		errors.Code = 190
		errors.Message = err.Error()
		ServerErrResponse(errors, writer)
	}
}

// ServerErrResponse -> server error formatter
func ServerErrResponse(error models.Error, writer http.ResponseWriter) {
	temp := &models.UniversalDTO{Success: false, StatusCode: 500, Error: error}
	//Send header, status code and output to writer
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusInternalServerError)
	err := json.NewEncoder(writer).Encode(temp)
	if err != nil {
		fmt.Println(err.Error())
	}
}

// ValidationResponse -> user input validation
func ValidationResponse(fields map[string][]string, writer http.ResponseWriter) {
	var errors models.Error
	response := make(map[string]interface{})
	response["errors"] = fields
	response["status"] = 422
	response["msg"] = "validation error"
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusUnprocessableEntity)
	err := json.NewEncoder(writer).Encode(response)
	if err != nil {
		errors.Code = 170
		errors.Message = err.Error()
		ServerErrResponse(errors, writer)
	}
}

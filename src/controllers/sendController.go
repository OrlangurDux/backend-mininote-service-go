package controllers

import (
	"bytes"
	"encoding/json"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"io"
	"net/http"
	middlewares "orlangur.link/services/mini.note/handlers"
	"orlangur.link/services/mini.note/helpers"
	"orlangur.link/services/mini.note/models"
	"reflect"
	"strings"
)

// SendRequest godoc
// @Summary      Send request
// @Description  Send message to email
// @Tags         Request
// @Accept       json
// @Produce      json
// @Param        request   body      models.Request  true  "Request param"
// @Success      200  {object}  models.UniversalDTO "ok"
// @Failure      400  {object}  models.UniversalDTO "error"
// @Failure      404  {object}  models.UniversalDTO "error"
// @Failure      500  {object}  models.UniversalDTO "error"
// @Security BearerAuth
// @Router       /send/request [post]
func (c *Controller) SendRequest(response http.ResponseWriter, request *http.Request) {
	var req models.Request
	var errors models.Error
	reqBody, _ := io.ReadAll(request.Body)
	request.Body.Close()
	request.Body = io.NopCloser(bytes.NewBuffer(reqBody))
	err := json.NewDecoder(request.Body).Decode(&req)
	if err != nil {
		request.Body = io.NopCloser(bytes.NewBuffer(reqBody))
		if err = request.ParseForm(); err != nil {
			errors.Code = 220
			errors.Message = err.Error()
			middlewares.ServerErrResponse(errors, response)
			return
		}
		pointer := reflect.ValueOf(&req)
		fields := pointer.Elem()
		if fields.Kind() == reflect.Struct {
			for k, v := range request.PostForm {
				k = cases.Title(language.Und, cases.NoLower).String(strings.ToLower(k))
				field := fields.FieldByName(k)
				if field.IsValid() {
					if field.CanSet() {
						if field.Kind() == reflect.String {
							field.SetString(v[0])
						}
					}
				}
			}
		}
	}
	subject := middlewares.DotEnvVariable("SMTP_SUBJECT")
	message := middlewares.DotEnvVariable("SMTP_MESSAGE")
	unpackReq := reflect.ValueOf(req)
	unpackKeys := unpackReq.Type()
	for i := 0; i < unpackReq.NumField(); i++ {
		message = strings.ReplaceAll(message, "#"+strings.ToLower(unpackKeys.Field(i).Name)+"#", unpackReq.Field(i).String())
	}
	to := []string{middlewares.DotEnvVariable("SMTP_TO")}
	err = helpers.Mail(to, subject, message)
	if err != nil {
		errors.Code = 220
		errors.Message = err.Error()
		middlewares.ServerErrResponse(errors, response)
		return
	}
	middlewares.SuccessResponse("Send success", response)
}

package controllers

import (
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	middlewares "orlangur.link/services/mini.note/handlers"
	"orlangur.link/services/mini.note/models"
)

// Controller -> base controller
type Controller struct {
	MG *mongo.Client
}

// BaseController -> base controller
func BaseController(mg *mongo.Client) Controller {
	return Controller{mg}
}

// GetVersion godoc
// @Summary      Get version
// @Description  Get version app
// @Tags         Help
// @Produce      json
// @Success      200  {object}  models.UniversalDTO "Version"
// @Failure      400  {object}  models.UniversalDTO "error"
// @Failure      404  {object}  models.UniversalDTO "error"
// @Failure      500  {object}  models.UniversalDTO "error"
// @Router       /version [get]
func (c Controller) GetVersion(response http.ResponseWriter, request *http.Request) {
	data := struct {
		Version string `json:"version"`
		Author  string `json:"author"`
		Contact string `json:"contact"`
	}{
		Version: models.Version,
		Author:  "Alexey (Orlangur)",
		Contact: "o@orlangur.link",
	}
	middlewares.SuccessResponse(data, response)
}

// Pagination -> create pagination for item
func Pagination(records interface{}, request *http.Request) (interface{}, error) {
	return nil, nil
}

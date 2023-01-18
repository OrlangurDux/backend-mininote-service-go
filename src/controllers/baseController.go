package controllers

import (
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

// Controller -> base controller
type Controller struct {
	MG *mongo.Client
}

// BaseController -> base controller
func BaseController(mg *mongo.Client) Controller {
	return Controller{mg}
}

func Pagination(records interface{}, request *http.Request) (interface{}, error) {
	return nil, nil
}

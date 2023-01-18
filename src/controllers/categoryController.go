package controllers

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	middlewares "orlangur.link/services/mini.note/handlers"
	"orlangur.link/services/mini.note/helpers"
	"orlangur.link/services/mini.note/models"
	"strconv"
)

// CategoryCreateEndpoint godoc
// @Summary      Category record
// @Description  Create category record
// @Tags         Categories
// @Accept       x-www-form-urlencoded
// @Produce      json
// @Param        name  formData   string  true  "Category name"
// @Param        parent_id   formData   string  false  "Category parent_id"
// @Param        sort formData   number  false  "Category sort"
// @Success      200  {object}  models.UniversalDTO "Create category item"
// @Failure      400  {object}  models.UniversalDTO "error"
// @Failure      404  {object}  models.UniversalDTO "error"
// @Failure      500  {object}  models.UniversalDTO "error"
// @Security BearerAuth
// @Router       /categories [post]
func (c Controller) CategoryCreateEndpoint(response http.ResponseWriter, request *http.Request) {
	var errors models.Error
	var category models.Category
	err := request.ParseForm()
	if err != nil {
		errors.Code = 470
		errors.Message = err.Error()
		middlewares.ErrorResponse(errors, response)
		return
	}
	name := request.FormValue("name")
	parentId := request.FormValue("parent_id")
	sort := request.FormValue("sort")

	category.Id = primitive.NewObjectID()
	category.UserId, err = helpers.GetUserID()
	if err != nil {
		errors.Code = 475
		errors.Message = err.Error()
		middlewares.ErrorResponse(errors, response)
		return
	}
	category.Name = name
	if parentId != "" {
		category.ParentId, err = primitive.ObjectIDFromHex(parentId)
		if err != nil {
			errors.Code = 480
			errors.Message = err.Error()
			middlewares.ErrorResponse(errors, response)
			return
		}
	}
	category.Sort, err = strconv.Atoi(sort)
	if err != nil {
		errors.Code = 490
		errors.Message = err.Error()
		middlewares.ErrorResponse(errors, response)
		return
	}

	collection := c.MG.Database("notes").Collection("categories")
	one, err := collection.InsertOne(context.TODO(), category)
	if err != nil {
		errors.Code = 500
		errors.Message = err.Error()
		middlewares.ErrorResponse(errors, response)
		return
	}
	id := one.InsertedID.(primitive.ObjectID).Hex()
	obj := struct {
		Id      string `json:"id"`
		Message string `json:"message"`
	}{
		Id:      id,
		Message: "Category added",
	}
	middlewares.SuccessResponse(obj, response)
}

// CategoryReadEndpoint -> note read
// @Summary      Category element
// @Description  View category element
// @Tags         Categories
// @Produce      json
// @Param		 id path string true "category id"
// @Success      200  {object}  models.UniversalDTO{data=models.Category} "Category item"
// @Failure      400  {object}  models.UniversalDTO "error"
// @Failure      404  {object}  models.UniversalDTO "error"
// @Failure      500  {object}  models.UniversalDTO "error"
// @Security BearerAuth
// @Router       /categories/{id} [get]
func (c Controller) CategoryReadEndpoint(response http.ResponseWriter, request *http.Request) {
	var errors models.Error
	var category models.Category

	param := mux.Vars(request)
	id, err := primitive.ObjectIDFromHex(param["id"])
	if err != nil {
		errors.Code = 580
		errors.Message = err.Error()
		middlewares.ErrorResponse(errors, response)
		return
	}
	userId, err := helpers.GetUserID()
	if err != nil {
		errors.Code = 590
		errors.Message = err.Error()
		middlewares.ErrorResponse(errors, response)
		return
	}
	filter := bson.M{"_id": id, "user_id": userId}
	collection := c.MG.Database("notes").Collection("categories")
	item := collection.FindOne(context.TODO(), filter)
	err = item.Decode(&category)
	if err != nil {
		errors.Code = 600
		errors.Message = err.Error()
		middlewares.ErrorResponse(errors, response)
		return
	}
	middlewares.SuccessResponse(category, response)
}

// CategoryUpdateEndpoint godoc
// @Summary      Category record
// @Description  Update category record
// @Tags         Categories
// @Accept       x-www-form-urlencoded
// @Produce      json
// @Param		 id path string true "category id"
// @Param        name  formData   string  true  "Category name"
// @Param        parent_id   formData   string  false  "Category parent id"
// @Param        sort formData number false "Category sort"
// @Success      200  {object}  models.UniversalDTO "Update category item"
// @Failure      400  {object}  models.UniversalDTO "error"
// @Failure      404  {object}  models.UniversalDTO "error"
// @Failure      500  {object}  models.UniversalDTO "error"
// @Security BearerAuth
// @Router       /categories/{id} [put]
func (c Controller) CategoryUpdateEndpoint(response http.ResponseWriter, request *http.Request) {
	var errors models.Error
	var iSort int
	param := mux.Vars(request)
	id, err := primitive.ObjectIDFromHex(param["id"])
	if err != nil {
		errors.Code = 525
		errors.Message = err.Error()
		middlewares.ErrorResponse(errors, response)
		return
	}
	err = request.ParseForm()
	if err != nil {
		errors.Code = 530
		errors.Message = err.Error()
		middlewares.ErrorResponse(errors, response)
		return
	}
	name := request.FormValue("name")
	parentId := request.FormValue("parent_id")
	sort := request.FormValue("sort")
	if sort != "" {
		iSort, err = strconv.Atoi(request.FormValue("sort"))
		if err != nil {
			errors.Code = 540
			errors.Message = err.Error()
			middlewares.ErrorResponse(errors, response)
			return
		}
	}
	data := bson.M{}
	data["name"] = name
	if parentId != "" {
		data["parent_id"] = parentId
	}
	if sort != "" {
		data["sort"] = iSort
	}
	update := bson.M{"$set": data}
	collection := c.MG.Database("notes").Collection("categories")
	_, err = collection.UpdateByID(context.TODO(), id, update)
	if err != nil {
		errors.Code = 550
		errors.Message = err.Error()
		middlewares.ErrorResponse(errors, response)
		return
	}
	success := fmt.Sprintf("Category ID:%s updated", id)
	middlewares.SuccessResponse(success, response)
}

// CategoryDeleteEndpoint godoc
// @Summary      Category delete
// @Description  Delete category record
// @Tags         Categories
// @Produce      json
// @Param		 id path string true "category id"
// @Success      200  {object}  models.UniversalDTO "Delete category record"
// @Failure      400  {object}  models.UniversalDTO "error"
// @Failure      404  {object}  models.UniversalDTO "error"
// @Failure      500  {object}  models.UniversalDTO "error"
// @Security BearerAuth
// @Router       /categories/{id} [delete]
func (c Controller) CategoryDeleteEndpoint(response http.ResponseWriter, request *http.Request) {
	var errors models.Error
	param := mux.Vars(request)
	id, err := primitive.ObjectIDFromHex(param["id"])
	if err != nil {
		errors.Code = 560
		errors.Message = err.Error()
		middlewares.ErrorResponse(errors, response)
		return
	}

	collection := c.MG.Database("notes").Collection("categories")
	del := bson.M{"_id": id}
	_, err = collection.DeleteOne(context.TODO(), del)
	if err != nil {
		errors.Code = 570
		errors.Message = err.Error()
		middlewares.ErrorResponse(errors, response)
		return
	}
	middlewares.SuccessResponse("Category deleted", response)
}

// CategoryListEndpoint godoc
// @Summary      Categories list
// @Description  View categories list
// @Tags         Categories
// @Produce      json
// @Success      200  {object}  models.UniversalDTO{data=models.Categories} "Categories items list"
// @Failure      400  {object}  models.UniversalDTO "error"
// @Failure      404  {object}  models.UniversalDTO "error"
// @Failure      500  {object}  models.UniversalDTO "error"
// @Security BearerAuth
// @Router       /categories [get]
func (c Controller) CategoryListEndpoint(response http.ResponseWriter, request *http.Request) {
	var errors models.Error
	var categories models.Categories
	var category []*models.Category
	collection := c.MG.Database("notes").Collection("categories")
	userId, _ := helpers.GetUserID()
	filter := bson.M{"user_id": userId}
	items, err := collection.Find(context.TODO(), filter)
	if err != nil {
		errors.Code = 510
		errors.Message = err.Error()
		middlewares.ErrorResponse(errors, response)
		return
	}
	err = items.All(context.TODO(), &category)
	if err != nil {
		errors.Code = 520
		errors.Message = err.Error()
		middlewares.ErrorResponse(errors, response)
		return
	}
	categories.Total = len(category)
	categories.Items = category
	middlewares.SuccessResponse(categories, response)
}

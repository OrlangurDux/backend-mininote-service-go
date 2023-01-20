package controllers

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	options2 "go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	middlewares "orlangur.link/services/mini.note/handlers"
	"orlangur.link/services/mini.note/helpers"
	"orlangur.link/services/mini.note/models"
	"strconv"
	"time"
)

// NoteListEndpoint godoc
// @Summary      Notes list
// @Description  View notes list
// @Accept       x-www-form-urlencoded
// @Tags         Note
// @Produce      json
// @Param        page	query	string	false	"Number page"
// @Param        per_page	query	string 	false	"Number per page"
// @Success      200  {object}  models.UniversalDTO{data=models.Notes} "Note items list"
// @Failure      400  {object}  models.UniversalDTO "error"
// @Failure      404  {object}  models.UniversalDTO "error"
// @Failure      500  {object}  models.UniversalDTO "error"
// @Security BearerAuth
// @Router       /notes [get]
func (c Controller) NoteListEndpoint(response http.ResponseWriter, request *http.Request) {
	var errors models.Error
	var notes []*models.Note
	var noteList models.Notes
	userID, err := helpers.GetUserID()
	if err != nil {
		errors.Code = 295
		errors.Message = err.Error()
		middlewares.ErrorResponse(errors, response)
		return
	}
	page, _ := strconv.ParseInt(request.URL.Query().Get("page"), 10, 32)
	perPage, _ := strconv.ParseInt(request.URL.Query().Get("per_page"), 10, 32)
	if perPage == 0 {
		sPerPage := middlewares.DotEnvVariable("DEFAULT_PER_PAGE")
		if sPerPage != "" {
			perPage, err = strconv.ParseInt(sPerPage, 10, 32)
			if err != nil {
				errors.Code = 296
				errors.Message = err.Error()
				middlewares.ErrorResponse(errors, response)
				return
			}
		} else {
			perPage = 10
		}
	}
	if page > 0 {
		page = page - 1
	}
	offset := page * perPage
	collection := c.MG.Database("notes").Collection("notes")
	options := options2.Find()
	options.SetSkip(offset)
	options.SetLimit(perPage)
	filter := bson.M{"user_id": userID}
	total, err := collection.CountDocuments(context.TODO(), filter)
	if err != nil {
		errors.Code = 297
		errors.Message = err.Error()
		middlewares.ErrorResponse(errors, response)
		return
	}
	find, err := collection.Find(context.TODO(), filter, options)
	if err != nil {
		errors.Code = 298
		errors.Message = err.Error()
		middlewares.ErrorResponse(errors, response)
		return
	}
	err = find.All(context.TODO(), &notes)
	if err != nil {
		errors.Code = 299
		errors.Message = err.Error()
		middlewares.ErrorResponse(errors, response)
		return
	}
	noteList.Total = total
	noteList.Page = page
	noteList.PerPage = perPage
	noteList.Items = notes
	middlewares.SuccessResponse(noteList, response)
}

// NoteReadEndpoint -> note read
// @Summary      Note element
// @Description  View note element
// @Tags         Note
// @Produce      json
// @Param		 id path string true "note id"
// @Success      200  {object}  models.UniversalDTO{data=models.Note} "Note item"
// @Failure      400  {object}  models.UniversalDTO "error"
// @Failure      404  {object}  models.UniversalDTO "error"
// @Failure      500  {object}  models.UniversalDTO "error"
// @Security BearerAuth
// @Router       /notes/{id} [get]
func (c Controller) NoteReadEndpoint(response http.ResponseWriter, request *http.Request) {
	var note models.Note
	var errors models.Error
	param := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(param["id"])
	collection := c.MG.Database("notes").Collection("notes")
	filter := bson.M{"_id": id}
	err := collection.FindOne(context.TODO(), filter).Decode(&note)
	if err != nil {
		errors.Code = 300
		errors.Message = err.Error()
		middlewares.ErrorResponse(errors, response)
		return
	}
	middlewares.SuccessResponse(note, response)
}

// NoteCreateEndpoint godoc
// @Summary      Note record
// @Description  Create note record
// @Tags         Note
// @Accept       x-www-form-urlencoded
// @Produce      json
// @Param        title  formData   string  true  "Note title"
// @Param        note   formData   string  true  "Note body"
// @Param        category_id formData string false "Note category"
// @Param        status formData   string  true  "Note status" Enums(draft,public,archive)
// @Success      200  {object}  models.UniversalDTO "Create note item"
// @Failure      400  {object}  models.UniversalDTO "error"
// @Failure      404  {object}  models.UniversalDTO "error"
// @Failure      500  {object}  models.UniversalDTO "error"
// @Security BearerAuth
// @Router       /notes [post]
func (c Controller) NoteCreateEndpoint(response http.ResponseWriter, request *http.Request) {
	var record models.Note
	var errors models.Error
	err := request.ParseForm()
	if err != nil {
		errors.Code = 270
		errors.Message = err.Error()
		middlewares.ErrorResponse(errors, response)
		return
	}
	title := request.PostFormValue("title")
	note := request.PostFormValue("note")
	status := request.PostFormValue("status")
	categoryID := request.PostFormValue("category_id")
	uid, err := helpers.GetUserID()
	if err != nil {
		errors.Code = 280
		errors.Message = err.Error()
		middlewares.ErrorResponse(errors, response)
		return
	}
	record.Id = primitive.NewObjectID()
	record.Title = title
	record.Note = note
	record.Status = status
	record.UserID = uid
	if categoryID != "" {
		record.CategoryID, err = primitive.ObjectIDFromHex(categoryID)
		if err != nil {
			errors.Code = 285
			errors.Message = err.Error()
			middlewares.ErrorResponse(errors, response)
			return
		}
	}
	record.CreatedAt = time.Now()
	record.UpdatedAt = time.Now()

	collection := c.MG.Database("notes").Collection("notes")
	one, err := collection.InsertOne(context.TODO(), record)
	if err != nil {
		errors.Code = 290
		errors.Message = err.Error()
		middlewares.ErrorResponse(errors, response)
		return
	}
	id := one.InsertedID.(primitive.ObjectID).Hex()
	obj := struct {
		ID      string `json:"id"`
		Message string `json:"message"`
	}{
		ID:      id,
		Message: "Note added",
	}
	middlewares.SuccessResponse(obj, response)
}

// NoteUpdateEndpoint godoc
// @Summary      Note record
// @Description  Update note record
// @Tags         Note
// @Accept       x-www-form-urlencoded
// @Produce      json
// @Param		 id path string true "note id"
// @Param        title  formData   string  true  "Note title"
// @Param        note   formData   string  true  "Note body"
// @Param        status formData   string  true  "Note status" Enums(draft,public,archive)
// @Success      200  {object}  models.UniversalDTO "Update note item"
// @Failure      400  {object}  models.UniversalDTO "error"
// @Failure      404  {object}  models.UniversalDTO "error"
// @Failure      500  {object}  models.UniversalDTO "error"
// @Security BearerAuth
// @Router       /notes/{id} [put]
func (c Controller) NoteUpdateEndpoint(response http.ResponseWriter, request *http.Request) {
	//var note models.Note
	var errors models.Error
	data := bson.M{}
	param := mux.Vars(request)
	id, err := primitive.ObjectIDFromHex(param["id"])
	if err != nil {
		errors.Code = 300
		errors.Message = err.Error()
		middlewares.ErrorResponse(errors, response)
		return
	}
	err = request.ParseForm()
	if err != nil {
		errors.Code = 310
		errors.Message = err.Error()
		middlewares.ErrorResponse(errors, response)
		return
	}

	for key, value := range request.PostForm {
		if value[0] != "" {
			data[key] = value[0]
		}
	}
	collection := c.MG.Database("notes").Collection("notes")
	update := bson.M{"$set": data}
	_, err = collection.UpdateByID(context.TODO(), id, update)
	if err != nil {
		errors.Code = 320
		errors.Message = err.Error()
		middlewares.ErrorResponse(errors, response)
		return
	}
	success := fmt.Sprintf("Note ID:%s updated", id)
	middlewares.SuccessResponse(success, response)
}

// NoteDeleteEndpoint godoc
// @Summary      Note delete
// @Description  Delete note record
// @Tags         Note
// @Produce      json
// @Param		 id path string true "note id"
// @Success      200  {object}  models.UniversalDTO "Delete note record"
// @Failure      400  {object}  models.UniversalDTO "error"
// @Failure      404  {object}  models.UniversalDTO "error"
// @Failure      500  {object}  models.UniversalDTO "error"
// @Security BearerAuth
// @Router       /notes/{id} [delete]
func (c Controller) NoteDeleteEndpoint(response http.ResponseWriter, request *http.Request) {
	var errors models.Error
	param := mux.Vars(request)
	id, err := primitive.ObjectIDFromHex(param["id"])
	if err != nil {
		errors.Code = 325
		errors.Message = err.Error()
		middlewares.ErrorResponse(errors, response)
		return
	}
	userID, err := helpers.GetUserID()
	if err != nil {
		errors.Code = 330
		errors.Message = err.Error()
		middlewares.ErrorResponse(errors, response)
		return
	}
	collection := c.MG.Database("notes").Collection("notes")
	filter := bson.M{"_id": id, "user_id": userID}
	_, err = collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		errors.Code = 340
		errors.Message = err.Error()
		middlewares.ErrorResponse(errors, response)
		return
	}
	middlewares.SuccessResponse("Deleted success", response)
}

// NoteSearchEndpoint godoc
// @Summary      Note search
// @Description  Search note records
// @Tags         Note
// @Produce      json
// @Param		 q	query string false "note search by q" Format(text)
// @Success      200  {object}  models.UniversalDTO "Search note items list"
// @Failure      400  {object}  models.UniversalDTO "error"
// @Failure      404  {object}  models.UniversalDTO "error"
// @Failure      500  {object}  models.UniversalDTO "error"
// @Security BearerAuth
// @Router       /notes/search [get]
func (c Controller) NoteSearchEndpoint(response http.ResponseWriter, request *http.Request) {
	var errors models.Error
	var notes []*models.Note
	var pipeline []bson.M
	q := request.URL.Query().Get("q")
	userID, err := helpers.GetUserID()
	if err != nil {
		errors.Code = 295
		errors.Message = err.Error()
		middlewares.ErrorResponse(errors, response)
		return
	}
	collection := c.MG.Database("notes").Collection("notes")
	var orKey []interface{}
	orKey = append(orKey, bson.M{"title": bson.M{"$regex": q, "$options": "i"}})
	orKey = append(orKey, bson.M{"note": bson.M{"$regex": q, "$options": "i"}})
	matchStage := bson.M{"$match": bson.M{"$and": []bson.M{{"user_id": userID}, {"$or": orKey}}}}
	sortStage := bson.M{"$sort": bson.M{"updated_at": -1}}
	pipeline = append(pipeline, matchStage, sortStage)
	find, err := collection.Aggregate(context.TODO(), pipeline)
	if err != nil {
		errors.Code = 298
		errors.Message = err.Error()
		middlewares.ErrorResponse(errors, response)
		return
	}
	err = find.All(context.TODO(), &notes)
	if err != nil {
		errors.Code = 299
		errors.Message = err.Error()
		middlewares.ErrorResponse(errors, response)
		return
	}
	middlewares.SuccessResponse(notes, response)
}

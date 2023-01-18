package controllers

import (
	"context"
	"crypto/md5"
	"fmt"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	middlewares "orlangur.link/services/mini.note/handlers"
	"orlangur.link/services/mini.note/helpers"
	"orlangur.link/services/mini.note/models"
)

// UserRegisterEndpoint godoc
// @Summary      Register user account
// @Description  Registration user by email and password
// @Tags         User
// @Accept       x-www-form-urlencoded
// @Produce      json
// @Param        email   formData  string true  "Email"
// @Param        password formData  string  true  "Password"
// @Success      200  {object}  models.UniversalDTO "ok"
// @Failure      400  {object}  models.UniversalDTO "error"
// @Failure      404  {object}  models.UniversalDTO "error"
// @Failure      500  {object}  models.UniversalDTO "error"
// @Router       /users/register [post]
func (c Controller) UserRegisterEndpoint(response http.ResponseWriter, request *http.Request) {
	var user models.User
	var errors models.Error
	err := request.ParseForm()

	if err != nil {
		errors.Code = 5
		errors.Message = err.Error()
		middlewares.ErrorResponse(errors, response)
		return
	}

	email := request.PostFormValue("email")
	password := request.PostFormValue("password")

	filter := bson.M{"email": email}
	collection := c.MG.Database("notes").Collection("users")

	err = collection.FindOne(context.TODO(), filter).Decode(&user)

	if err != nil {
		hash := md5.Sum([]byte(password))
		user.Id = primitive.NewObjectID()
		user.Email = email
		user.Password = fmt.Sprintf("%x", hash)
		user.Active = true
		user.CreatedAt = time.Now()
		user.UpdatedAt = time.Now()
		_, err := collection.InsertOne(context.TODO(), user)
		if err != nil {
			errors.Code = 10
			errors.Message = err.Error()
			middlewares.ErrorResponse(errors, response)
			return
		}
		middlewares.SuccessResponse("User added", response)
	} else {
		errors.Code = 20
		errors.Message = "User " + email + " exist."
		middlewares.ErrorResponse(errors, response)
		return
	}
}

// UserLoginEndpoint godoc
// @Summary      Login user account
// @Description  Authorization user account by email and password
// @Tags         User
// @Accept       x-www-form-urlencoded
// @Produce      json
// @Param        email   formData  string true  "Email"
// @Param        password   formData  string true  "Password"
// @Success      200  {object}  models.JWT "ok"
// @Failure      400  {object}  models.UniversalDTO "error"
// @Failure      404  {object}  models.UniversalDTO "error"
// @Failure      500  {object}  models.UniversalDTO "error"
// @Router       /users/login [post]
func (c Controller) UserLoginEndpoint(response http.ResponseWriter, request *http.Request) {
	var user models.User
	var errors models.Error
	err := request.ParseForm()

	if err != nil {
		errors.Code = 30
		errors.Message = err.Error()
		middlewares.ErrorResponse(errors, response)
		return
	}

	email := request.PostFormValue("email")
	password := request.PostFormValue("password")

	filter := bson.M{"email": email}
	collection := c.MG.Database("notes").Collection("users")

	err = collection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		errors.Code = 40
		errors.Message = "User " + email + " not exist"
		middlewares.ErrorResponse(errors, response)
		return
	}

	filter = bson.M{"email": email, "active": true}
	err = collection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		errors.Code = 45
		errors.Message = "User " + email + " not active"
		middlewares.ErrorResponse(errors, response)
		return
	}

	hash := fmt.Sprintf("%x", md5.Sum([]byte(password)))
	filter = bson.M{"email": email, "password": hash}
	err = collection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		errors.Code = 50
		errors.Message = "Password incorrect"
		middlewares.ErrorResponse(errors, response)
		return
	}

	jwt, err := middlewares.GenerateJWT(user)
	if err != nil {
		errors.Code = 60
		errors.Message = err.Error()
		middlewares.ErrorResponse(errors, response)
		return
	}

	update := bson.M{"$set": bson.M{"authorized_at": time.Now()}}
	_, err = collection.UpdateByID(context.TODO(), user.Id, update)

	if err != nil {
		errors.Code = 70
		errors.Message = err.Error()
		middlewares.ErrorResponse(errors, response)
		return
	}

	middlewares.SuccessResponseJwt(jwt, response)
}

// UserForgotEndpoint godoc
// @Summary      Forgot user password
// @Description  Recovery user password by email
// @Tags         User
// @Accept       x-www-form-urlencoded
// @Produce      json
// @Param        email           formData   string  false  "Email"
// @Param 		 restore_token   formData   string  false  "Restore Token"
// @Param		 password        formData   string  false  "New password"
// @Success      200  {object}  models.UniversalDTO "ok"
// @Failure      400  {object}  models.UniversalDTO "error"
// @Failure      404  {object}  models.UniversalDTO "error"
// @Failure      500  {object}  models.UniversalDTO "error"
// @Router       /users/forgot [post]
func (c Controller) UserForgotEndpoint(response http.ResponseWriter, request *http.Request) {
	var errors models.Error
	var user models.User
	err := request.ParseForm()
	if err != nil {
		errors.Code = 71
		errors.Message = err.Error()
		middlewares.ErrorResponse(errors, response)
		return
	}
	email := request.PostFormValue("email")
	restoreToken := request.PostFormValue("restore_token")
	newPassword := request.PostFormValue("password")
	responseMessage := ""

	collection := c.MG.Database("notes").Collection("users")
	if (email != "" && restoreToken != "" && newPassword != "") || (newPassword != "" && restoreToken != "") {
		filter := bson.M{"restore_token": restoreToken}
		err := collection.FindOne(context.TODO(), filter).Decode(&user)
		if err != nil {
			errors.Code = 73
			errors.Message = "Restore Token incorrect please try again"
			middlewares.ErrorResponse(errors, response)
			return
		}
		hash := fmt.Sprintf("%x", md5.Sum([]byte(newPassword)))
		update := bson.M{"$set": bson.M{"password": hash, "restore_token": ""}}
		_, err = collection.UpdateByID(context.TODO(), user.Id, update)
		if err != nil {
			errors.Code = 74
			errors.Message = err.Error()
			middlewares.ErrorResponse(errors, response)
			return
		}
		responseMessage = "Update password success"
	} else if email != "" {
		restoreToken = helpers.RandomString(32)
		subject := fmt.Sprintf("Token fo recovery password %s", middlewares.DotEnvVariable("HOST"))
		message := fmt.Sprintf("Link for recovery password %sforgot/?token=%s", middlewares.DotEnvVariable("HOST"), restoreToken)
		filter := bson.M{"email": email}
		update := bson.M{"$set": bson.M{"restore_token": restoreToken}}
		err := collection.FindOneAndUpdate(context.TODO(), filter, update).Decode(&user)
		if err != nil {
			errors.Code = 75
			errors.Message = err.Error()
			middlewares.ErrorResponse(errors, response)
			return
		}
		err = helpers.Mail([]string{email}, subject, message)
		if err != nil {
			errors.Code = 76
			errors.Message = err.Error()
			middlewares.ErrorResponse(errors, response)
			return
		}
		responseMessage = "Restore Token send to email address"
	} else {
		errors.Code = 77
		errors.Message = "Fields not be empty"
		middlewares.ErrorResponse(errors, response)
		return
	}
	obj := struct {
		Token   string `json:"token"`
		Message string `json:"message"`
	}{
		Token:   restoreToken,
		Message: responseMessage,
	}
	middlewares.SuccessResponse(obj, response)
}

// UserCheckByEmailEndpoint godoc
// @Summary      User check
// @Description  Check exist user by email
// @Tags         User
// @Produce      json
// @Param        email           formData   string  false  "Email"
// @Success      200  {object}  models.UniversalDTO "ok"
// @Failure      400  {object}  models.UniversalDTO "error"
// @Failure      404  {object}  models.UniversalDTO "error"
// @Failure      500  {object}  models.UniversalDTO "error"
// @Router       /users/check [post]
func (c Controller) UserCheckByEmailEndpoint(response http.ResponseWriter, request *http.Request) {
	var user models.User
	var errors models.Error
	err := request.ParseForm()
	if err != nil {
		errors.Code = 160
		errors.Message = err.Error()
		middlewares.ErrorResponse(errors, response)
		return
	}
	email := request.PostFormValue("email")
	collection := c.MG.Database("notes").Collection("users")
	if email != "" {
		filter := bson.M{"email": email}
		err = collection.FindOne(context.TODO(), filter).Decode(&user)
		if err != nil {
			errors.Code = 40
			errors.Message = "User " + email + " not exist"
			middlewares.ErrorResponse(errors, response)
			return
		}
	} else {
		errors.Code = 77
		errors.Message = "Fields not empty"
		middlewares.ErrorResponse(errors, response)
		return
	}
	success := "User " + email + " exists"
	middlewares.SuccessResponse(success, response)
}

// UserProfileReadEndpoint godoc
// @Summary      User profile
// @Description  View user profile information
// @Tags         User
// @Produce      json
// @Success      200  {object}  models.UniversalDTO "ok"
// @Failure      400  {object}  models.UniversalDTO "error"
// @Failure      404  {object}  models.UniversalDTO "error"
// @Failure      500  {object}  models.UniversalDTO "error"
// @Security BearerAuth
// @Router       /users/profile [get]
func (c Controller) UserProfileReadEndpoint(response http.ResponseWriter, request *http.Request) {
	var user models.User
	var errors models.Error
	userId, err := helpers.GetUserID()
	if err != nil {
		errors.Code = 78
		errors.Message = err.Error()
		middlewares.ErrorResponse(errors, response)
		return
	}

	collection := c.MG.Database("notes").Collection("users")
	filter := bson.M{"_id": userId}
	err = collection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		errors.Code = 80
		errors.Message = err.Error()
		middlewares.ErrorResponse(errors, response)
		return
	}
	middlewares.SuccessResponse(user, response)
}

// UserProfileUpdateEndpoint godoc
// @Summary      User update
// @Description  Update user profile information
// @Tags         User
// @Accept       x-www-form-urlencoded
// @Produce      json
// @Param        name   formData  string false  "Name"
// @Param        avatar   formData  file false  "Avatar"
// @Success      200  {object}  models.UniversalDTO "ok"
// @Failure      400  {object}  models.UniversalDTO "error"
// @Failure      404  {object}  models.UniversalDTO "error"
// @Failure      500  {object}  models.UniversalDTO "error"
// @Security BearerAuth
// @Router       /users/profile [put]
func (c Controller) UserProfileUpdateEndpoint(response http.ResponseWriter, request *http.Request) {
	var errors models.Error
	data := bson.M{}
	if err := request.ParseMultipartForm(32 << 20); err != nil {
		errors.Code = 85
		errors.Message = err.Error()
		middlewares.ErrorResponse(errors, response)
		return
	} else {
		avatar, handler, err := request.FormFile("avatar")
		if err == nil {
			path, err := helpers.UploadAvatar(avatar, handler)
			defer avatar.Close()
			if err != nil {
				errors.Code = 87
				errors.Message = err.Error()
				middlewares.ErrorResponse(errors, response)
				return
			}
			data["avatar"] = path
		}
		for key, value := range request.PostForm {
			if value[0] != "" {
				data[key] = value[0]
			}
		}
	}

	userId, err := helpers.GetUserID()
	if err != nil {
		errors.Code = 90
		errors.Message = err.Error()
		middlewares.ErrorResponse(errors, response)
		return
	}
	collection := c.MG.Database("notes").Collection("users")
	update := bson.M{"$set": data}
	_, err = collection.UpdateByID(context.TODO(), userId, update)
	if err != nil {
		errors.Code = 100
		errors.Message = err.Error()
		middlewares.ErrorResponse(errors, response)
		return
	}
	middlewares.SuccessResponse("User updated", response)
}

// UserProfileDeleteEndpoint godoc
// @Summary      User delete
// @Description  Delete user profile and notes information
// @Tags         User
// @Accept       x-www-form-urlencoded
// @Produce      json
// @Success      200  {object}  models.UniversalDTO "ok"
// @Failure      400  {object}  models.UniversalDTO "error"
// @Failure      404  {object}  models.UniversalDTO "error"
// @Failure      500  {object}  models.UniversalDTO "error"
// @Security BearerAuth
// @Router       /users/profile [delete]
func (c Controller) UserProfileDeleteEndpoint(response http.ResponseWriter, request *http.Request) {
	var errors models.Error
	userId, err := helpers.GetUserID()
	if err != nil {
		errors.Code = 120
		errors.Message = err.Error()
		middlewares.ErrorResponse(errors, response)
		return
	}
	filter := bson.M{"_id": userId}
	collection := c.MG.Database("notes").Collection("users")
	_, err = collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		errors.Code = 120
		errors.Message = err.Error()
		middlewares.ErrorResponse(errors, response)
		return
	}
	filter = bson.M{"user_id": userId}
	collection = c.MG.Database("notes").Collection("notes")
	_, err = collection.DeleteMany(context.TODO(), filter)
	if err != nil {
		errors.Code = 120
		errors.Message = err.Error()
		middlewares.ErrorResponse(errors, response)
		return
	}
	collection = c.MG.Database("notes").Collection("categories")
	_, err = collection.DeleteMany(context.TODO(), filter)
	if err != nil {
		errors.Code = 125
		errors.Message = err.Error()
		middlewares.ErrorResponse(errors, response)
		return
	}
	middlewares.SuccessResponse("User and notes deleted", response)
}

// UserPasswordUpdateEndpoint godoc
// @Summary      User password
// @Description  Update user password for authorization
// @Tags         User
// @Accept       x-www-form-urlencoded
// @Produce      json
// @Param        password   formData  string true  "Password"
// @Success      200  {object}  models.UniversalDTO "ok"
// @Failure      400  {object}  models.UniversalDTO "error"
// @Failure      404  {object}  models.UniversalDTO "error"
// @Failure      500  {object}  models.UniversalDTO "error"
// @Security BearerAuth
// @Router       /users/password [put]
func (c Controller) UserPasswordUpdateEndpoint(response http.ResponseWriter, request *http.Request) {
	var errors models.Error
	err := request.ParseForm()
	if err != nil {
		errors.Code = 130
		errors.Message = err.Error()
		middlewares.ErrorResponse(errors, response)
		return
	}
	password := request.PostFormValue("password")

	userId, err := helpers.GetUserID()
	if err != nil {
		errors.Code = 140
		errors.Message = err.Error()
		middlewares.ErrorResponse(errors, response)
		return
	}

	hash := fmt.Sprintf("%x", md5.Sum([]byte(password)))
	collection := c.MG.Database("notes").Collection("users")
	update := bson.M{"$set": bson.M{"password": hash}}
	_, err = collection.UpdateByID(context.TODO(), userId, update)
	if err != nil {
		errors.Code = 150
		errors.Message = err.Error()
		middlewares.ErrorResponse(errors, response)
		return
	}
	middlewares.SuccessResponse("Update password success", response)
}

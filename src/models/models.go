package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Version print version with start program
var Version = "development"

// UniversalDTO -> model for response json
type UniversalDTO struct {
	StatusCode int         `json:"status"`
	Success    bool        `json:"success"`
	Error      Error       `json:"error"`
	Data       interface{} `json:"data"`
} //@name UniversalDTO

// Request -> request message
type Request struct {
	Name  string `json:"name" example:"Jhon Doe"`
	Phone string `json:"phone" example:"+79000000000"`
} //@name Request

// Error -> model error
type Error struct {
	Code    int    `json:"code" example:"1"`
	Message string `json:"message" example:"Note not found."`
} //@name Error

// JWT -> model jwt
type JWT struct {
	Success     bool   `json:"success"`
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
	TokenType   string `json:"token_type"`
} //@name JWT

// User -> model user
type User struct {
	ID           primitive.ObjectID `json:"id" bson:"_id"`
	Name         string             `json:"name" bson:"name"`
	Email        string             `json:"email" bson:"email"`
	Password     string             `json:"-" bson:"password"`
	Avatar       string             `json:"avatar" bson:"avatar"`
	Active       bool               `json:"active" bson:"active"`
	RestoreToken string             `json:"restore_token" bson:"restore_token"`
	CreatedAt    time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at" bson:"updated_at"`
	AuthorizedAt time.Time          `json:"authorized_at" bson:"authorized_at"`
} //@name User

// Note -> model note
type Note struct {
	ID         primitive.ObjectID `json:"id" bson:"_id"`
	CategoryID primitive.ObjectID `json:"category_id" bson:"category_id"`
	UserID     primitive.ObjectID `json:"user_id" bson:"user_id"`
	Title      string             `json:"title" bson:"title"`
	Note       string             `json:"note" bson:"note"`
	Status     string             `json:"status" bson:"status"`
	CreatedAt  time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt  time.Time          `json:"updated_at" bson:"updated_at"`
} //@name Note

// Notes -> model notes
type Notes struct {
	Total   int64   `json:"total"`
	Page    int64   `json:"page"`
	PerPage int64   `json:"per_page"`
	Items   []*Note `json:"items"`
} //@name Notes

// Category -> model category
type Category struct {
	ID       primitive.ObjectID `json:"id" bson:"_id"`
	ParentID primitive.ObjectID `json:"parent_id" bson:"parent_id"`
	UserID   primitive.ObjectID `json:"user_id" bson:"user_id"`
	Name     string             `json:"name" bson:"name"`
	Sort     int                `json:"sort" bson:"sort"`
} //@name Category

// Categories -> model categories
type Categories struct {
	Total int         `json:"total"`
	Items []*Category `json:"items"`
} //@name Categories

package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"orlangur.link/services/mini.note/models"
	"orlangur.link/services/mini.note/routes"
	"strconv"
	"strings"
	"testing"
)

var (
	email           = "test@test.loc"
	password        = "12345678"
	name            = "Test"
	newPassword     = "12345678"
	categoryName    = "Parent Category"
	categoryUpdName = "Parent Category update"
	categorySort    = "10"
	categoryUpdSort = "20"
	noteTitle       = "Test title"
	noteUpdTitle    = "Test title update"
	noteDesc        = "Description test"
	noteUpdDesc     = "Description test update"
	noteStatus      = "public"
	noteUpdStatus   = "draft"
	noteID          = ""
	categoryID      = ""
)

func TestRegisterEndpoint(t *testing.T) {
	body := url.Values{}
	body.Set("email", email)
	body.Set("password", password)
	req, err := http.NewRequest("POST", "/api/v1/users/register", strings.NewReader(body.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(body.Encode())))
	rr := httptest.NewRecorder()
	client := routes.Routes()
	client.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	if !strings.Contains(rr.Body.String(), "success") {
		t.Errorf("handler returned unexpected body: got %v",
			rr.Body.String())
	}
}

func TestUserLoginEndpoint(t *testing.T) {
	_, err := getJWTToken()
	if err != nil {
		t.Error(err)
	}
	/*var token models.JWT
	body := url.Values{}
	body.Set("email", email)
	body.Set("password", password)
	req, err := http.NewRequest("POST", "/api/v1/users/login", strings.NewReader(body.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(body.Encode())))
	rr := httptest.NewRecorder()
	client := routes.Routes()
	client.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	err = json.NewDecoder(rr.Body).Decode(&token)
	if err != nil {
		t.Errorf("error decode response: got %v",
			rr.Body.String())
	}
	if !token.Success {
		t.Errorf("handler returned unexpected body: got %v",
			rr.Body.String())
	}*/
}

func TestUserForgotEndpoint(t *testing.T) {
	var object interface{}
	body := url.Values{}
	body.Set("email", email)
	req, err := http.NewRequest("POST", "/api/v1/users/forgot", strings.NewReader(body.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(body.Encode())))
	rr := httptest.NewRecorder()
	client := routes.Routes()
	client.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	err = json.NewDecoder(rr.Body).Decode(&object)
	if err != nil {
		t.Fatal(err)
	}
	decode := object.(map[string]interface{})
	data := decode["data"].(map[string]interface{})
	if !decode["success"].(bool) {
		t.Errorf("handler returned unexpected body: got %v",
			rr.Body.String())
	}
	body.Set("restore_token", data["token"].(string))
	body.Set("password", password)
	req, err = http.NewRequest("POST", "/api/v1/users/forgot", strings.NewReader(body.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(body.Encode())))
	rr = httptest.NewRecorder()
	client.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	if !strings.Contains(rr.Body.String(), "success") {
		t.Errorf("handler returned unexpected body: got %v",
			rr.Body.String())
	}
}

func TestUserCheckByEmailEndpoint(t *testing.T) {
	var response models.UniversalDTO
	body := url.Values{}
	body.Set("email", email)
	req, err := http.NewRequest("POST", "/api/v1/users/check", strings.NewReader(body.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(body.Encode())))
	rr := httptest.NewRecorder()
	client := routes.Routes()
	client.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	err = json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Fatal(err)
	}
	if !response.Success {
		t.Errorf("handler returned unexpected body: got %v", response)
	}
}

func TestUserProfileReadEndpoint(t *testing.T) {
	var token models.JWT
	body := url.Values{}
	body.Set("email", email)
	body.Set("password", password)
	req, err := http.NewRequest("POST", "/api/v1/users/login", strings.NewReader(body.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(body.Encode())))
	rr := httptest.NewRecorder()
	client := routes.Routes()
	client.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	err = json.NewDecoder(rr.Body).Decode(&token)
	if err != nil {
		t.Errorf("error decode response: got %v",
			rr.Body.String())
	}
	if !token.Success {
		t.Errorf("handler returned unexpected body: got %v",
			rr.Body.String())
	}
	jwt := token.AccessToken

	req, err = http.NewRequest("GET", "/api/v1/users/profile", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Authorization", "Bearer "+jwt)
	rr = httptest.NewRecorder()
	client.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	if !strings.Contains(rr.Body.String(), "success") {
		t.Errorf("handler returned unexpected body: got %v",
			rr.Body.String())
	}
}

func TestUserProfileUpdateEndpoint(t *testing.T) {
	var response models.UniversalDTO
	jwt, err := getJWTToken()
	if err != nil {
		t.Error(err)
	}

	update := &bytes.Buffer{}
	writer := multipart.NewWriter(update)
	fw, err := writer.CreateFormField("name")
	if err != nil {
		t.Fatal(err)
	}
	_, err = io.Copy(fw, strings.NewReader(name))
	if err != nil {
		t.Fatal(err)
	}
	writer.Close()
	req, err := http.NewRequest("PUT", "/api/v1/users/profile", bytes.NewReader(update.Bytes()))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Content-Type", writer.FormDataContentType())
	req.Header.Add("Authorization", "Bearer "+jwt)
	rr := httptest.NewRecorder()
	client := routes.Routes()
	client.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	err = json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Fatal(err)
	}
	if !response.Success {
		t.Errorf("handler returned unexpected body: got %v", response)
	}
}

func TestUserPasswordUpdateEndpoint(t *testing.T) {
	var response models.UniversalDTO
	jwt, err := getJWTToken()
	if err != nil {
		t.Error(err)
	}
	/*update := &bytes.Buffer{}
	writer := multipart.NewWriter(update)
	fw, err := writer.CreateFormField("password")
	if err != nil {
		t.Fatal(err)
	}
	_, err = io.Copy(fw, strings.NewReader(newPassword))
	if err != nil {
		t.Fatal(err)
	}
	writer.Close()*/
	update := url.Values{}
	update.Set("password", newPassword)
	//req, err = http.NewRequest("PUT", "/api/v1/users/password", bytes.NewReader(update.Bytes()))
	req, err := http.NewRequest("PUT", "/api/v1/users/password", strings.NewReader(update.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	//25d55ad283aa400af464c76d713c07ad
	//req.Header.Add("Content-Type", writer.FormDataContentType())
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(update.Encode())))
	req.Header.Add("Authorization", "Bearer "+jwt)
	rr := httptest.NewRecorder()
	client := routes.Routes()
	client.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v wont %v", status, http.StatusOK)
	}

	err = json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Fatal(err)
	}

	if !response.Success {
		t.Errorf("handler returned unexpected body: got %v", response)
	}
}

func TestCategoryCreateEndpoint(t *testing.T) {
	var response models.UniversalDTO
	jwt, err := getJWTToken()
	if err != nil {
		t.Error(err)
	}
	body := url.Values{}
	body.Set("name", categoryName)
	body.Set("sort", categorySort)
	req, err := http.NewRequest("POST", "/api/v1/categories", strings.NewReader(body.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Authorization", "Bearer "+jwt)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(body.Encode())))
	rr := httptest.NewRecorder()
	client := routes.Routes()
	client.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	// Check the response body is what we expect.
	err = json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Fatal(err)
	}
	if !response.Success {
		t.Errorf("handler returned unexpected body: got %v", response)
	}
	data := response.Data.(map[string]interface{})
	categoryID = data["id"].(string)
}

func TestCategoryUpdateEndpoint(t *testing.T) {
	var response models.UniversalDTO
	jwt, err := getJWTToken()
	if err != nil {
		t.Error(err)
	}
	body := url.Values{}
	body.Set("name", categoryUpdName)
	body.Set("sort", categoryUpdSort)
	req, err := http.NewRequest("PUT", "/api/v1/categories/"+categoryID, strings.NewReader(body.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Authorization", "Bearer "+jwt)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(body.Encode())))
	rr := httptest.NewRecorder()
	client := routes.Routes()
	client.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	// Check the response body is what we expect.
	err = json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Fatal(err)
	}
	if !response.Success {
		t.Errorf("handler returned unexpected body: got %v", response)
	}
}

func TestCategoryReadEndpoint(t *testing.T) {
	var response models.UniversalDTO
	jwt, err := getJWTToken()
	if err != nil {
		t.Error(err)
	}
	req, err := http.NewRequest("GET", "/api/v1/categories/"+categoryID, nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Authorization", "Bearer "+jwt)
	rr := httptest.NewRecorder()
	client := routes.Routes()
	client.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	// Check the response body is what we expect.
	err = json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Fatal(err)
	}
	if !response.Success {
		t.Errorf("handler returned unexpected body: got %v", response)
	}
}

func TestCategoryListEndpoint(t *testing.T) {
	var response models.UniversalDTO
	jwt, err := getJWTToken()
	if err != nil {
		t.Error(err)
	}
	req, err := http.NewRequest("GET", "/api/v1/categories", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Authorization", "Bearer "+jwt)
	rr := httptest.NewRecorder()
	client := routes.Routes()
	client.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	// Check the response body is what we expect.
	err = json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Fatal(err)
	}
	if !response.Success {
		t.Errorf("handler returned unexpected body: got %v", response)
	}
}

func TestCategoryDeleteEndpoint(t *testing.T) {
	var response models.UniversalDTO
	jwt, err := getJWTToken()
	if err != nil {
		t.Error(err)
	}
	req, err := http.NewRequest("DELETE", "/api/v1/categories/"+categoryID, nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Authorization", "Bearer "+jwt)
	rr := httptest.NewRecorder()
	client := routes.Routes()
	client.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	// Check the response body is what we expect.
	err = json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Fatal(err)
	}
	if !response.Success {
		t.Errorf("handler returned unexpected body: got %v", response)
	}
}

func TestNoteCreateEndpoint(t *testing.T) {
	var response models.UniversalDTO
	jwt, err := getJWTToken()
	if err != nil {
		t.Error(err)
	}
	body := url.Values{}
	body.Set("title", noteTitle)
	body.Set("note", noteDesc)
	body.Set("status", noteStatus)
	req, err := http.NewRequest("POST", "/api/v1/notes", strings.NewReader(body.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Authorization", "Bearer "+jwt)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(body.Encode())))
	rr := httptest.NewRecorder()
	client := routes.Routes()
	client.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	// Check the response body is what we expect.
	err = json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Fatal(err)
	}
	if !response.Success {
		t.Errorf("handler returned unexpected body: got %v", response)
	}
	data := response.Data.(map[string]interface{})
	noteID = data["id"].(string)
}

func TestNoteUpdateEndpoint(t *testing.T) {
	var response models.UniversalDTO
	jwt, err := getJWTToken()
	if err != nil {
		t.Error(err)
	}
	body := url.Values{}
	body.Set("title", noteUpdTitle)
	body.Set("note", noteUpdDesc)
	body.Set("status", noteUpdStatus)
	req, err := http.NewRequest("PUT", "/api/v1/notes/"+noteID, strings.NewReader(body.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Authorization", "Bearer "+jwt)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(body.Encode())))
	rr := httptest.NewRecorder()
	client := routes.Routes()
	client.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	// Check the response body is what we expect.
	err = json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Fatal(err)
	}
	if !response.Success {
		t.Errorf("handler returned unexpected body: got %v", response)
	}
}

func TestNoteReadEndpoint(t *testing.T) {
	var response models.UniversalDTO
	jwt, err := getJWTToken()
	if err != nil {
		t.Error(err)
	}
	req, err := http.NewRequest("GET", "/api/v1/notes/"+noteID, nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Authorization", "Bearer "+jwt)
	rr := httptest.NewRecorder()
	client := routes.Routes()
	client.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	// Check the response body is what we expect.
	err = json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Fatal(err)
	}
	if !response.Success {
		t.Errorf("handler returned unexpected body: got %v", response)
	}
}

func TestNoteListEndpoint(t *testing.T) {
	var response models.UniversalDTO
	jwt, err := getJWTToken()
	if err != nil {
		t.Error(err)
	}
	req, err := http.NewRequest("GET", "/api/v1/notes", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Authorization", "Bearer "+jwt)
	rr := httptest.NewRecorder()
	client := routes.Routes()
	client.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	// Check the response body is what we expect.
	err = json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Fatal(err)
	}
	if !response.Success {
		t.Errorf("handler returned unexpected body: got %v", response)
	}
}

func TestNoteDeleteEndpoint(t *testing.T) {
	var response models.UniversalDTO
	jwt, err := getJWTToken()
	if err != nil {
		t.Error(err)
	}
	req, err := http.NewRequest("DELETE", "/api/v1/notes/"+noteID, nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Authorization", "Bearer "+jwt)
	rr := httptest.NewRecorder()
	client := routes.Routes()
	client.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	// Check the response body is what we expect.
	err = json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Fatal(err)
	}
	if !response.Success {
		t.Errorf("handler returned unexpected body: got %v", response)
	}
}

func TestUserProfileDeleteEndpoint(t *testing.T) {
	var token models.JWT
	body := url.Values{}
	body.Set("email", email)
	body.Set("password", password)
	req, err := http.NewRequest("POST", "/api/v1/users/login", strings.NewReader(body.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(body.Encode())))
	rr := httptest.NewRecorder()
	client := routes.Routes()
	client.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	err = json.NewDecoder(rr.Body).Decode(&token)
	if err != nil {
		t.Errorf("error decode response: got %v",
			rr.Body.String())
	}
	if !token.Success {
		t.Errorf("handler returned unexpected body: got %v",
			rr.Body.String())
	}
	jwt := token.AccessToken
	req, err = http.NewRequest("DELETE", "/api/v1/users/profile", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	req.Header.Add("Authorization", "Bearer "+jwt)
	client.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	if !strings.Contains(rr.Body.String(), "success") {
		t.Errorf("handler returned unexpected body: got %v",
			rr.Body.String())
	}
}

func getJWTToken() (string, error) {
	var token models.JWT
	body := url.Values{}
	body.Set("email", email)
	body.Set("password", password)
	req, err := http.NewRequest("POST", "/api/v1/users/login", strings.NewReader(body.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(body.Encode())))
	rr := httptest.NewRecorder()
	client := routes.Routes()
	client.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		sErr := fmt.Sprintf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
		err = errors.New(sErr)
		return "", err
	}
	err = json.NewDecoder(rr.Body).Decode(&token)
	if err != nil {
		sErr := fmt.Sprintf("error decode response: got %v",
			rr.Body.String())
		err = errors.New(sErr)
		return "", err
	}
	if !token.Success {
		sErr := fmt.Sprintf("handler returned unexpected body: got %v",
			rr.Body.String())
		err = errors.New(sErr)
		return "", err
	}
	jwt := token.AccessToken
	return jwt, nil
}

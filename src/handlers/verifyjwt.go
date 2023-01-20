package middlewares

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
	"orlangur.link/services/mini.note/models"
)

var mySigningKey []byte

// ParseJwt -> decode to string map JWT token
var ParseJwt map[string]interface{}

/* //Roles general role
var Roles []interface{}

//Permissions access permission
var Permissions []interface{} */

// IsAuthorized -> verify jwt header
func IsAuthorized(next http.Handler) http.Handler {
	var error models.Error
	mySigningKey = []byte(DotEnvVariable("JWT_SECRET"))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if check, err := IsCheckJWTHS256(w, r); check {
			next.ServeHTTP(w, r)
		} else {
			error.Code = 0
			error.Message = err
			AuthorizationResponse(error, w)
		}
	})
}

// IsCheckJWTHS256 -> check JWT token algorithm HS256
func IsCheckJWTHS256(response http.ResponseWriter, request *http.Request) (bool, string) {
	if request.Header["Authorization"] != nil {
		aToken := strings.Split(request.Header["Authorization"][0], " ")
		if len(aToken) > 1 {
			token, err := jwt.Parse(aToken[1], func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return false, fmt.Errorf("unexpected siging method: %v", token.Header["alg"])
				}
				return mySigningKey, nil
			})

			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				ParseJwt = claims
				return true, ""
			}
			return false, err.Error()
		}
		return false, "Invalid JWT token"
	} else {
		return false, "Not authorized"
	}
}

// IsCheckJWTRS256 -> check JWT token algorithm RS256
func IsCheckJWTRS256(response http.ResponseWriter, request *http.Request) (bool, string) {
	SecretKey := "-----BEGIN CERTIFICATE-----\n" + DotEnvVariable("JWT_SECRET") + "\n-----END CERTIFICATE-----"

	if request.Header["Authorization"] != nil {

		key, err := jwt.ParseRSAPublicKeyFromPEM([]byte(SecretKey))
		if err != nil {
			fmt.Println(err)
			return false, "Not Authorized"
		}

		aToken := strings.Split(request.Header["Authorization"][0], " ")
		if len(aToken) > 1 {
			token, err := jwt.Parse(aToken[1], func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
					return nil, fmt.Errorf("there was an error")
				}
				return key, nil
			})

			if err != nil {
				return false, err.Error()
			}

			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				/* 				realmAccess := claims["realm_access"].(map[string]interface{})
				   				Roles = realmAccess["roles"].([]interface{})
				   				resourceAccess := claims["resource_access"].(map[string]interface{})
				   				resourcePermissions := resourceAccess["react-backend"].(map[string]interface{})
				   				Permissions = resourcePermissions["roles"].([]interface{}) */
				ParseJwt = claims
				return true, ""
			}
		} else {
			//AuthorizationResponse("Invalid JWT token", response)
			return false, "Invalid JWT token"
		}
	} else {
		//AuthorizationResponse("Not Authorized", response)
		return false, "Not Authorized"
	}
	return false, ""
}

// GenerateJWT -> generate jwt
func GenerateJWT(user models.User) (models.JWT, error) {
	var JWT models.JWT

	mySigningKey = []byte(DotEnvVariable("JWT_SECRET"))
	minute, _ := strconv.Atoi(DotEnvVariable("JWT_LIFETIME"))
	expirationTime := time.Now().Add(time.Minute * time.Duration(minute)).Unix()

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["user"] = user
	claims["exp"] = expirationTime

	tokenString, err := token.SignedString(mySigningKey)

	if err != nil {
		return JWT, err
	}

	JWT.AccessToken = tokenString
	JWT.ExpiresIn = expirationTime
	JWT.TokenType = "Bearer"
	JWT.Success = true

	return JWT, nil
}

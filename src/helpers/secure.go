package helpers

import (
	"math/rand"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	middlewares "orlangur.link/services/mini.note/handlers"
)

// GetUserID -> return convert object id
func GetUserID() (primitive.ObjectID, error) {
	if middlewares.ParseJwt["user"] != nil {
		userObject := middlewares.ParseJwt["user"].(map[string]interface{})
		userId, err := primitive.ObjectIDFromHex(userObject["id"].(string))
		return userId, err
	} else {
		return primitive.ObjectID{}, nil
	}
}

// RandomString -> generate random string
func RandomString(n int) string {
	rand.Seed(time.Now().Unix())
	var alphabet = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789!@#$%^&*()_-+")
	var sb strings.Builder
	alphabetSize := len(alphabet)
	for i := 0; i < n; i++ {
		ch := alphabet[rand.Intn(alphabetSize)]
		sb.WriteRune(ch)
	}
	s := sb.String()
	return s
}

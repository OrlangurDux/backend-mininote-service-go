package helpers

import (
	"errors"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	middlewares "orlangur.link/services/mini.note/handlers"
)

var allowedTypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"image/gif":  true,
	"image/bmp":  true,
	"image/webp": true,
}

// UploadAvatar -> upload avatar for user
func UploadAvatar(file multipart.File, handler *multipart.FileHeader) (string, error) {
	execDir, _ := os.Getwd()
	avatarDir := middlewares.DotEnvVariable("AVATAR_DIR", "/uploaded/avatars/")
	pathAvatarDir := execDir + avatarDir
	defer func() {
		_, _ = file.Seek(0, io.SeekCurrent)
	}()

	buffer := make([]byte, 512)
	n, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return "", err
	}
	contentType := http.DetectContentType(buffer[:n])
	if _, ok := allowedTypes[contentType]; !ok {
		return "", errors.New("invalid file type")
	}

	if _, err := os.Stat(pathAvatarDir); os.IsNotExist(err) {
		err := os.MkdirAll(pathAvatarDir, os.ModePerm)
		if err != nil {
			return "", err
		}
	}
	userID, err := GetUserID()
	if err != nil {
		return "", err
	}
	ext := filepath.Ext(handler.Filename)
	fileAvatar := userID.Hex() + ext
	shortFileAvatar := avatarDir + fileAvatar
	fullFileAvatar := pathAvatarDir + fileAvatar
	f, err := os.OpenFile(fullFileAvatar, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return "", err
	}

	defer func() {
		if err := f.Close(); err != nil {
			log.Println(err)
		}
	}()
	_, _ = io.Copy(f, file)

	return shortFileAvatar, nil
}

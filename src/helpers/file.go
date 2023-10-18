package helpers

import (
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	middlewares "orlangur.link/services/mini.note/handlers"
)

// UploadAvatar -> upload avatar for user
func UploadAvatar(file multipart.File, handler *multipart.FileHeader) (string, error) {
	execDir, _ := os.Getwd()
	avatarDir := middlewares.DotEnvVariable("AVATAR_DIR", "/uploaded/avatars/")
	pathAvatarDir := execDir + avatarDir
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
	defer f.Close()
	_, _ = io.Copy(f, file)

	return shortFileAvatar, nil
}

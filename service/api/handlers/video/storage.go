package video

import (
	"errors"
	"first/pkg/constants"
	"first/pkg/upload"
	"mime/multipart"
	"os"
)

type defaultStorage struct {
}

func (svc defaultStorage) UploadFile(fileHeader *multipart.FileHeader) (*FileInfo, error) {
	fileName := upload.GetFileName(fileHeader.Filename)

	uploadSavePath := upload.GetSavePath()
	if upload.CheckSavePath(uploadSavePath) {
		err := upload.CreateSavePath(uploadSavePath, os.ModePerm)
		if err != nil {
			return nil, errors.New("failed to create save directory.")
		}
	}

	if upload.CheckPermission(uploadSavePath) {
		return nil, errors.New("insufficient file permissions.")
	}

	dst := uploadSavePath + "/" + fileName

	err := upload.SaveFile(fileHeader, dst)
	if err != nil {
		return nil, err
	}
	accessUrl := constants.UploadServerUrl + "/" + fileName
	return &FileInfo{AccessUrl: accessUrl}, nil
}

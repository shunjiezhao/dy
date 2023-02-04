package storage

import (
	"first/pkg/upload"
	"mime/multipart"
	"os"
	"time"
)

type defaultFileFactory struct {
}

func (f defaultFileFactory) Factory() Storage {
	return defaultFileStorage{}

}

type defaultFileStorage struct {
	UploadServerUrl string
}

func (svc defaultFileStorage) UploadFile(title string, fileHeader *multipart.FileHeader, i int64, time time.Time) {
	fileName := upload.GetFileName(fileHeader.Filename)
	ext := upload.GetFileExt(fileName)

	uploadSavePath := upload.GetSavePath()
	if upload.CheckSavePath(uploadSavePath) {
		err := upload.CreateSavePath(uploadSavePath, os.ModePerm)
		if err != nil {
			return
		}
	}

	if upload.CheckPermission(uploadSavePath) {
		return
	}

	dst := uploadSavePath + "/" + title + "." + ext

	err := upload.SaveFile(fileHeader, dst)
	if err != nil {
		return
	}
	//accessUrl := constants.UploadServerUrl + "/" + dst
	return
}

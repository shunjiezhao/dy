package storage

import (
	"first/pkg/upload"
	"os"
)

type defaultFileFactory struct {
}

func (f defaultFileFactory) Factory() Storage {
	return defaultFileStorage{}

}

type defaultFileStorage struct {
	UploadServerUrl string
}

func (svc defaultFileStorage) UploadFile(info *Info) (string, string) {
	title := info.Title
	fileName := upload.GetFileName(title)
	ext := upload.GetFileExt(fileName)

	uploadSavePath := upload.GetSavePath()
	if upload.CheckSavePath(uploadSavePath) {
		err := upload.CreateSavePath(uploadSavePath, os.ModePerm)
		if err != nil {
			return "", ""
		}
	}

	if upload.CheckPermission(uploadSavePath) {
		return "", ""

	}

	dst := uploadSavePath + "/" + title + "." + ext

	err := upload.SaveFile(info.Data, dst)
	if err != nil {
		return "", ""

	}
	//accessUrl := constants.UploadServerUrl + "/" + dst
	return "", ""
}

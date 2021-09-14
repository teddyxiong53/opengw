package controller

import (
	"archive/zip"
	"goAdapter/httpServer/model"
	"goAdapter/pkg/mylog"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func fileExist(path string) bool {
	_, err := os.Lstat(path)
	return !os.IsNotExist(err)
}

func unZip(zipFile string, destDir string) error {

	zipReader, err := zip.OpenReader(zipFile)
	if err != nil {
		mylog.Logger.Errorf("OpenReader err,", err)
		return err
	}
	defer zipReader.Close()

	for _, f := range zipReader.File {
		fpath := filepath.Join(destDir, f.Name)
		log.Println("fpath ", fpath)
		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
		} else {
			if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
				log.Println("mkdir err", err)
				return err
			}
			inFile, err := f.Open()
			if err != nil {
				log.Println("open err,", err)
				return err
			}
			defer inFile.Close()

			outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				log.Println("openFile err,", err)
				return err
			}
			defer outFile.Close()

			_, err = io.Copy(outFile, inFile)
			if err != nil {
				log.Println("copy err,", err)
				return err
			}
		}
	}

	return nil
}

func UpdatePlugin(context *gin.Context) {

	// 获取文件头
	file, err := context.FormFile("file")
	if err != nil {
		context.JSON(http.StatusOK, model.Response{
			Code:    "1",
			Message: "",
			Data:    "",
		})

		return
	}
	// 获取文件名
	exeCurDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	fileDir := exeCurDir + "/plugin"
	fileName := fileDir + "/" + file.Filename
	log.Println(fileName)

	if fileExist(fileDir) == false {
		os.MkdirAll(fileDir, os.ModePerm)
	}

	//保存文件到服务器本地
	if err := context.SaveUploadedFile(file, fileName); err != nil {
		log.Println(err)
		context.JSON(http.StatusOK, model.Response{
			Code:    "1",
			Message: "save error",
		})
		return
	}

	unZip(fileName, fileDir)
	err = os.Remove(fileName)
	if err != nil {
		mylog.Logger.Errorf("removeFile err,%s\n", fileName)
	}

	context.JSON(http.StatusOK, model.Response{
		Code:    "0",
		Message: "save sucess",
	})
}

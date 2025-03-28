/*
@Description: This is auto comment by koroFileHeader.
@Author: Linn
@Date: 2021-09-10 09:28:15
@LastEditors: WalkMiao
@LastEditTime: 2021-09-14 14:26:18
@FilePath: /goAdapter-Raw/pkg/backup/recover.go
*/
package backup

import (
	"archive/zip"
	"goAdapter/pkg/mylog"
	"io"
	"log"
	"os"
	"path/filepath"
)

// Unzip解压缩zip文件到指定的目录destDir
func Unzip(zipFile string, destDir string) error {
	zipReader, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}
	defer zipReader.Close()

	for _, f := range zipReader.File {
		fpath := filepath.Join(destDir, f.Name)
		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
		} else {
			if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
				return err
			}

			inFile, err := f.Open()
			if err != nil {
				return err
			}
			defer inFile.Close()

			outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer outFile.Close()

			_, err = io.Copy(outFile, inFile)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func UnZipFiles(zipFile string, destDir string) error {

	zipReader, err := zip.OpenReader(zipFile)
	if err != nil {
		mylog.Logger.Errorf("OpenReader err,%v", err)
		return err
	}
	defer zipReader.Close()

	for _, f := range zipReader.File {
		fpath := filepath.Join(destDir, f.Name)
		//log.Println("fpath ", fpath)
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

func RecoverFiles(name string) bool {

	exeCurDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))

	fileName := exeCurDir + "/selfpara/" + name
	fileAbsoluteDir := exeCurDir + "/"
	mylog.Logger.Debugf("fileName %v", fileName)
	if err := UnZipFiles(fileName, fileAbsoluteDir); err != nil {
		mylog.Logger.Errorf("err %v", err)
		return false
	}
	err := os.Remove(fileName)
	if err != nil {
		log.Printf("removeFile err,%s\n", fileName)
	}

	return true
}

package backup

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Zip备份指定的文件或者目录到zip文件 支持多源
// srcFiles 要备份的文件或者目录路径 可以是一个也可以是多个 可以是相对路径也可以是绝对路径
// destZip 指定的zip文件路径
func Zip(destZip string, srcFiles ...string) error {
	zipfile, err := os.Create(destZip)
	if err != nil {
		return err
	}
	defer zipfile.Close()

	archive := zip.NewWriter(zipfile)
	defer archive.Close()

	for _, srcFile := range srcFiles {
		err := filepath.Walk(srcFile, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			header, err := zip.FileInfoHeader(info)
			if err != nil {
				return err
			}

			header.Name = strings.TrimPrefix(path, filepath.Dir(srcFile)+"/")
			// header.Name = path
			if info.IsDir() {
				header.Name += "/"
			} else {
				header.Method = zip.Deflate
			}

			writer, err := archive.CreateHeader(header)
			if err != nil {
				return err
			}

			if !info.IsDir() {
				file, err := os.Open(path)
				if err != nil {
					return err
				}
				defer file.Close()
				_, err = io.Copy(writer, file)
				if err != nil {
					return err
				}
			}
			return err
		})
		if err != nil {
			return err
		}
	}
	return err
}

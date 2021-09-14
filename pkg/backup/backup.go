package backup

import (
	"archive/zip"
	"goAdapter/pkg/mylog"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

//压缩多个文件到一个文件里面
//Param 1: 输出的zip文件的名字
//Param 2: 需要添加到zip文件里面的文件
//Param 3: 由于file是绝对路径，打包后可能不是想要的目录，oldform就是filename中需要被替换的掉的路径
//Param 4: 要替换成的路径
func ZipFiles(filename string, files []string, oldform, newform string) error {

	newZipFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer newZipFile.Close()

	zipWriter := zip.NewWriter(newZipFile)
	defer zipWriter.Close()

	// 把files添加到zip中
	for _, file := range files {

		zipfile, err := os.Open(file)
		if err != nil {
			return err
		}
		defer zipfile.Close()

		// 获取file的基础信息
		info, err := zipfile.Stat()
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		//使用上面的FileInforHeader() 就可以把文件保存的路径替换成我们自己想要的了，如下面
		header.Name = strings.Replace(file, oldform, newform, -1)

		// 优化压缩
		// 更多参考see http://golang.org/pkg/archive/zip/#pkg-constants
		header.Method = zip.Deflate

		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}
		if _, err = io.Copy(writer, zipfile); err != nil {
			return err
		}
	}
	return nil
}

func updataConfigFile(path string, fileName []string) ([]string, error) {

	rd, err := ioutil.ReadDir(path)
	if err != nil {
		log.Println("readDir err,", err)
		return fileName, err
	}

	for _, fi := range rd {
		if fi.IsDir() {
			fullDir := path + "/" + fi.Name()
			fileName, _ = updataConfigFile(fullDir, fileName)
		} else {
			fullName := path + "/" + fi.Name()
			if strings.Contains(fi.Name(), ".json") {
				//log.Println("fullName ",fullName)
				fileName = append(fileName, fullName)
			} else if strings.Contains(fi.Name(), ".lua") {
				//log.Println("fullName ",fullName)
				fileName = append(fileName, fullName)
			}
		}
	}

	return fileName, nil
}

func BackupFiles() (bool, string) {

	exeCurDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))

	//遍历json文件
	configPath := exeCurDir + "/selfpara"
	fileNameMap := make([]string, 0)
	fileNameMap, _ = updataConfigFile(configPath, fileNameMap)

	/*
		fileList := []string{
			exeCurDir + "/config/collInterface.json",
			exeCurDir + "/config/commSerialInterface.json",
			exeCurDir + "/config/deviceNodeType.json",
			exeCurDir + "/config/networkParam.json",
			exeCurDir + "/config/ntpHostAddr.json",
			exeCurDir + "/config/reportServiceParamListAliyun.json",
			//exeCurDir + "/config/serverConfig.json",
		}
	*/

	//保留原来文件的结构
	err := ZipFiles(exeCurDir+"/selfpara/selfpara.zip", fileNameMap, exeCurDir+"/selfpara", "selfpara/")
	if err != nil {
		mylog.Logger.Errorf("zipFile err,%v", err)
		return false, ""
	}

	return true, exeCurDir + "/selfpara/selfpara.zip"
}

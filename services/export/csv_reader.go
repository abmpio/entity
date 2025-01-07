package export

import (
	"encoding/csv"
	"errors"
	"log"
	"os"
	"path/filepath"
)

type csvReader struct {
	Path string
}

// 创建一个新的csvReader对象
func newCSVReader(path string) (*csvReader, error) {
	if len(path) <= 0 {
		return nil, errors.New("path参数不能为空")
	}
	return &csvReader{
		Path: path,
	}, nil
}

func (h *csvReader) ReadData() (columns []string, rows [][]string, err error) {
	openCast, err := os.Open(h.Path)
	if err != nil {
		log.Fatalln("无法打开csv文件,文件路径:", h.Path)
		return
	}

	defer openCast.Close()
	//创建csv读取实例
	csvReader := csv.NewReader(openCast)

	//获取一行内容，第一行为header
	columns, err = csvReader.Read()
	if err != nil {
		return nil, nil, err
	}

	rows, err = csvReader.ReadAll()
	if err != nil {
		return columns, nil, err
	}
	return
}

func directoryIsExist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}

// 确保目录已经存在，如果不存在，则创建
func ensureDirectoryExist(path string) error {
	currentDirectory := filepath.Dir(path)
	exist, err := directoryIsExist(currentDirectory)
	if err != nil {
		return err
	}
	if !exist {
		//不存在，则创建
		err = os.MkdirAll(currentDirectory, os.ModePerm)
	}
	if err != nil {
		return err
	}
	return nil
}

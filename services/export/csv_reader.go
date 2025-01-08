package export

import (
	"encoding/csv"
	"errors"
	"log"
	"os"
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

func (h *csvReader) ReadData() ([]map[string]string, error) {
	list := make([]map[string]string, 0)
	csvFile, err := os.Open(h.Path)
	if err != nil {
		log.Fatalln("无法打开csv文件,文件路径:", h.Path)
		return list, err
	}

	defer csvFile.Close()
	//创建csv读取实例
	csvReader := csv.NewReader(csvFile)

	//获取一行内容，第一行为header
	columns, err := csvReader.Read()
	if err != nil {
		return list, err
	}

	if len(columns) <= 0 {
		return list, nil
	}
	rows, err := csvReader.ReadAll()
	if err != nil {
		return list, err
	}
	if len(rows) <= 0 {
		return list, nil
	}
	for i := range rows {
		currentRecord := make(map[string]string)
		for j, cName := range columns {
			if len(rows[i][j]) > 0 {
				currentRecord[cName] = rows[i][j]
			}
		}
		if len(currentRecord) > 0 {
			list = append(list, currentRecord)
		}
	}
	return list, nil
}

package export

import (
	"errors"
	"fmt"
	"time"

	"github.com/ReneKroon/ttlcache"
	"github.com/abmpio/abmp/pkg/log"
	jsonUtil "github.com/abmpio/libx/json"
	"github.com/abmpio/mongodbr"
	uuid "github.com/satori/go.uuid"
)

// 数据源类型
type DataSourceType string

const (
	DataSourceType_CSV DataSourceType = "csv"
)

type ImportOptions struct {
	Type  DataSourceType `json:"type"`
	Async bool           `json:"async"`

	FilePath string `json:"filePath"`
	// 导入数据源的columnName与目标对象的name的映射关系
	// key为导入数据源的columnName,如csv中第一行的column
	// value为目标对象的字段名
	FieldNameTitleMap map[string]string `json:"fieldNameTitleMap"`
	// 用来将原始的行转换到对象的属性中
	TransformRecordFunc func(src map[string]string, dest map[string]interface{}) error
}

func (o *ImportOptions) SafeFieldNameMap(fieldName string) string {
	if len(o.FieldNameTitleMap) <= 0 {
		return ""
	}
	v, ok := o.FieldNameTitleMap[fieldName]
	if !ok {
		return ""
	}
	return v
}

const (
	ImportStatus_Error    = "error"
	ImportStatus_Running  = "running"
	ImportStatus_Finished = "finished"
)

type EntityImport struct {
	ImportOptions

	Id        string     `json:"id"`
	Status    string     `json:"status"`
	StartTime time.Time  `json:"startTime"`
	EndTime   *time.Time `json:"endTime"`
}

type IEntityImportService[TEntity mongodbr.IEntity] interface {
	ImportFrom(options *ImportOptions) (importId string, err error)
	GetImport(importId string) (*EntityImport, error)

	SaveListFunc() func(list []*TEntity) error
}

var _ IEntityImportService[*mongodbr.Entity] = (*EntityImportService[*mongodbr.Entity])(nil)

type EntityImportService[T mongodbr.IEntity] struct {
	cache        *ttlcache.Cache
	saveListFunc func(list []*T) error
}

func NewEntityImportService[T mongodbr.IEntity](saveListFunc func(list []*T) error) IEntityImportService[T] {
	s := &EntityImportService[T]{
		cache:        ttlcache.NewCache(),
		saveListFunc: saveListFunc,
	}
	s.cache.SetTTL(time.Minute * 5)

	return s
}

func (s *EntityImportService[T]) SaveListFunc() func(list []*T) error {
	return s.saveListFunc
}

func (s *EntityImportService[T]) ImportFrom(options *ImportOptions) (importId string, err error) {
	importId = s.generateId()
	entityImport := &EntityImport{
		ImportOptions: *options,

		Id:        importId,
		Status:    ImportStatus_Running,
		StartTime: time.Now(),
	}

	s.cache.Set(importId, entityImport)
	if options.Async {
		//new threading to start export
		go func() {
			defer func() {
				if p := recover(); p != nil {
					msg := fmt.Sprint(p)
					log.Logger.Error(msg)
				}
			}()
			if options.Type == DataSourceType_CSV {
				// csv
				s.importFromCsv(entityImport)
			} else {
				log.Logger.Warn(fmt.Sprintf("不支持的导入类型:%s", options.Type))
			}
		}()
	} else {
		if options.Type == DataSourceType_CSV {
			s.importFromCsv(entityImport)
		} else {
			log.Logger.Warn(fmt.Sprintf("不支持的导入类型:%s", options.Type))
		}
	}
	return importId, nil
}

func (s *EntityImportService[T]) GetImport(importId string) (*EntityImport, error) {
	res, ok := s.cache.Get(importId)
	if !ok {
		return nil, errors.New("export not found")
	}
	export := res.(*EntityImport)
	return export, nil
}

func (s *EntityImportService[T]) importFromCsv(entityImport *EntityImport) error {
	list, err := s.extractEntityFromCSV(entityImport)
	if err != nil {
		return nil
	}
	if len(list) <= 0 {
		log.Logger.Info(fmt.Sprintf("导入的文件中没有包含任何要导入的数据,filePath:%s",
			entityImport.ImportOptions.FilePath))
		return nil
	} else {
		log.Logger.Info(fmt.Sprintf("导入的文件中包含 %d 条导入的记录,准备保存,filePath:%s",
			len(list),
			entityImport.ImportOptions.FilePath))
	}
	err = s.saveListFunc(list)
	if err != nil {
		log.Logger.Warn(fmt.Sprintf("将导入的数据列表保存到db中时出现异常,filePath:%s,err:%s",
			entityImport.ImportOptions.FilePath,
			err.Error()))
		return err
	} else {
		log.Logger.Info(fmt.Sprintf("已成功将导入的数据列表保存到db,filePath:%s",
			entityImport.ImportOptions.FilePath))
	}
	return nil
}

func (s *EntityImportService[T]) extractEntityFromCSV(entityImport *EntityImport) ([]*T, error) {
	csvReader, err := newCSVReader(entityImport.FilePath)
	if err != nil {
		return nil, err
	}
	rows, err := csvReader.ReadData()
	if err != nil {
		err = fmt.Errorf("在读取csv文件的数据时出现异常,文件:%s,异常信息:%s", csvReader.Path, err.Error())
		return nil, err
	}
	if len(rows) <= 1 {
		err = fmt.Errorf("文件中没有可导入的数据,文件:%s", csvReader.Path)
		return nil, err
	}
	var entityList []*T
	for i, row := range rows {
		newEntityMap := make(map[string]interface{})
		if entityImport.ImportOptions.TransformRecordFunc != nil {
			// used transform func
			err := entityImport.ImportOptions.TransformRecordFunc(row, newEntityMap)
			if err != nil {
				return nil, fmt.Errorf("导入的数据中包含了无效的数据,行:%d,err:%s", i+1, err.Error())
			}
		} else {
			for eachKey, eachValue := range row {
				cName := eachKey
				cMapName := entityImport.ImportOptions.SafeFieldNameMap(eachKey)
				if len(cMapName) > 0 {
					cName = cMapName
				}
				if len(cName) <= 0 {
					continue
				}
				newEntityMap[cName] = eachValue
			}
		}
		if len(newEntityMap) <= 0 {
			log.Logger.Warn(fmt.Sprintf("导入的数据中包含了空行,将忽略此行的数据,行:%d", i+1))
			continue
		}
		newEntity := new(T)
		err = jsonUtil.ConvertObjectTo(newEntityMap, newEntity)
		if err != nil {
			return nil, fmt.Errorf("无效的数据,行:%d,err:%s", i, err)
		} else {
			log.Logger.Debug(fmt.Sprintf("已成功校验导入的数据,行:%d,数据:%s", i+1, jsonUtil.ObjectToJson(newEntityMap)))
		}
		err = mongodbr.Validate(newEntity)
		if err != nil {
			err = fmt.Errorf("导入的数据中存在着无效的行,中止导入.行:%d,数据:%s,err:%s",
				i+1,
				jsonUtil.ObjectToJson(newEntityMap),
				err.Error())
			return nil, err
		}
		entityList = append(entityList, newEntity)
	}
	return entityList, nil
}

func (s *EntityImportService[T]) generateId() string {
	exportId := uuid.NewV4().String()
	return exportId
}

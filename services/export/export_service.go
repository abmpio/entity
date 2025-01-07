package export

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/ReneKroon/ttlcache"
	"github.com/abmpio/abmp/pkg/log"
	"github.com/abmpio/abmp/pkg/utils/lang"
	jsonUtil "github.com/abmpio/libx/json"
	"github.com/abmpio/libx/lang/tuple"
	"github.com/abmpio/libx/slicex"
	"github.com/abmpio/mongodbr"
	uuid "github.com/satori/go.uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ExportOptions struct {
	Type   string      `json:"type"`
	Target string      `json:"target"`
	Filter interface{} `json:"filter"`
	Async  bool        `json:"async"`

	GetFieldNameFunc  func(entity interface{}, name string) string `json:"-"`
	Skip              int                                          `json:"skip"`
	FieldNameList     []string                                     `json:"fieldNameList"`
	FieldNameTitleMap map[string]string                            `json:"fieldNameTitleMap"`
}

const (
	ExportStatus_Error    = "error"
	ExportStatus_Running  = "running"
	ExportStatus_Finished = "finished"
)

type EntityExport struct {
	ExportOptions

	Id           string     `json:"id"`
	Status       string     `json:"status"`
	StartTime    time.Time  `json:"startTime"`
	EndTime      *time.Time `json:"endTime"`
	FileName     string     `json:"fileName"`
	DownloadPath string     `json:"-"`
	Limit        int        `json:"-"`
}

type IEntityExportService[TEntity mongodbr.IEntity] interface {
	GetRepository() mongodbr.IRepository

	ExportToCSV(options ExportOptions) (exportId string, err error)
	GetExport(exportId string) (*EntityExport, error)
}

type EntityExportService[T mongodbr.IEntity] struct {
	repository mongodbr.IRepository

	cache *ttlcache.Cache
}

func NewEntityExportService[T mongodbr.IEntity](repository mongodbr.IRepository) IEntityExportService[T] {
	s := &EntityExportService[T]{
		repository: repository,
		cache:      ttlcache.NewCache(),
	}
	s.cache.SetTTL(time.Minute * 5)

	return s
}

func (s *EntityExportService[T]) GetRepository() mongodbr.IRepository {
	return s.repository
}

func (s *EntityExportService[T]) ExportToCSV(options ExportOptions) (exportId string, err error) {
	exportId = s.generateId()
	entityExport := &EntityExport{
		Id:            exportId,
		ExportOptions: options,
		Status:        ExportStatus_Running,
		StartTime:     time.Now(),
		FileName:      s.getFileName(exportId),
		DownloadPath:  s.getDownloadPath(exportId),
	}

	s.cache.Set(exportId, entityExport)
	if options.Async {
		//new threading to start export
		go func() {
			defer func() {
				if p := recover(); p != nil {
					msg := fmt.Sprint(p)
					log.Logger.Error(msg)
				}
			}()
			s.export(entityExport)
		}()
	} else {
		s.export(entityExport)
	}

	return exportId, nil
}

func (s *EntityExportService[T]) GetExport(exportId string) (*EntityExport, error) {
	res, ok := s.cache.Get(exportId)
	if !ok {
		return nil, errors.New("export not found")
	}
	export := res.(*EntityExport)
	return export, nil
}

func (s *EntityExportService[T]) export(export *EntityExport) {
	cursor := s.repository.FindByFilter(export.Filter).GetCursor()

	//csv writer
	csvWriter, csvFile, err := s.getCsvWriter(export)
	defer func() {
		if csvWriter != nil {
			csvWriter.Flush()
		}
		if csvFile != nil {
			_ = csvFile.Close()
		}
	}()
	if err != nil {
		export.Status = ExportStatus_Error
		export.EndTime = lang.NowToPtr()
		log.Logger.Error(fmt.Sprintf("export error (id: %s),err:%s", export.Id, err.Error()))
		s.cache.Set(export.Id, export)
		return
	}

	//write header
	columns, err := s.mapColumns(export)
	if err != nil {
		export.Status = ExportStatus_Error
		export.EndTime = lang.NowToPtr()
		log.Logger.Error(fmt.Sprintf("export error (id: %s),err:%s", export.Id, err.Error()))
		s.cache.Set(export.Id, export)
		return
	}
	columnTitle := slicex.ToSliceV[tuple.T2[string, string], string](columns, func(item tuple.T2[string, string]) string {
		return item.V2
	})
	err = csvWriter.Write(columnTitle)
	if err != nil {
		export.Status = ExportStatus_Error
		export.EndTime = lang.NowToPtr()
		log.Logger.Error(fmt.Sprintf("export error (id: %s),err: %s", export.Id, err.Error()))
		s.cache.Set(export.Id, export)
		return
	}
	csvWriter.Flush()

	i := 0
	for {
		i++

		err := cursor.Err()
		if err != nil {
			if err != mongo.ErrNoDocuments {
				export.Status = ExportStatus_Error
				export.EndTime = lang.NowToPtr()
				log.Logger.Error(fmt.Sprintf("export error (id: %s),err: %s", export.Id, err.Error()))
			} else {
				export.Status = ExportStatus_Finished
				export.EndTime = lang.NowToPtr()
				log.Logger.Debug(fmt.Sprintf("export finished (id: %s)", export.Id))
			}
			s.cache.Set(export.Id, export)
			return
		}

		//has data
		if !cursor.Next(context.Background()) {
			export.Status = ExportStatus_Finished
			export.EndTime = lang.NowToPtr()
			log.Logger.Debug(fmt.Sprintf("export finished (id: %s)", export.Id))
			s.cache.Set(export.Id, export)
			return
		}

		entityItem := new(T)
		err = cursor.Decode(entityItem)
		if err != nil {
			export.Status = ExportStatus_Error
			export.EndTime = lang.NowToPtr()
			log.Logger.Error(fmt.Sprintf("export error (id: %s),err: %s", export.Id, err.Error()))
			s.cache.Set(export.Id, export)
			return
		}

		columnNameList := slicex.ToSliceV[tuple.T2[string, string], string](columns, func(item tuple.T2[string, string]) string {
			return item.V1
		})
		cells := s.getRowCells(columnNameList, entityItem, &export.ExportOptions)
		err = csvWriter.Write(cells)
		if err != nil {
			export.Status = ExportStatus_Error
			export.EndTime = lang.NowToPtr()
			log.Logger.Error(fmt.Sprintf("export error (id: %s),err: %s", export.Id, err.Error()))
			s.cache.Set(export.Id, export)
			return
		}

		//flush
		if export.Limit > 0 && i >= export.Limit {
			csvWriter.Flush()
			i = 0
		}
	}
}

func (s *EntityExportService[T]) generateId() string {
	exportId := uuid.NewV4().String()
	return exportId
}

func (s *EntityExportService[T]) getExportDir() (dir string, err error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	exportDir := path.Join(wd, "_temp", "export", "csv")
	if !s.exists(exportDir) {
		err := os.MkdirAll(exportDir, 0755)
		if err != nil {
			return "", err
		}
	}
	return exportDir, nil
}

func (svc *EntityExportService[T]) getFileName(exportId string) (fileName string) {
	return exportId + "_" + time.Now().Format("20060102150405") + ".csv"
}

// getDownloadPath returns the download path for the export
// format: /wd/<tempDir>/export/<exportId>/<exportId>_<timestamp>.csv
func (s *EntityExportService[T]) getDownloadPath(exportId string) (downloadPath string) {
	exportDir, err := s.getExportDir()
	if err != nil {
		return ""
	}
	downloadPath = path.Join(exportDir, s.getFileName(exportId))
	return downloadPath
}

func (s *EntityExportService[T]) getCsvWriter(export *EntityExport) (csvWriter *csv.Writer, csvFile *os.File, err error) {
	// open file
	csvFile, err = os.Create(export.DownloadPath)
	if err != nil {
		return nil, nil, err
	}

	// create csv writer
	csvWriter = csv.NewWriter(csvFile)

	return csvWriter, csvFile, nil
}

func (s *EntityExportService[T]) mapColumns(export *EntityExport) (columns []tuple.T2[string, string], err error) {
	if len(export.FieldNameList) > 0 {
		for _, eachColumn := range export.FieldNameList {
			columnTitle := eachColumn
			if len(export.FieldNameTitleMap) > 0 {
				columnTitleMap, ok := export.FieldNameTitleMap[eachColumn]
				if ok && len(columnTitle) > 0 {
					columnTitle = columnTitleMap
				}
			}
			columns = append(columns, tuple.New2(eachColumn, columnTitle))
		}
		return columns, nil
	}

	var data []bson.M
	if err := s.repository.FindByFilter(export.Filter, mongodbr.FindOptionWithLimit(10)).All(&data); err != nil {
		return nil, err
	}

	// columns set
	columnsSet := make(map[string]bool)
	for _, d := range data {
		for k := range d {
			columnsSet[k] = true
		}
	}

	// columns
	columns = make([]tuple.T2[string, string], 0)
	for k := range columnsSet {
		columns = append(columns, tuple.New2(k, k))
	}

	return columns, nil
}

func (s *EntityExportService[T]) getRowCells(columns []string, entityItem *T, options *ExportOptions) []string {
	var cells []string
	var data bson.M
	if options.GetFieldNameFunc == nil {
		jsonUtil.ConvertObjectTo(entityItem, &data)
	}
	for _, c := range columns {
		if options.GetFieldNameFunc != nil {
			cellValue := options.GetFieldNameFunc(entityItem, c)
			cells = append(cells, cellValue)
		} else {
			cValue, ok := data[c]
			if !ok {
				cells = append(cells, "")
				continue
			}
			switch vValue := cValue.(type) {
			case string:
				cells = append(cells, vValue)
			case time.Time:
				cells = append(cells, vValue.Format("2006-01-02 15:04:05"))
			case int:
				cells = append(cells, strconv.Itoa(vValue))
			case int32:
				cells = append(cells, strconv.Itoa(int(vValue)))
			case int64:
				cells = append(cells, strconv.FormatInt(vValue, 10))
			case float32:
				cells = append(cells, strconv.FormatFloat(float64(vValue), 'f', -1, 32))
			case float64:
				cells = append(cells, strconv.FormatFloat(vValue, 'f', -1, 64))
			case bool:
				cells = append(cells, strconv.FormatBool(vValue))
			case primitive.ObjectID:
				cells = append(cells, vValue.Hex())
			case primitive.DateTime:
				cells = append(cells, vValue.Time().Format("2006-01-02 15:04:05"))
			default:
				cells = append(cells, fmt.Sprintf("%v", vValue))
			}
		}
	}
	return cells
}

func (svc *EntityExportService[T]) exists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return os.IsExist(err)
	}
	return true
}

package filter

import (
	"encoding/json"
	"reflect"
	"strings"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type FieldSpecification struct {
	TypeIsString bool
}

type GetFilterQueryOptions struct {
	SpecificationList map[string]*FieldSpecification
}

func (o *GetFilterQueryOptions) FieldIsString(fieldName string) bool {
	if len(o.SpecificationList) <= 0 {
		return false
	}
	specification := o.FieldSpecification(fieldName)
	if specification == nil {
		return false
	}
	return specification.TypeIsString
}

func (o *GetFilterQueryOptions) FieldSpecification(fieldName string) *FieldSpecification {
	specification, ok := o.SpecificationList[fieldName]
	if !ok {
		return nil
	}
	return specification
}

type GetFilterQueryOption = func(options *GetFilterQueryOptions)

// force specify field type is string
func GetFilterQueryOptionWithTypeIsString(fieldName string, fieldIsString bool) GetFilterQueryOption {
	return func(options *GetFilterQueryOptions) {
		specification := options.FieldSpecification(fieldName)
		if specification == nil {
			options.SpecificationList[fieldName] = &FieldSpecification{
				TypeIsString: fieldIsString,
			}
		} else {
			specification.TypeIsString = fieldIsString
		}
	}
}

// GetFilter Get entity.Filter
func GetFilter(getKeyFn func(key string) string, opts ...GetFilterQueryOption) (f *Filter, err error) {
	options := &GetFilterQueryOptions{
		SpecificationList: map[string]*FieldSpecification{},
	}
	for _, eachOpt := range opts {
		eachOpt(options)
	}
	// bind
	condStr := getKeyFn(FilterQueryFieldConditions)
	var conditions []*Condition
	if err := json.Unmarshal([]byte(condStr), &conditions); err != nil {
		return nil, err
	}

	// attempt to convert object id
	for i, cond := range conditions {
		v := reflect.ValueOf(cond.Value)
		switch v.Kind() {
		case reflect.String:
			item := cond.Value.(string)
			usedString := options.FieldIsString(cond.Key)
			if usedString {
				// used string
				conditions[i].Value = item
				continue
			}
			// mongodb object id
			id, err := bson.ObjectIDFromHex(item)
			if err == nil {
				conditions[i].Value = id
			} else {
				conditions[i].Value = item
			}
		case reflect.Slice, reflect.Array:
			var items []interface{}
			for i := 0; i < v.Len(); i++ {
				vItem := v.Index(i)
				item := vItem.Interface()

				// string
				stringItem, ok := item.(string)
				if ok {
					items = append(items, stringItem)
					continue
				}

				// default
				items = append(items, item)
			}
			conditions[i].Value = items
		}
	}

	return &Filter{
		IsOr:       false,
		Conditions: conditions,
	}, nil
}

// GetFilterQuery Get bson.M
func GetFilterQuery(getKeyFn func(key string) string, opts ...GetFilterQueryOption) (q map[string]interface{}, err error) {
	f, err := GetFilter(getKeyFn, opts...)
	if err != nil {
		return nil, err
	}

	if f == nil {
		return nil, nil
	}

	// TODO: implement logic OR

	return FilterToQuery(f)
}

func MustGetFilterQuery(getKeyFn func(key string) string, opts ...GetFilterQueryOption) (q map[string]interface{}) {
	q, err := GetFilterQuery(getKeyFn, opts...)
	if err != nil {
		return nil
	}
	return q
}

// GetFilterAll Get all
func GetFilterAll(getKeyFn func(key string) string) (res bool, err error) {
	resStr := getKeyFn(FilterQueryFieldAll)
	switch strings.ToUpper(resStr) {
	case "1":
		return true, nil
	case "0":
		return false, nil
	case "Y":
		return true, nil
	case "N":
		return false, nil
	case "T":
		return true, nil
	case "F":
		return false, nil
	case "TRUE":
		return true, nil
	case "FALSE":
		return false, nil
	default:
		return false, ErrorFilterInvalidOperation
	}
}

func MustGetFilterAll(getKeyFn func(key string) string) (res bool) {
	res, err := GetFilterAll(getKeyFn)
	if err != nil {
		return false
	}
	return res
}

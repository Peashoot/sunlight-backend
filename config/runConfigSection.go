package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

type RunConfigSection struct {
	ParamName  string       // 参数名称
	ParamValue interface{}  // 参数值
	ParamType  reflect.Type // 参数类型
}

func NewRunConfigSection(paramName string, dftVal interface{}) *RunConfigSection {
	section := &RunConfigSection{
		ParamName:  paramName,
		ParamValue: dftVal,
		ParamType:  reflect.TypeOf(dftVal),
	}
	section.sectionInit()
	return section
}

// createOrGetFromDB 创建或从数据库获取配置项
func (section *RunConfigSection) createOrGetFromDB() (SystemRunParamModel, error) {
	var paramConfig SystemRunParamModel
	err := MysqlDB.Where("name = ?", section.ParamName).First(&paramConfig).Error
	// 如果没有找到记录，需要初始化一条记录到数据库中
	if err == gorm.ErrRecordNotFound {
		paramConfig = SystemRunParamModel{
			Name:  section.ParamName,
			Value: fmt.Sprintf("%v", section.ParamValue),
			Kind:  section.ParamType.String(),
		}
		return paramConfig, MysqlDB.Create(&paramConfig).Error
	}
	return paramConfig, err
}

// sectionInit 初始化参数值
func (section *RunConfigSection) sectionInit() error {
	paramConfig, err := section.createOrGetFromDB()
	if err != nil {
		return err
	}
	section.convertFromString(paramConfig.Value)
	return nil
}

// Update 更新值
func (section *RunConfigSection) Update(newVal interface{}) error {
	if !section.ParamType.AssignableTo(reflect.TypeOf(newVal)) {
		return errors.New("unmatched kind")
	}
	section.ParamValue = newVal
	paramConfig, err := section.createOrGetFromDB()
	if err != nil {
		return err
	}
	paramConfig.Value = fmt.Sprintf("%v", section.ParamValue)
	return MysqlDB.Save(&paramConfig).Error
}

// convertFromString 从字符串转到具体数组
func (section *RunConfigSection) convertFromString(dbVal string) error {
	if section.ParamType.Kind() == reflect.String {
		section.ParamValue = dbVal
		return nil
	}
	if section.ParamType.Kind() == reflect.Bool {
		dbVal = strings.ToLower(dbVal)
		if dbVal == "true" || dbVal == "false" {
			section.ParamValue = dbVal == "true"
			return nil
		}
		return errors.New("unmatched value")
	}
	if len(dbVal) > 2 && (dbVal[0] == '[' || dbVal[0] == '{') {
		temp := reflect.Zero(section.ParamType).Interface()
		if err := json.Unmarshal([]byte(dbVal), &temp); err != nil {
			return err
		}
		section.ParamValue = temp
		return nil
	}
	flt64, err := strconv.ParseFloat(dbVal, 64)
	if err != nil {
		return err
	}
	switch section.ParamType.Kind() {
	case reflect.Int:
		section.ParamValue = int(flt64)
		return nil
	case reflect.Int8:
		section.ParamValue = int8(flt64)
		return nil
	case reflect.Int16:
		section.ParamValue = int16(flt64)
		return nil
	case reflect.Int32:
		section.ParamValue = int32(flt64)
		return nil
	case reflect.Int64:
		section.ParamValue = int64(flt64)
		return nil
	case reflect.Uint:
		section.ParamValue = uint(flt64)
		return nil
	case reflect.Uint8:
		section.ParamValue = uint8(flt64)
		return nil
	case reflect.Uint16:
		section.ParamValue = uint16(flt64)
		return nil
	case reflect.Uint32:
		section.ParamValue = uint32(flt64)
		return nil
	case reflect.Uint64:
		section.ParamValue = uint64(flt64)
		return nil
	case reflect.Float32:
		section.ParamValue = float32(flt64)
		return nil
	case reflect.Float64:
		section.ParamValue = flt64
		return nil
	}
	return errors.New("unmatched value")
}

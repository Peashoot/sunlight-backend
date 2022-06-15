package config

import "errors"

// RunConfiguration 运行配置集合
type RunConfiguration struct {
	innerMap map[string]*RunConfigSection
}

var RunConfigCollection = NewRunConfiguration()

// NewRunConfiguration 创建运行配置集合
func NewRunConfiguration() *RunConfiguration {
	return &RunConfiguration{
		innerMap: make(map[RunConfigName]*RunConfigSection),
	}
}

// Register 注册配置项
func Register[T any](paramName RunConfigName, defaultVal T) {
	RunConfigCollection.innerMap[paramName] = NewRunConfigSection(paramName, defaultVal)
}

// Update 更新配置项
func Update[T any](paramName RunConfigName, newVal T) error {
	if section := RunConfigCollection.innerMap[paramName]; section != nil {
		return section.Update(newVal)
	}
	panic(errors.New("not exists"))
}

// GetValue 获取配置值
func GetValue[T any](paramName string) T {
	var val interface{}
	if section := RunConfigCollection.innerMap[paramName]; section != nil {
		val = section.ParamValue
	} else {
		panic(errors.New("not exists"))
	}
	ret, _ := val.(T)
	return ret
}

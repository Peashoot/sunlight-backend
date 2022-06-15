package services

import (
	"github.com/peashoot/sunlight/config"
	"github.com/peashoot/sunlight/entity/db"
	"github.com/peashoot/sunlight/utils"
)

type EnumDicService struct {
}

func NewEnumDicService() *EnumDicService {
	return &EnumDicService{}
}

// GetAllEnumDic 获取所有字典
// 字典表根据不同的类型放到不同的键值中
func (enumDicService *EnumDicService) GetAllEnumDic(fromCache bool) ([]db.EnumDicModel, error) {
	keys, _, err := utils.RedisScan(0, utils.EnumDictionaryCachePrefix+"*", 128)
	enumModels := make([]db.EnumDicModel, 0)
	// 当缓存中存在字典表信息，并且要求从缓存取
	if fromCache && err == nil && len(keys) > 0 {
		for _, key := range keys {
			dic, err := utils.RedisGet[map[string]db.EnumDicModel](key)
			if err != nil {
				continue
			}
			for _, val := range dic {
				enumModels = append(enumModels, val)
			}
		}
		return enumModels, nil
	}
	// 不是从缓存取的，先把缓存里的全部清空
	if err == nil && len(keys) > 0 {
		utils.RedisRemove(keys...)
	}
	// 从数据库查询所有字典
	if err := config.MysqlDB.Model(&db.EnumDicModel{}).Find(&enumModels).Error; err != nil {
		return enumModels, err
	}
	// 将数据库查询结果写到缓存中
	preMap := make(map[string]map[string]db.EnumDicModel)
	for _, model := range enumModels {
		if preMap[model.TypeName] == nil {
			preMap[model.TypeName] = make(map[string]db.EnumDicModel)
		}
		preMap[model.TypeName][model.Value] = model
	}
	for enumType, enumDic := range preMap {
		utils.RedisSet(enumType, enumDic)
	}
	return enumModels, err
}

// GetEnumDic 更新缓存
func (enumDicService *EnumDicService) GetEnumDic(enumType string, fromCache bool) (map[string]db.EnumDicModel, error) {
	if fromCache {
		enumDic, err := utils.RedisGet[map[string]db.EnumDicModel](utils.EnumDictionaryCachePrefix + enumType)
		if err == nil {
			return enumDic, nil
		}
	}
	enumModels := make([]db.EnumDicModel, 0)
	if err := config.MysqlDB.Model(&db.EnumDicModel{}).Where("type_name = ?", enumType).Find(&enumModels).Error; err != nil {
		return nil, err
	}
	enumDic := make(map[string]db.EnumDicModel)
	for _, model := range enumModels {
		enumDic[model.Value] = model
	}
	return enumDic, nil
}

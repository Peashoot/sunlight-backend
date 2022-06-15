package services

import (
	"time"

	"github.com/google/uuid"
	"github.com/peashoot/sunlight/config"
	"github.com/peashoot/sunlight/entity/db"
	"github.com/peashoot/sunlight/utils"
	"gorm.io/gorm"
)

type FileStorageService struct {
}

func NewFileStorageService() FileStorageService {
	return FileStorageService{}
}

// AddRecord 保存文件
func (fileService *FileStorageService) AddRecord(name, subffix, fileURL, categoryCode, opUserCode string) (db.FileStorageRecord, error) {
	storage := db.FileStorageRecord{
		BaseModel: db.BaseModel{
			Code:      uuid.NewString(),
			CreatedBy: opUserCode,
		},
		Name:         name,
		Subffix:      subffix,
		FileURL:      fileURL,
		CategoryCode: categoryCode,
	}
	storage.FileType = utils.GetFileType(storage.Subffix)
	// 将文件信息保存到数据库中
	if err := config.MysqlDB.Create(&storage).Error; err != nil {
		return storage, err
	}
	return storage, nil
}

// Modify 修改文件
func (fileService *FileStorageService) Modify(fileCode, categoryCode, opUserCode string) (db.FileStorageRecord, error) {
	var storage db.FileStorageRecord
	if err := config.MysqlDB.Where("code = ?", fileCode).First(&storage).Error; err != nil {
		return storage, err
	}
	storage.UpdatedBy = opUserCode
	storage.CategoryCode = categoryCode
	storage.UpdatedAt = time.Now()
	err := config.MysqlDB.Save(&storage).Error
	return storage, err
}

// Remove 移除文件
func (fileService *FileStorageService) Remove(fileCode, opUserCode string) (db.FileStorageRecord, error) {
	var storage db.FileStorageRecord
	if err := config.MysqlDB.Where("code = ?", fileCode).First(&storage).Error; err != nil {
		return storage, err
	}
	storage.DeletedBy = opUserCode
	storage.DeletedAt = gorm.DeletedAt{Time: time.Now(), Valid: true}
	err := config.MysqlDB.Save(&storage).Error
	return storage, err
}

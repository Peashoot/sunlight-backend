package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	ConfigFilePath      string
	ListenPort          int
	MysqlDB             *gorm.DB
	RedisDB             *redis.Client
	AppRunLogFileWriter *os.File
)

type logWriter struct {
}

func (w logWriter) Printf(format string, args ...interface{}) {
	log.Printf(format, args...)
}

func init() {
	flag.StringVar(&ConfigFilePath, "c", "config.json", "the path of config file")
	var err error
	var startUpConfig StartUpConfiguration
	var configFileBytes []byte
	configFileBytes, err = os.ReadFile(ConfigFilePath)
	if err != nil {
		log.Fatalln("can't read file", ConfigFilePath, ", please make the path correct, err:", err)
	}
	err = json.Unmarshal(configFileBytes, &startUpConfig)
	if err != nil {
		log.Fatalln("fail to unmarshal config file, err:", err)
	}
	ListenPort = startUpConfig.ListenPort
	// 创建日志文件
	AppRunLogFileWriter, err = os.OpenFile(startUpConfig.LogPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0744)
	if err != nil {
		log.Fatalln("open log file failed, err:", err)
		return
	}
	// 连接mysql
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local&timeout=%ds",
		startUpConfig.Mysql.Username,
		startUpConfig.Mysql.Password,
		startUpConfig.Mysql.Host,
		startUpConfig.Mysql.Port,
		startUpConfig.Mysql.DBName,
		startUpConfig.Mysql.ConnectTimeout)
	gormLogger := logger.New(
		logWriter{},
		logger.Config{
			SlowThreshold:             200 * time.Millisecond, // Slow SQL threshold
			LogLevel:                  logger.Warn,            // Log level
			IgnoreRecordNotFoundError: true,                   // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,                  // Disable color
		},
	)
	MysqlDB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger:                                   gormLogger,
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		log.Fatalln("try to connect mysql failured, ex:", err)
	}
	// 连接redis
	RedisDB = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", startUpConfig.Redis.Host, startUpConfig.Redis.Port), // 指定
		Password: startUpConfig.Redis.Password,
		DB:       startUpConfig.Redis.DBIndex,
	})
	if _, err = RedisDB.Ping().Result(); err != nil {
		log.Fatal("try to connect redis failured, ex:", err)
	}
}

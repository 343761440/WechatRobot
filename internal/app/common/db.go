package common

import (
	"errors"
	"log"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/ini.v1"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func openDB(debug, singularTable bool, opt ...DBOption) (*gorm.DB, error) {
	opts := defaultDBOptions
	for _, o := range opt {
		o.apply(&opts)
	}

	logLevel := logger.Warn
	if debug {
		logrus.Info("openDB with debug mode")
		logLevel = logger.Info
	}

	logger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logLevel,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	logrus.Info("open db: ", opts.dialect)
	var db *gorm.DB
	switch opts.dialect {
	case "mysql":
		var err error
		sqlconnstr := opts.mysqlOptions.User + ":" + opts.mysqlOptions.Password + "@tcp(" + opts.mysqlOptions.Host + ":" + opts.mysqlOptions.Port + ")/" + opts.mysqlOptions.Database + "?charset=utf8mb4&parseTime=true&loc=Local"
		db, err = gorm.Open(mysql.Open(sqlconnstr), &gorm.Config{
			Logger: logger,
			NamingStrategy: schema.NamingStrategy{
				SingularTable: singularTable,
			},
		})
		db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;")
		if err != nil {
			return nil, err
		}

	case "sqlite3":
		var err error
		db, err = gorm.Open(sqlite.Open(opts.sqliteOptions.Path), &gorm.Config{
			Logger: logger,
			NamingStrategy: schema.NamingStrategy{
				SingularTable: singularTable,
			},
		})
		if err != nil {
			return nil, err
		}

	default:
		var errUnsupportDialect = errors.New("not support dialect")
		return nil, errUnsupportDialect
	}
	return db, nil
}

func InitDB(debug, singularTable bool) (*gorm.DB, error) {
	cfg, err := ini.Load(kCommonConfigPath)
	if err != nil {
		logrus.Error("Fail to read config file: ", err)
		os.Exit(1)
	}

	if mysqlCfg := cfg.Section("MYSQL"); len(mysqlCfg.Keys()) > 0 { // 优先使用MYSQL
		log.Printf("using database mysql: %+v", mysqlCfg)
		opts := []DBOption{MysqlOptions(&MySQLOptions{
			User:     mysqlCfg.Key("UserName").String(),
			Password: mysqlCfg.Key("Pwd").String(),
			Host:     mysqlCfg.Key("HostName").String(),
			Port:     mysqlCfg.Key("Port").String(),
			Database: mysqlCfg.Key("DatabaseName").String(),
		})}
		return openDB(debug, singularTable, opts...)
	}

	logrus.Error("Fail to get database config")
	os.Exit(1)
	return nil, errors.New("Fail to get database config")
}

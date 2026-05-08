package databaseHelper

import (
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var (
	Db *gorm.DB

	databaseName  = "template_db"
	parseTimeFlag = "?parseTime=true"
	connectionUrl = os.Getenv("DB") + databaseName + parseTimeFlag

	loc = time.FixedZone("", -3*60*60)
)

func InitDB() *gorm.DB {

	// Create GORM configuration
	config := gorm.Config{
		NowFunc: func() time.Time {
			timeDateNow := time.Now().In(loc)

			return timeDateNow
		},
		TranslateError: false,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	}

	db, err := gorm.Open(mysql.Open(connectionUrl), &config)
	db.Debug()
	if err != nil {
		panic(err.Error())
	}
	sqlDB, err := db.DB()
	if err != nil {
		panic(err.Error())
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)
	migrate(db)
	return db
}

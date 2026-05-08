package databaseHelper

import (
	"log"
	"os"
	"sales-manager-back/pkg/domain/agenda"
	"sales-manager-back/pkg/domain/client"
	"sales-manager-back/pkg/domain/collection"
	"sales-manager-back/pkg/domain/commission"
	"sales-manager-back/pkg/domain/dashboard"
	"sales-manager-back/pkg/domain/pipeline"
	"sales-manager-back/pkg/domain/principal"
	"sales-manager-back/pkg/domain/sale"
	"sales-manager-back/pkg/domain/tenant"
	"sales-manager-back/pkg/domain/user"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var (
	Db *gorm.DB

	databaseName  = "sales_manager"
	parseTimeFlag = "?parseTime=true"
	connectionUrl = os.Getenv("DB") + "/" + databaseName + parseTimeFlag

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
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	db.Debug()

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Failed to get database instance:", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	migrate(db)
	Db = db
	return db
}

func migrate(db *gorm.DB) {
	// Run all migrations
	err := db.AutoMigrate(
		&tenant.Tenant{},
		&user.User{},
		&principal.Principal{},
		&principal.PriceList{},
		&principal.Catalog{},
		&principal.Promotion{},
		&client.Client{},
		&client.ClientPrincipalCondition{},
		&client.AccountMovement{},
		&sale.Sale{},
		&sale.SaleItem{},
		&commission.Commission{},
		&agenda.AgendaEvent{},
		&pipeline.Deal{},
		&dashboard.Alert{},
		&dashboard.QuickNote{},
		&dashboard.Goal{},
		&collection.Collection{},
	)
	if err != nil {
		log.Fatal("Failed to run migrations:", err)
	}
	log.Println("Migrations completed successfully")
}

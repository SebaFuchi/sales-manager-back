package databaseHelper

import (
	"fmt"
	"template/pkg/domain/sample"

	"gorm.io/gorm"
)

// SeedTracker is used to keep track of which seeders have been executed
type SeedTracker struct {
	Id        uint   `gorm:"primaryKey"`
	SeederKey string `gorm:"type:varchar(255);uniqueIndex"`
	Executed  bool
}

func migrate(db *gorm.DB) {
	// First migrate the SeedTracker table
	err := db.AutoMigrate(&SeedTracker{})
	if err != nil {
		panic("Error migrating seed tracker: " + err.Error())
	}

	// Migrate base tables first (tables with no foreign keys)
	baseModels := []any{}

	err = db.AutoMigrate(baseModels...)
	if err != nil {
		panic("Error migrating base tables: " + err.Error())
	}

	// Then migrate dependent tables (tables with foreign keys)
	dependentModels := []any{
		&sample.Animals{},
	}

	err = db.AutoMigrate(dependentModels...)
	if err != nil {
		panic("Error migrating dependent tables: " + err.Error())
	}

	fmt.Println("Migration run successfully")

	// Run seeders only if they haven't been run before
	runSeedersIfNeeded(db)
	fmt.Println("Seeder check completed")
}

func seedSample(db *gorm.DB) {
	templateSeeders := []sample.Animals{
		{
			Name: "Action",
			Age:  5,
		},
	}
	for _, templateSeeder := range templateSeeders {
		db.Save(&templateSeeder)
	}
}

func runSeedersIfNeeded(db *gorm.DB) {
	if !isSeederExecuted(db, "sample_seeder_v1") {
		fmt.Println("Running sample seeder...")
		seedSample(db)
		markSeederAsExecuted(db, "sample_seeder_v1")
	}

}

func isSeederExecuted(db *gorm.DB, key string) bool {
	var tracker SeedTracker
	result := db.Where("seeder_key = ?", key).First(&tracker)
	return result.RowsAffected > 0 && tracker.Executed
}

// Mark a seeder as executed
func markSeederAsExecuted(db *gorm.DB, key string) {
	tracker := SeedTracker{
		SeederKey: key,
		Executed:  true,
	}
	db.Save(&tracker)
}

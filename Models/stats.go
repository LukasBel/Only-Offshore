package Models

import (
	"gorm.io/gorm"
	"log"
	"time"
)

type Stats struct {
	ID                uint8     `gorm:"primaryKey;autoIncrement" json:"id"`
	User              string    `json:"user"`
	FlatBench         int       `json:"flatBench"`
	InclineBench      int       `json:"inclineBench"`
	Squat             int       `json:"squat"`
	PullUps           int       `json:"pullUps"`
	WeightedPullUpMax int       `json:"weightedPullUpMax"`
	BodyWeight        int       `json:"bodyWeight"`
	CreatedAt         time.Time `json:"createdAt"`
}

func MigrateStats(db *gorm.DB) error {
	err := db.AutoMigrate(&Stats{})
	if err != nil {
		log.Fatal("Failed to Auto Migrate Stats")
		return err
	}
	return nil
}

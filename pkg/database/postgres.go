package database

import (
	"Beside-Mom-BE/configs"
	"Beside-Mom-BE/modules/entities"
	"errors"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitDB(config configs.PostgreSQL) {
	sslMode := config.SSLMode
	if sslMode == "" {
		sslMode = "disable"
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Bangkok",
		config.Host,
		config.Username,
		config.Password,
		config.Database,
		config.Port,
		sslMode,
	)

	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	_ = db.AutoMigrate(
		&entities.User{},
		&entities.Question{},
		&entities.Video{},
		&entities.Likes{},
		&entities.Kid{},
		&entities.Appointment{},
		&entities.Care{},
		&entities.OTP{},
		&entities.Quiz{},
		&entities.Evaluate{},
		&entities.History{},
		&entities.Growth{},
		&entities.Period{},
		&entities.Category{},
	)

	insertRoles()
	insertPeriods()
	insertCategories()
	log.Println("Database connection established successfully!")
}

func GetDB() *gorm.DB {
	if db == nil {
		log.Fatal("Database is not initialized")
	}

	return db
}

func insertRoles() {
	var adminRole entities.Role
	var userRole entities.Role

	if err := db.First(&adminRole, "role_name = ?", "Admin").Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			adminRole = entities.Role{RoleName: "Admin"}
			if err := db.Create(&adminRole).Error; err != nil {
				log.Fatalf("Failed to insert Admin role: %v", err)
			}

			log.Println("Admin role created successfully!")
		} else {
			log.Fatalf("Error checking Admin role: %v", err)
		}
	}

	if err := db.First(&userRole, "role_name = ?", "User").Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			userRole = entities.Role{RoleName: "User"}
			if err := db.Create(&userRole).Error; err != nil {
				log.Fatalf("Failed to insert User role: %v", err)
			}

			log.Println("User role created successfully!")
		} else {
			log.Fatalf("Error checking User role: %v", err)
		}
	}
}

func insertPeriods() {
	periodNames := []string{
		"แรกเกิด", "1 เดือน", "2 เดือน", "3 - 4 เดือน",
		"5 - 6 เดือน", "7 - 8 เดือน", "9 เดือน", "10 - 12 เดือน",
	}

	for _, name := range periodNames {
		var existing entities.Period
		if err := db.First(&existing, "period = ?", name).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				newPeriod := entities.Period{
					Period: name,
				}
				if err := db.Create(&newPeriod).Error; err != nil {
					log.Printf("Failed to insert period '%s': %v", name, err)
					continue
				}
				log.Printf("Inserted period: %s", name)
			} else {
				log.Printf("Error checking period '%s': %v", name, err)
			}
		}
	}
}

func insertCategories() {
	periodNames := []string{
		"ด้านการเคลื่อนไหว Gross Motor (GM)", "ด้านการใช้กล้ามเนื้อมัดเล็ก และสติปัญญา Fine Motor (FM)", "ด้านการเข้าใจภาษา Receptive Language (RL)",
		"ด้านการใช้ภาษา Expression Language (EL)", "ด้านการช่วยเหลือตนเองและสังคม Personal and Social (PS)",
	}

	for _, name := range periodNames {
		var existing entities.Category
		if err := db.First(&existing, "category = ?", name).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				newCategory := entities.Category{
					Category: name,
				}
				if err := db.Create(&newCategory).Error; err != nil {
					log.Printf("Failed to insert category '%s': %v", name, err)
					continue
				}
				log.Printf("Inserted category: %s", name)
			} else {
				log.Printf("Error checking category '%s': %v", name, err)
			}
		}
	}
}

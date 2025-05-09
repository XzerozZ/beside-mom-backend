package repositories

import (
	"Beside-Mom-BE/modules/entities"
	"sort"

	"gorm.io/gorm"
)

type GormGrowthRepository struct {
	db *gorm.DB
}

func NewGormGrowthRepository(db *gorm.DB) *GormGrowthRepository {
	return &GormGrowthRepository{db: db}
}

type GrowthRepository interface {
	CreateGrowth(growth *entities.Growth) (*entities.Growth, error)
	GetGrowthByID(id string) (*entities.Growth, error)
	GetAllGrowth(kidID string) ([]entities.Growth, error)
	GetLatestGrowthByKidID(kidID string, month int) (*entities.Growth, error)
	GetSummary(kidID string) ([]map[string]interface{}, error)
	UpdateGrowth(growth *entities.Growth) (*entities.Growth, error)
}

func (r *GormGrowthRepository) CreateGrowth(growth *entities.Growth) (*entities.Growth, error) {
	if err := r.db.Create(growth).Error; err != nil {
		return nil, err
	}

	return growth, nil
}

func (r *GormGrowthRepository) GetGrowthByID(id string) (*entities.Growth, error) {
	var growth *entities.Growth
	if err := r.db.Where("id = ?", id).Find(&growth).Error; err != nil {
		return nil, err
	}

	return growth, nil
}

func (r *GormGrowthRepository) GetAllGrowth(kidID string) ([]entities.Growth, error) {
	var growth []entities.Growth
	if err := r.db.Where("kid_id = ?", kidID).Order("created_at").Find(&growth).Error; err != nil {
		return nil, err
	}

	return growth, nil
}

func (r *GormGrowthRepository) GetLatestGrowthByKidID(kidID string, month int) (*entities.Growth, error) {
	var growth *entities.Growth
	if err := r.db.Where("kid_id = ? AND months = ?", kidID, month).Order("created_at desc").First(&growth).Error; err != nil {
		return nil, err
	}

	return growth, nil
}

func (r *GormGrowthRepository) UpdateGrowth(growth *entities.Growth) (*entities.Growth, error) {
	if err := r.db.Save(&growth).Error; err != nil {
		return nil, err
	}

	return growth, nil
}

func (r *GormGrowthRepository) GetSummary(kidID string) ([]map[string]interface{}, error) {
	var kid entities.Kid
	if err := r.db.First(&kid, "id = ?", kidID).Error; err != nil {
		return nil, err
	}

	var growths []entities.Growth
	if err := r.db.Where("kid_id = ?", kidID).Find(&growths).Error; err != nil {
		return nil, err
	}

	var summary []map[string]interface{}
	hasMonth0 := false

	for _, g := range growths {
		if g.Months == 0 {
			hasMonth0 = true
		}
		summary = append(summary, map[string]interface{}{
			"length": g.Length,
			"weight": g.Weight,
			"months": g.Months,
		})
	}

	if !hasMonth0 {
		summary = append(summary, map[string]interface{}{
			"length": kid.BirthLength,
			"weight": kid.BirthWeight,
			"months": 0,
		})
	}

	sort.Slice(summary, func(i, j int) bool {
		return summary[i]["months"].(int) < summary[j]["months"].(int)
	})

	return summary, nil
}

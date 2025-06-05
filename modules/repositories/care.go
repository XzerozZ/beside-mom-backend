package repositories

import (
	"Beside-Mom-BE/configs"
	"Beside-Mom-BE/modules/entities"
	"Beside-Mom-BE/pkg/utils"

	"gorm.io/gorm"
)

type GormCareRepository struct {
	db   *gorm.DB
	supa configs.Supabase
}

func NewGormCareRepository(db *gorm.DB, supa configs.Supabase) *GormCareRepository {
	return &GormCareRepository{
		db:   db,
		supa: supa,
	}
}

type CareRepository interface {
	CreateCare(care *entities.Care, asset []entities.Asset) (*entities.Care, error)
	GetCareByID(id string) (*entities.Care, error)
	GetAllCare() ([]entities.Care, error)
	UpdateCare(care *entities.Care) (*entities.Care, error)
	AddAssets(id string, asset []entities.Asset) (*entities.Care, error)
	RemoveAssets(id string, imageID *string) error
	DeleteCare(id string) error
}

func (r *GormCareRepository) CreateCare(care *entities.Care, assets []entities.Asset) (*entities.Care, error) {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(care).Error; err != nil {
			return err
		}

		for _, asset := range assets {
			if err := tx.Create(&asset).Error; err != nil {
				return err
			}
		}

		if err := tx.Model(care).Association("Assets").Append(assets); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return r.GetCareByID(care.ID)
}

func (r *GormCareRepository) GetCareByID(id string) (*entities.Care, error) {
	var care entities.Care
	if err := r.db.Preload("Assets").Where("id = ?", id).First(&care).Error; err != nil {
		return nil, err
	}

	return &care, nil
}

func (r *GormCareRepository) GetAllCare() ([]entities.Care, error) {
	var cares []entities.Care
	if err := r.db.Preload("Assets").Find(&cares).Error; err != nil {
		return nil, err
	}

	return cares, nil
}

func (r *GormCareRepository) AddAssets(id string, assets []entities.Asset) (*entities.Care, error) {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		var care entities.Care
		if err := tx.Preload("Assets").Where("id = ?", id).First(&care).Error; err != nil {
			return err
		}

		for i := range assets {
			if err := tx.Create(&assets[i]).Error; err != nil {
				return err
			}
		}

		return tx.Model(&care).Association("Assets").Append(assets)
	})

	if err != nil {
		return nil, err
	}

	return r.GetCareByID(id)
}

func (r *GormCareRepository) RemoveAssets(id string, imageID *string) error {
	var assetsToDelete []entities.Asset
	err := r.db.Transaction(func(tx *gorm.DB) error {
		var care entities.Care
		if err := tx.Preload("Assets").Where("id = ?", id).First(&care).Error; err != nil {
			return err
		}

		for _, img := range care.Assets {
			if img.ID == *imageID {
				assetsToDelete = append(assetsToDelete, img)
				break
			}
		}

		var assetsIDs []string
		for _, img := range assetsToDelete {
			if err := utils.DeleteImage(img.Link, r.supa); err != nil {
				return err
			}

			assetsIDs = append(assetsIDs, img.ID)
		}

		if err := tx.Model(&care).Association("Assets").Delete(assetsToDelete); err != nil {
			return err
		}

		if err := tx.Where("id IN ?", assetsIDs).Delete(&entities.Care{}).Error; err != nil {
			return err
		}

		if err := tx.Where("id = ?", assetsIDs).Delete(&entities.Asset{}).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (r *GormCareRepository) UpdateCare(care *entities.Care) (*entities.Care, error) {
	if err := r.db.Save(&care).Error; err != nil {
		return nil, err
	}

	return r.GetCareByID(care.ID)
}

func (r *GormCareRepository) DeleteCare(id string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var care entities.Care
		if err := tx.Preload("Assets").Where("id = ?", id).First(&care).Error; err != nil {
			return err
		}

		if err := tx.Exec("DELETE FROM care_assets WHERE care_id = ?", id).Error; err != nil {
			return err
		}

		if len(care.Assets) > 0 {
			var assetIDs []string
			for _, asset := range care.Assets {
				if err := utils.DeleteImage(asset.Link, r.supa); err != nil {
					return err
				}

				assetIDs = append(assetIDs, asset.ID)
			}

			if err := tx.Where("id IN ?", assetIDs).Delete(&entities.Asset{}).Error; err != nil {
				return err
			}
		}

		if err := tx.Delete(&care).Error; err != nil {
			return err
		}

		return nil
	})
}

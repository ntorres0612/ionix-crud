package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

type Drug struct {
	ID          uint64    `gorm:"primary_key;auto_increment" json:"id"`
	Name        string    `gorm:"size:255;not null;unique" json:"name"`
	Approved    bool      `gorm:"size:255;not null;" json:"approved"`
	MinDose     int       `gorm:"not null" json:"min_dose"`
	MaxDose     int       `gorm:"not null" json:"max_dose"`
	AvailableAt time.Time `gorm:"not null" json:"available_at"`
	CreatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (p *Drug) Prepare() {
	p.ID = 0
	p.Name = html.EscapeString(strings.TrimSpace(p.Name))
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
}

func (p *Drug) Validate() error {

	if p.Name == "" {
		return errors.New("required Title")
	}
	if p.MinDose > 1 {
		return errors.New("required Min Dose")
	}
	return nil
}

func (drug *Drug) SaveDrug(db *gorm.DB) (*Drug, error) {
	err := db.Debug().Model(&Drug{}).Create(&drug).Error
	if err != nil {
		return &Drug{}, err
	}

	return drug, nil
}

func (p *Drug) FindAllDrugs(db *gorm.DB) (*[]Drug, error) {
	drugs := []Drug{}
	err := db.Debug().Model(&Drug{}).Limit(100).Find(&drugs).Error
	if err != nil {
		return &[]Drug{}, err
	}

	return &drugs, nil
}

func (drug *Drug) FindDrugByID(db *gorm.DB, id uint64) (*Drug, error) {
	err := db.Debug().Model(&Drug{}).Where("id = ?", id).Take(&drug).Error
	if err != nil {
		return &Drug{}, err
	}

	return drug, nil
}

func (drug *Drug) UpdateADrug(db *gorm.DB) (*Drug, error) {

	err := db.Debug().Model(&Drug{}).Where("id = ?", drug.ID).Updates(Drug{
		Name:      drug.Name,
		MinDose:   drug.MinDose,
		MaxDose:   drug.MaxDose,
		UpdatedAt: time.Now()}).Error
	if err != nil {
		return &Drug{}, err
	}

	return drug, nil
}

func (drug *Drug) DeleteADrug(db *gorm.DB, id uint64) (int64, error) {

	db = db.Debug().Model(&Drug{}).Where("id = ?", id).Take(&Drug{}).Delete(&Drug{})

	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("Drug not found")
		}
		return 0, db.Error
	}
	return db.RowsAffected, nil
}

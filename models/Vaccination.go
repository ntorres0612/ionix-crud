package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

type Vaccination struct {
	ID        uint64    `gorm:"primary_key;auto_increment" json:"id"`
	Name      string    `gorm:"size:255;not null;unique" json:"name"`
	Drug      Drug      `json:"drug"`
	DrugID    uint32    `gorm:"not null" json:"drug_id"`
	Dose      int       `gorm:"not null" json:"dose"`
	Date      time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"date"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (p *Vaccination) Prepare() {
	p.ID = 0
	p.Name = html.EscapeString(strings.TrimSpace(p.Name))
	p.Drug = Drug{}
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
}

func (vaccination *Vaccination) Validate() error {

	if vaccination.Name == "" {
		return errors.New("required Title")
	}
	if vaccination.DrugID <= 0 {
		return errors.New("required Drug")
	}

	return nil
}

func (vaccination *Vaccination) SaveVaccination(db *gorm.DB) (*Vaccination, error) {

	drug := Drug{}
	err := db.Debug().Model(&Drug{}).Where("id = ?", vaccination.DrugID).Take(&drug).Error
	if err != nil {
		return &Vaccination{}, err
	}
	if drug.ID > 0 {
		if vaccination.Dose > drug.MaxDose || vaccination.Dose < drug.MinDose {
			return &Vaccination{}, errors.New("the dose is not within the range")
		}
	}

	if vaccination.Date.Before(drug.AvailableAt) {
		return &Vaccination{}, errors.New("date not allowed")
	}

	err = db.Debug().Model(&Vaccination{}).Create(&vaccination).Error
	if err != nil {
		return &Vaccination{}, err
	}

	return vaccination, nil
}

func (vaccination *Vaccination) FindAllSaveVaccinations(db *gorm.DB) (*[]Vaccination, error) {
	vaccinations := []Vaccination{}
	err := db.Debug().Model(&Vaccination{}).Limit(100).Find(&vaccinations).Error
	if err != nil {
		return &[]Vaccination{}, err
	}
	if len(vaccinations) > 0 {
		for i := range vaccinations {
			err := db.Debug().Model(&User{}).Where("id = ?", vaccinations[i].DrugID).Take(&vaccinations[i].Drug).Error
			if err != nil {
				return &[]Vaccination{}, err
			}
		}
	}
	return &vaccinations, nil
}

func (vaccination *Vaccination) FindVaccinationByID(db *gorm.DB, id uint64) (*Vaccination, error) {
	err := db.Debug().Model(&Vaccination{}).Where("id = ?", id).Take(&vaccination).Error
	if err != nil {
		return &Vaccination{}, err
	}

	return vaccination, nil
}

func (vaccination *Vaccination) UpdateAVaccination(db *gorm.DB) (*Vaccination, error) {

	drug := Drug{}
	err := db.Debug().Model(&Drug{}).Where("id = ?", vaccination.DrugID).Take(&drug).Error
	if err != nil {
		return &Vaccination{}, err
	}
	if drug.ID > 0 {
		if vaccination.Dose > drug.MaxDose || vaccination.Dose < drug.MinDose {
			return &Vaccination{}, errors.New("the dose is not within the range")
		}
	}

	if vaccination.Date.Before(drug.AvailableAt) {
		return &Vaccination{}, errors.New("date not allowed")
	}

	err = db.Debug().Model(&Vaccination{}).Where("id = ?", vaccination.ID).Updates(Vaccination{
		Name:      vaccination.Name,
		Dose:      vaccination.Dose,
		DrugID:    vaccination.DrugID,
		Date:      vaccination.Date,
		UpdatedAt: time.Now()}).Error
	if err != nil {
		return &Vaccination{}, err
	}

	return vaccination, nil
}

func (p *Vaccination) DeleteAVaccination(db *gorm.DB, id uint64) (int64, error) {

	db = db.Debug().Model(&Vaccination{}).Where("id = ? ", id).Take(&Vaccination{}).Delete(&Vaccination{})

	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("Vaccination not found")
		}
		return 0, db.Error
	}
	return db.RowsAffected, nil
}

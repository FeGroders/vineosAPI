package models

import (
	"errors"
	"html"
	"strings"
	"github.com/jinzhu/gorm"
)

type Wine struct {
	ID     	    uint64    `gorm:"primary_key;auto_increment" json:"id"`
	Name  	    string    `gorm:"size:50;not null;" json:"name"`
	Description string    `gorm:"size:200;not null;" json:"description"`
	Year        string    `gorm:"size:4;not null;" json:"year"`
	Price 	    string    `gorm:"size:100;not null;" json:"price"`
	Image 	    string    `gorm:"size:200;not null;" json:"image"`
	Available   bool      `gorm:"type:bool;default:false" json:"available"`
}

func (w *Wine) Prepare() {
	w.ID = 0
	w.Name = html.EscapeString(strings.TrimSpace(w.Name))
	w.Description = html.EscapeString(strings.TrimSpace(w.Description))
	w.Image = html.EscapeString(strings.TrimSpace(w.Image))
	// w.Year = html.EscapeString(strings.TrimSpace(w.Year))
}

func (w *Wine) Validate() error {
	if w.Name == "" {
		return errors.New("Required Name")
	}
	if w.Description == "" {
		return errors.New("Required Description")
	}
	return nil
}

func (w *Wine) SaveWine(db *gorm.DB) (*Wine, error) {
	var err error
	err = db.Debug().Create(&w).Error
	if err != nil {
		return &Wine{}, err
	}
	return w, nil
}

func (w *Wine) FindAllWines(db *gorm.DB) (*[]Wine, error) {
	var err error
	wines := []Wine{}
	err = db.Debug().Model(&Wine{}).Limit(100).Find(&wines).Error
	if err != nil {
		return &[]Wine{}, err
	}
	return &wines, err
}

func (w *Wine) FindWineByID(db *gorm.DB, uid uint64) (*Wine, error) {
	var err error
	
	err = db.Debug().Model(&Wine{}).Where("id = ?", uid).Take(&w).Error
	if err != nil {
		return &Wine{}, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return &Wine{}, errors.New("Wine Not Found")
	}
	return w, err
}

func (w *Wine) UpdateWine(db *gorm.DB) (*Wine, error) {
	var err error

	err = db.Debug().Model(&Wine{}).Where("id = ?", w.ID).Updates(Wine{Name: w.Name, Description: w.Description, Year: w.Year, Price: w.Price, Image: w.Image, Available: w.Available}).Error
	if err != nil {
		return &Wine{}, err
	}
	if w.ID != 0 {
		err = db.Debug().Model(&Wine{}).Where("id = ?", w.ID).Take(&w.ID).Error
		if err != nil {
			return &Wine{}, err
		}
	}
	return w, nil
}

func (w *Wine) DeleteWine(db *gorm.DB, pid uint64) (int64, error) {
	db = db.Debug().Model(&Wine{}).Where("id = ?", pid).Take(&Wine{}).Delete(&Wine{})

	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("Wine not found")
		}
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
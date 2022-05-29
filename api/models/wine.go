package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

type Wine struct {
	ID     	    uint64    `gorm:"primary_key;auto_increment" json:"id"`
	Nome  	    string    `gorm:"size:200;not null;unique" json:"title"`
	Descricao   string    `gorm:"size:200;not null;" json:"content"`
	Ano         string    `gorm:"size:200;not null;" json:"ano"`
	Preco 	    string    `gorm:"size:200;not null;" json:"preco"`
	Imagem 	    string    `gorm:"size:200;not null;" json:"imagem"`
	Disponivel  bool      `gorm:"type:bool;default:false"`
}

func (w *Wine) Prepare() {
	w.ID = 0
	w.Nome = html.EscapeString(strings.TrimSpace(w.Nome))
	w.Descricao = html.EscapeString(strings.TrimSpace(w.Descricao))
}

func (w *Wine) Validate() error {

	if w.Nome == "" {
		return errors.New("Required Nome")
	}
	if w.Descricao == "" {
		return errors.New("Required Descricao")
	}
	return nil
}

func (w *Wine) SaveWine(db *gorm.DB) (*Wine, error) {
	var err error
	err = db.Debug().Model(&Wine{}).Create(&w).Error
	if err != nil {
		return &Wine{}, err
	}
	if w.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", w.AuthorID).Take(&w.Author).Error
		if err != nil {
			return &Wine{}, err
		}
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
	if len(wines) > 0 {
		for i, _ := range wines {
			err := db.Debug().Model(&User{}).Where("id = ?", wines[i].AuthorID).Take(&wines[i].Author).Error
			if err != nil {
				return &[]Wine{}, err
			}
		}
	}
	return &wines, nil
}

func (w *Wine) FindWineByID(db *gorm.DB, pid uint64) (*Wine, error) {
	var err error
	err = db.Debug().Model(&Wine{}).Where("id = ?", pid).Take(&w).Error
	if err != nil {
		return &Wine{}, err
	}
	if w.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", w.AuthorID).Take(&w.Author).Error
		if err != nil {
			return &Wine{}, err
		}
	}
	return w, nil
}

func (w *Wine) UpdateWine(db *gorm.DB) (*Wine, error) {

	var err error

	err = db.Debug().Model(&Wine{}).Where("id = ?", w.ID).Updates(Wine{Nome: w.Nome, Descricao: w.Descricao, UpdatedAt: time.Now()}).Error
	if err != nil {
		return &Wine{}, err
	}
	if w.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", w.AuthorID).Take(&w.Author).Error
		if err != nil {
			return &Wine{}, err
		}
	}
	return w, nil
}

func (w *Wine) DeleteWine(db *gorm.DB, pid uint64, uid uint32) (int64, error) {

	db = db.Debug().Model(&Wine{}).Where("id = ? and author_id = ?", pid, uid).Take(&Wine{}).Delete(&Wine{})

	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("Wine not found")
		}
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
package models

import "time"

type Comic struct {
	ID         uint       `json:"id" gorm:"index;primaryKey"`
	Slug       string     `json:"slug" gorm:"index;unique;not null"`
	Title      string     `json:"title" gorm:"not null"`
	Author     string     `json:"author"`
	ThumbUrl   string     `json:"thumb_url"`
	PosterUrl  string     `json:"poster_url"`
	Categories []Category `gorm:"many2many:comics_categories;"`
	Nation     []Nation   `gorm:"many2many:comics_nations;"`
	Chapters   []Chapter  `gorm:"foreignKey:ComicID"`
	CreateAt   time.Time  `json:"create_at"`
	UpdateAt   time.Time  `json:"update_at"`
}

package models

import "time"

type Chapter struct {
	ID          string    `json:"id" gorm:"primaryKey"`
	ChapterPath string    `json:"chapter_api_data"`
	ChapterName string    `json:"chapter_name"`
	ComicID     uint      `json:"comic_id" gorm:"not null"`
	CreatedAt   time.Time `json:"created_at"`
}

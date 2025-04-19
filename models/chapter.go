package models

type Chapter struct {
	ID           string `json:"id" gorm:"primaryKey"`
	ChapterTitle string `json:"chapter_title" gorm:"not null"`
	ChapterName  string `json:"chapter_name"`
	ChapterPath  string `json:"chapter_path"`
	ComicID      uint   `json:"comic_id" gorm:"not null"`
	Comic        Comic  `gorm:"foreignKey:ComicID"`
}

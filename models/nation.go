package models

type Nation struct {
	ID    int     `gorm:"primary_key" json:"id"`
	Name  string  `json:"name"`
	Slug  string  `json:"slug" gorm:"index" gorm:"unique"`
	Comic []Comic `gorm:"many2many:comics_nations;"`
}

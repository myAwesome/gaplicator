package models


type Post struct {
	ID uint `gorm:"primarykey" json:"id"`
	Title                string      `gorm:"column:title;size:200" json:"title"`
}

func (Post) TableName() string { return "posts" }


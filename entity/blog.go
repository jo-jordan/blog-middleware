package entity

type Blog struct {
	ID       uint64
	ParentId uint64 `gorm:"column:parent_id"`
	Type     uint   `gorm:"column:type"`
	Name     string `gorm:"column:blog_name"`
	Children []Blog `gorm:"_"`
}

func (Blog) TableName() string {
	return "blogs"
}

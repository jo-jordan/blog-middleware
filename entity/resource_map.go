package entity

type ResourceMap struct {
	Key   string `gorm:"column:res_key"`
	Value string `gorm:"column:res_value"`
}

func (ResourceMap) TableName() string {
	return "resource_map"
}

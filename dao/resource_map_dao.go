package dao

import (
	"github.com/lzjlxebr/blog-middleware/common"
	"github.com/lzjlxebr/blog-middleware/entity"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func ResourceMapGetValueBy(key string) *entity.ResourceMap {
	db, err := gorm.Open(mysql.Open(common.MySQLConnURL), &gorm.Config{})
	common.ErrorBus(err)

	var res entity.ResourceMap
	db.Where(&entity.ResourceMap{Key: key}).First(&res)
	return &res
}

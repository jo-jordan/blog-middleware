package dao

import (
	"github.com/lzjlxebr/blog-middleware/common"
	"github.com/lzjlxebr/blog-middleware/entity"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func AccessLogSave(al entity.AccessLog) {
	db, err := gorm.Open(mysql.Open(common.MySQLConnURL), &gorm.Config{})
	common.ErrorBus(err)

	db.Create(al)
}

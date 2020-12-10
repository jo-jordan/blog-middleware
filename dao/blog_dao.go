package dao

import (
	"github.com/lzjlxebr/blog-middleware/common"
	"github.com/lzjlxebr/blog-middleware/entity"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func BlogDeleteAll() {
	db, err := gorm.Open(mysql.Open(common.MySQLConnURL), &gorm.Config{})
	common.ErrorBus(err)

	db.Where("1 = 1").Delete(&entity.Blog{})
}

func BlogFindAll() []entity.Blog {
	db, err := gorm.Open(mysql.Open(common.MySQLConnURL), &gorm.Config{})
	common.ErrorBus(err)

	var blogs []entity.Blog
	db.Find(&blogs)

	return blogs
}

func BlogSaveAll(blogs []entity.Blog) {
	db, err := gorm.Open(mysql.Open(common.MySQLConnURL), &gorm.Config{})
	common.ErrorBus(err)

	db.CreateInBatches(blogs, 200)
}

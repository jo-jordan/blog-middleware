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

func BlogFindById(id uint64) entity.Blog {
	db, err := gorm.Open(mysql.Open(common.MySQLConnURL), &gorm.Config{})
	common.ErrorBus(err)

	var blog entity.Blog
	db.Where(&entity.Blog{ID: id, Type: 1}).Find(&blog)
	blog.URL = common.NOTES_HOST + blog.Path
	return blog
}

func BlogFindByCategory(category uint64) []entity.Blog {
	db, err := gorm.Open(mysql.Open(common.MySQLConnURL), &gorm.Config{})
	common.ErrorBus(err)

	var blogs []entity.Blog
	db.Where(&entity.Blog{ParentId: category, Type: 1}).Find(&blogs)

	return blogs
}

func CategoryFindAll() []entity.Blog {
	db, err := gorm.Open(mysql.Open(common.MySQLConnURL), &gorm.Config{})
	common.ErrorBus(err)

	var blogs []entity.Blog
	db.Where(map[string]interface{}{"Type": 0}).Find(&blogs)

	return blogs
}

func BlogSaveAll(blogs []entity.Blog) {
	db, err := gorm.Open(mysql.Open(common.MySQLConnURL), &gorm.Config{})
	common.ErrorBus(err)

	_ = db.Transaction(func(tx *gorm.DB) error {

		db.Where("1 = 1").Delete(&entity.Blog{})
		db.CreateInBatches(blogs, 200)

		// 返回 nil 提交事务
		return nil
	})
}

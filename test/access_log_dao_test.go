package test

import (
	"github.com/bwmarrin/snowflake"
	"github.com/lzjlxebr/blog-middleware/common"
	"github.com/lzjlxebr/blog-middleware/dao"
	"github.com/lzjlxebr/blog-middleware/entity"
	"testing"
	"time"
)

func TestAccessLogSave(t *testing.T) {
	node, err := snowflake.NewNode(1)
	common.ErrorBus(err)

	dao.AccessLogSave(entity.AccessLog{
		ID:         uint64(node.Generate().Int64()),
		RealIP:     111110,
		Time:       time.Time{}.Local(),
		Type:       1110,
		ResourceID: 110,
	})
}

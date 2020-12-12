package test

import (
	"crypto/md5"
	"encoding/binary"
	"fmt"
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

func TestStringParse(t *testing.T) {

	bytes := []byte("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF")
	u := binary.LittleEndian.Uint64(md5.New().Sum(bytes))

	fmt.Println(u)
}

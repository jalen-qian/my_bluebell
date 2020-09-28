package snowflake

/**
雪花算法
*/
import (
	"github.com/bwmarrin/snowflake"
	"github.com/pkg/errors"
)

var node *snowflake.Node

func Init(nodeId int64) (err error) {
	node, err = snowflake.NewNode(nodeId)
	if err != nil {
		return err
	}
	return
}

func GetId() (id int64, err error) {
	if node == nil {
		return 0, errors.New("snowflake node not initialized!")
	}
	return node.Generate().Int64(), nil
}

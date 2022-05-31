package models

import (
	"app/library/consts"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"time"
)

func GetNode(values ...interface{}) *Node {
	for _, v := range values {
		switch vv := v.(type) {
		case uint32:
			var i Node
			if err := DB().Take(&i, "id=?", vv).Error; err != nil {
				return nil
			}
			return &i
		case *gin.Context:
			if i, exist := vv.Get(consts.ContextNode); exist {
				if instance, ok := i.(*Node); ok {
					return instance
				}
			}
			return nil
		}
	}
	return nil
}

type NodeMarshalJSON Node

//virtual Node
type Node struct {
	ID          uint32  `gorm:"primarykey;column:id" json:"id" form:"id"`
	Name        string  `gorm:"" json:"name" form:"name" binding:"required"`
	Desc        string  `gorm:"column:desc" json:"desc" form:"desc" binding:"required"`
	IP          string  `gorm:"" json:"ip" form:"ip" binding:"required"` //内网IP
	Uid         *uint32 `gorm:"" json:"uid" form:"-"`
	SshPort     string  `gorm:"" json:"sshPort" form:"required"`
	SshKey      string  `gorm:"" json:"sshKey" form:"required"` //private key
	SshUsername string  `gorm:"" json:"sshUsername" form:"required"`
	SshPassword string  `gorm:"" json:"sshPassword" form:"required"`

	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (t *Node) TableName() string {
	return "d_node"
}

// 实现 sql.Scanner 接口，Scan 将 value 扫描至 Jsonb
func (t *Node) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}
	json.Unmarshal(bytes, &t) //no error
	return nil
}

// 实现 driver.Valuer 接口，Value 返回 json value
func (t Node) Value() (driver.Value, error) {
	return json.Marshal(t)
}

//override marshal json
func (t Node) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		NodeMarshalJSON
	}{
		NodeMarshalJSON(t),
	})
}

//k8s cluster visual Node data validator check
func (t *Node) Validator(c *gin.Context) error {
	return nil
}

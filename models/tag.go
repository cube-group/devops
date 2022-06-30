package models

import (
	"app/library/consts"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"time"
)

func GetTag(values ...interface{}) *Tag {
	for _, v := range values {
		switch vv := v.(type) {
		case uint32:
			var i Tag
			if err := DB().Take(&i, "id=?", vv).Error; err != nil {
				return nil
			}
			return &i
		case string:
			var i Tag
			if err := DB().Take(&i, "name=?", vv).Error; err != nil {
				return nil
			}
			return &i
		case *gin.Context:
			if i, exist := vv.Get(consts.ContextTag); exist {
				if instance, ok := i.(*Tag); ok {
					return instance
				}
			}
			return nil
		}
	}
	return nil
}

func TagList() (res []Tag) {
	DB().Order("id DESC").Find(&res)
	return
}

func TagRelProject(tx *gorm.DB, pid, tid uint32) (err error) {
	if tid > 0 {
		if tag := GetTag(tid); tag != nil {
			if err = tag.RelProject(tx, pid); err != nil {
				return
			}
		}
	} else {
		if err = tx.Unscoped().Delete(&TagRel{}, "pid=?", pid).Error; err != nil {
			return
		}
	}
	return
}

type TagRel struct {
	ID        uint32         `gorm:"primarykey;column:id" json:"id" form:"id"`
	Tid       uint32         `gorm:"" json:"tid" form:"-"`
	Pid       uint32         `gorm:"" json:"pid" form:"-"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (t *TagRel) TableName() string {
	return "d_tag_rel"
}

type TagMarshalJSON Tag

type Tag struct {
	ID        uint32         `gorm:"primarykey;column:id" json:"id" form:"id"`
	Name      string         `gorm:"" json:"name" form:"name" binding:"required"`
	Desc      string         `gorm:"column:desc" json:"desc" form:"desc" binding:"required"`
	Uid       *uint32        `gorm:"" json:"uid" form:"-"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (t *Tag) TableName() string {
	return "d_tag"
}

//override marshal json
func (t Tag) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		TagMarshalJSON
	}{
		TagMarshalJSON(t),
	})
}

func (t *Tag) ProjectIds() (res []uint32) {
	var rel []TagRel
	if DB().Find(&rel, "tid=?", t.ID).Error != nil {
		return
	}
	for _, item := range rel {
		res = append(res, item.Pid)
	}
	return
}

func (t *Tag) RelProject(tx *gorm.DB, pid uint32) error {
	if tx == nil {
		tx = DB()
	}
	var find TagRel
	if tx.Take(&find, "tid=? AND pid=?", t.ID, pid).Error == nil {
		return nil
	}
	return tx.Save(&TagRel{Tid: t.ID, Pid: pid}).Error
}

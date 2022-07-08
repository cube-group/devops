package models

import (
	"app/library/log"
	"app/library/types/times"
	"errors"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
	"sync"
	"time"
)

func GetProjectCronjob(values ...interface{}) (res *ProjectCronjob) {
	for _, v := range values {
		switch vv := v.(type) {
		case uint32:
			var i ProjectCronjob
			if err := DB().Take(&i, "project_id=?", vv).Error; err != nil {
				return nil
			}
			res = &i
		}
	}
	return
}

type ProjectCronjob struct {
	ID        uint32         `gorm:"primarykey;column:id" json:"id"`
	Uid       uint32         `gorm:"" json:"uid"`
	NodeId    uint32         `gorm:"" json:"nodeId"`
	Node      *Node          `gorm:"" json:"node"`
	ProjectId uint32         `gorm:"" json:"projectId"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (t *ProjectCronjob) TableName() string {
	return "d_project_cronjob"
}

var c *cron.Cron
var cMap sync.Map

func InitCronjob() {
	go initSystemCronjob()
	initProjectCronjob()
}

func initSystemCronjob() {
	var systemCron = cron.New()
	systemCron.AddFunc("0 4 */1 * *", func() {
		NodeClean()
	})
	systemCron.Run()
}

func initProjectCronjob() {
	time.Sleep(10 * time.Second)
	var list []ProjectCronjob
	if err := DB().Find(&list).Error; err != nil {
		log.StdWarning("cronjob", "init", err)
	}
	c = cron.New()
	for _, item := range list {
		CronjobAdd(item)
	}
	c.Run()
}

//add cronjob
func CronjobAdd(cronjob ProjectCronjob) (err error) {
	var project = GetProject(cronjob.ProjectId)
	if project == nil {
		return errors.New("project not found")
	}
	if project.Mode != ProjectModeDocker && project.Mode != ProjectModeNative {
		return errors.New("cronjob not supports this mode")
	}
	//find projectCronjob
	if cronjob.ID == 0 {
		var db = DB()
		var findCronjob ProjectCronjob
		if db.Take(&findCronjob, "project_id=?", cronjob.ProjectId).Error == nil {
			cronjob.ID = findCronjob.ID
		}
		if err = db.Save(&cronjob).Error; err != nil {
			return err
		}
	}
	//remove before add
	if value, ok := cMap.Load(project.ID); ok {
		if entryID, ok := value.(cron.EntryID); ok {
			c.Remove(entryID)
		}
	}
	//add
	entryID, err := c.AddFunc(project.Cronjob, func() {
		project = GetProject(cronjob.ProjectId)
		if project == nil {
			return
		}
		var version = "cronjob-" + times.FormatFileDatetime(time.Now())
		var er = project.Apply(&History{
			Uid:       cronjob.Uid,
			Version:   version,
			Desc:      version,
			NodeId:    cronjob.NodeId,
			Node:      cronjob.Node,
			ProjectId: project.ID,
			Project:   project,
		}, false)
		if er != nil {
			log.StdWarning("cronjob", "projectID", project.ID, project.Cronjob, er)
		}
	})
	if err != nil {
		return
	}
	cMap.Store(project.ID, entryID)
	return nil
}

//remove or stop cronjob
func CronjobStop(projectID uint32) error {
	var i = GetProjectCronjob(projectID)
	if i != nil {
		if err := DB().Unscoped().Delete(i).Error; err != nil {
			return err
		}
	}
	if value, ok := cMap.Load(projectID); ok {
		if entryID, ok := value.(cron.EntryID); ok {
			c.Remove(entryID)
		}
	}
	return nil
}

func CronjobStopAll() {
	c.Stop()
}

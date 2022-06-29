package task

import (
	"app/models"
	"github.com/robfig/cron/v3"
)

func Init() {
	//cron.WithSeconds() //支持到秒级
	//cron.WithChain(cron.SkipIfStillRunning(cron.DefaultLogger))//没执行完跳过本次函数
	//cron.WithLogger(cron.VerbosePrintfLogger(log.New(os.Stdout, "cron: ", log.LstdFlags)))//打印任务日志
	c := cron.New()
	c.AddFunc("0 4 */1 * *", func() {
		models.NodeClean()
	})
	c.Run()
}

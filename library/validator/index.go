package validator

import (
	"github.com/robfig/cron/v3"
)

func IsCron(spec string) bool {
	parser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	_, err := parser.Parse(spec)
	if err != nil {
		return false
	}
	return true
}

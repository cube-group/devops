package times

import (
	"app/library/types/jsonutil"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"math"
	"strings"
	"time"
)

var months = map[string]int{
	"January":   1,
	"February":  2,
	"March":     3,
	"April":     4,
	"May":       5,
	"June":      6,
	"July":      7,
	"August":    8,
	"September": 9,
	"October":   10,
	"November":  11,
	"December":  12,
}

//创建JsonTime
func NewJsonTime(option ...time.Time) JsonTime {
	if len(option) == 0 {
		return JsonTime{Time: time.Now()}
	} else {
		return JsonTime{Time: option[0]}
	}
}

//json时间格式化
type JsonTime struct {
	time.Time
}

//写入数据库时会调用该方法将自定义时间类型转换并写入数据库
func (t JsonTime) Value() (driver.Value, error) {
	var zeroTime time.Time
	if t.Time.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return t.Time, nil
}

//读取数据库时会调用该方法将时间数据转换成自定义时间类型
func (t *JsonTime) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*t = JsonTime{Time: value}
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}

//实现json序列化方法,格式化时间
func (t JsonTime) MarshalJSON() ([]byte, error) {
	var stamp = gin.H{"datetime": FormatDatetime(t.Time), "str": TimeFormat(t.Time)}
	return jsonutil.ToBytes(stamp), nil
}

//TODO datetime默认值
func Default() time.Time {
	value, _ := time.Parse("0000-00-00 00:00:00", "1970-01-01 08:00:01")
	return value
}

//格式化为文件使用的年月日时分秒
func FormatFileDatetime(date time.Time) string {
	return date.Format("20060102150405")
}

//格式化为年月日时分秒
func FormatDatetime(date time.Time) string {
	return date.Format("2006-01-02 15:04:05")
}

//格式化为月日时分秒
func FormatDaytime(date time.Time) string {
	return date.Format("01-02 15:04:05")
}

//格式化为月日时分秒
func FormatTime(date time.Time) string {
	return date.Format("15:04:05")
}

//格式华为年月日
func FormatDate(date time.Time) string {
	return date.Format("2006-01-02")
}

//格式华为年月日时
func FormatDay(date time.Time) string {
	return date.Format("02")
}

//格式华为年月日时
func FormatDayHour(date time.Time) string {
	return date.Format("02 15:00")
}

//格式华为年月日时
func FormatHour(date time.Time) string {
	return date.Format("15:00")
}

//格式化为分秒
func FormatMinute(date time.Time) string {
	return date.Format("04:05")
}

func GetFormatMonth(options ...time.Time) int {
	var current time.Time
	if len(options) == 0 {
		current = time.Now()
	} else {
		current = options[0]
	}
	_, month, _ := current.Date()
	return months[month.String()]
}

/**
 * 1. 如果距当前时间60s内,则显示x秒内
 * 2. 如果距当前时间60m内,则显示x分钟内
 * 3. 如果距当前时间24h内,则显示x小时内
 * 4. 如果超过24小时,则显示x天前
 * @param $time
 * @return string
 */
func TimeFormat(v time.Time) string {
	if v == Default() {
		return ""
	}
	now := time.Now()
	diff := now.Sub(v).Seconds()

	if diff < 0 {
		return "非法时间"
	}
	if diff < 60 {
		return "刚刚"
	}
	if diff < (60 * 60) {
		return fmt.Sprintf("%v分钟前", math.Floor(diff/60))
	}
	if diff < (60 * 60 * 24) {
		return fmt.Sprintf("%v小时前", math.Floor(diff/(60*60)))
	}
	return fmt.Sprintf("%v天前", math.Floor(diff/(60*60*24)))

	if diff < (60 * 60 * 24 * 30) {
		return fmt.Sprintf("%v天前", math.Floor(diff/(60*60*24)))
	}
	if diff < (60 * 60 * 24 * 30 * 12) {
		return fmt.Sprintf("%v月前", math.Floor(diff/(60*60*24*30)))
	}
	return fmt.Sprintf("%v年前", math.Floor(diff/(60*60*24*30*12)))

}

//获取最近指定天数的datetime
func LastDatetime(currentTime time.Time, days int) string {
	return FormatDatetime(time.Unix(currentTime.Unix()-int64(days*24*3600), 0))
}

//格式化为当天开始时间
func DayStarTime(current time.Time) time.Time {
	return time.Date(current.Year(), current.Month(), current.Day(), 0, 0, 0, 0, current.Location())
}

//格式化为当天结束时间
func DayEndTime(current time.Time) time.Time {
	return time.Date(current.Year(), current.Month(), current.Day(), 23, 59, 59, 0, current.Location())
}

func Strtotime(v string) (time.Time, error) {
	res, err := time.ParseInLocation("2006-01-02 15:04:05", v, time.Local)
	return res, err
}
func StrTimeFormat(v string) string {
	res, _ := Strtotime(v)
	return TimeFormat(res)
}

//获取utc时间无小数点
//例如：2018-11-01T07:15:56
func GetNowUTCTimeString() (string, int64, error) {
	var cstZone = time.FixedZone("CST", 8*3600) // 东八
	t := time.Now().In(cstZone).UTC()
	source := []time.Time{t}
	bytes, err := json.Marshal(source)
	if err != nil {
		return "", 0, err
	}
	s := string(bytes)
	s = strings.Split(s, "[\"")[1]
	splits := strings.Split(s, ".")
	return fmt.Sprintf("%sZ", splits[0]), t.Unix(), nil
}

func GetLocalTimeFromUTCTime(utcString string) (time.Time, error) {
	utc, err := time.Parse("2006-01-02T15:04:05Z", utcString)
	if err != nil {
		return time.Now(), err
	}
	return utc, nil
}

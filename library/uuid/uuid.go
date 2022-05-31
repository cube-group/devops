package uuid

import (
	"app/library/crypt/md5"
	"app/library/types/jsonutil"
	"fmt"
	"math/rand"
	"time"
)

//获取唯一md5id
func GetUUID(prefix ...interface{}) string {
	return md5.MD5(fmt.Sprintf("%s%d%d", jsonutil.ToString(prefix), rand.Int63n(1000), time.Now().Nanosecond()))
}

package unique

import (
	"app/library/crypt/md5"
	"app/library/types/jsonutil"
	"app/library/uuid"
	"fmt"
	"math/rand"
	"time"
)

func Id(args ...string) string {
	return md5.MD5(fmt.Sprintf("%s%s%s%d%d", jsonutil.ToString(args), uuid.GetUUID(), time.Now().String(), rand.Intn(100000), time.Now().UnixNano()))
}

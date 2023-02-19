package common

import "time"

var TOKEN_EXPIRE_TIME time.Time = time.Now().Add(time.Minute * 30)

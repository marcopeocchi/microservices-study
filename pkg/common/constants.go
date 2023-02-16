package common

import "time"

const BCRYPT_ROUNDS int = 12

var TOKEN_EXPIRE_TIME time.Time = time.Now().Add(time.Minute * 30)

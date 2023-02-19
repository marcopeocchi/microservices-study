package common

import "time"

const TOKEN_COOKIE_NAME string = "jwt_token"

var TOKEN_EXPIRE_TIME time.Time = time.Now().Add(time.Minute * 30)

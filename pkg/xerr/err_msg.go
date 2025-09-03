package xerr

var codeText = map[int]string{
	SERVER_COMMON_ERROR: "Server Internal Error",
	REQUEST_PATH_ERROR:  "Path Error",
	DB_ERROR:            "DB Error",
}

func ErrMsg(errCode int) string {
	if msg, ok := codeText[errCode]; ok {
		return msg
	} else {
		return "Unknown Error"
	}
}

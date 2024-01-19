package consts

const (
	Success  = 200
	NotFound = 404
	Error    = 500 // 业务错误
	Failed   = 501 // 系统错误
	Warn     = 201 // 警告
)

var ResponseMap map[uint16]string = map[uint16]string{
	Success:  "success",
	NotFound: "not found",
	Failed:   "failed",
	Error:    "error",
	Warn:     "warn",
}

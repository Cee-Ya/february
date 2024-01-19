package exception

// BizException 常规业务异常
type BizException struct {
	Err error
}

func (b *BizException) Error() string {
	return b.Err.Error()
}

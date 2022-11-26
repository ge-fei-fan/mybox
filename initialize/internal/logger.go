package internal

import "gorm.io/gorm/logger"

type writer struct {
	logger.Writer
}

func NewWriter(w logger.Writer) *writer {
	return &writer{Writer: w}
}

func (w *writer) Printf(message string, data ...interface{}) {
	//todo 日志后面加上来
	//简单打印
	w.Writer.Printf(message, data...)
}

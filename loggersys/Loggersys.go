package loggersys

import (
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

// Log aa system printer log
var Log *logrus.Logger

func init() {
	NewLogger()
}

// NewLogger create file logger system
func NewLogger() *logrus.Logger {
	if Log != nil {
		return Log
	}

	path := "/home/qa/works/anthenacmc/anthena_%Y%m%d%H.log"
	lpath := "/home/qa/works/anthenacmc/anthena.log"
	writer, err := rotatelogs.New(
		path,
		// WithLinkName为最新的日志建立软连接,以方便随着找到当前日志文件
		rotatelogs.WithLinkName(lpath),

		// WithRotationTime设置日志分割的时间,这里设置为一小时分割一次
		rotatelogs.WithRotationTime(time.Hour),

		// WithMaxAge和WithRotationCount二者只能设置一个,
		// WithMaxAge设置文件清理前的最长保存时间,
		// WithRotationCount设置文件清理前最多保存的个数.
		//rotatelogs.WithMaxAge(time.Hour*24),
		rotatelogs.WithRotationCount(20),
	)

	if err != nil {
		logrus.Errorf("config local file system for logger error: %v", err)
	}

	Log = logrus.New()
	Log.Hooks.Add(lfshook.NewHook(lfshook.WriterMap{
		logrus.DebugLevel: writer,
		logrus.InfoLevel:  writer,
		logrus.WarnLevel:  writer,
		logrus.ErrorLevel: writer,
		logrus.FatalLevel: writer,
		logrus.PanicLevel: writer,
	}, &logrus.TextFormatter{DisableColors: true}))

	return Log
}

package loggersys

import (
	"fmt"
	"os"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

// Log aa system printer log
var Log *logrus.Logger

func init() {
	if _, err := os.Stat("logs"); os.IsNotExist(err) {
		// /var/log/anthena does not exist
		if err = os.Mkdir("logs", 0775); err != nil {
			fmt.Println(err)
		}
	}
}

// NewLogger create file logger system
func NewLogger(filename string) *logrus.Logger {
	if Log != nil {
		return Log
	}

	path := fmt.Sprintf("logs/%s", filename)
	path += "_%Y%m%d%H.log"
	lpath := fmt.Sprintf("logs/%s.log", filename)
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

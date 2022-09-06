package logger

import (
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	log "github.com/sirupsen/logrus"
	"time"
)

var Logger = log.New()

func init() {
	configLocalFilesystemLogger1("log/log")
}

//切割日志和清理过期日志
func configLocalFilesystemLogger1(filePath string) {
	writer, err := rotatelogs.New(
		filePath+"-%Y-%m-%d-%H.log",
		rotatelogs.WithLinkName(filePath),             // 生成软链，指向最新日志文件
		rotatelogs.WithMaxAge(time.Second*60*60*24*3), // 文件最大保存时间
		rotatelogs.WithRotationTime(time.Second*60),   // 日志切割时间间隔
	)
	if err != nil {
		Logger.Fatal("Init log failed, err:", err)
	}
	Logger.SetOutput(writer)
	Logger.SetLevel(log.InfoLevel)
	Logger.SetReportCaller(true)
	//Logger.SetFormatter(&log.JSONFormatter{})
	Logger.AddHook(&DefaultFieldsHook{})
}

type DefaultFieldsHook struct {
}

func (df *DefaultFieldsHook) Fire(entry *log.Entry) error {
	//entry.Data["appName"] = "GoFrame"
	return nil
}

func (df *DefaultFieldsHook) Levels() []log.Level {
	return log.AllLevels
}

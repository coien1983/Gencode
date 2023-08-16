package template

var LogTemplate = `package sysinit

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func LogInit(logger *LogC) (err error) {
	writerSyncer := getLogWriter(
		logger.Filename,
		logger.MaxSize,
		logger.MaxBackups,
		logger.MaxAge,
	)
	encoder := getEncoder()

	var l = new(zapcore.Level)
	err = l.UnmarshalText([]byte(logger.Level))
	if err != nil {
		return
	}

	core := zapcore.NewCore(encoder, writerSyncer, l)
	lg := zap.New(core, zap.AddCaller())
	//替换zap库中全局的logger
	zap.ReplaceGlobals(lg)
	return
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	return zapcore.NewJSONEncoder(encoderConfig)
}

func getLogWriter(filename string, maxSize, maxBackup, maxAge int) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    maxSize,
		MaxAge:     maxAge,
		MaxBackups: maxBackup,
	}

	return zapcore.AddSync(lumberJackLogger)
}
`

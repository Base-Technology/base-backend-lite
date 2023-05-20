package seelog

import (
	"fmt"
	"log"
	"path/filepath"
	"runtime"

	"github.com/cihub/seelog"
	"github.com/petermattis/goid"
)

func init() {
	seelog.RegisterCustomFormatter("ServiceName", createAppNameFormatter)
	logger, err := seelog.LoggerFromConfigAsString(debugConfig)
	if err != nil {
		log.Fatalf("parse seelog config error, %v", err)
	}
	seelog.ReplaceLogger(logger)
}

// SetDebugLevel 设置日志界别为debug
func SetDebugLevel() {
	logger, err := seelog.LoggerFromConfigAsString(debugConfig)
	if err != nil {
		log.Fatalf("parse seelog config error, %v", err)
	}
	seelog.ReplaceLogger(logger)
}

// SetInfoLevel 设置日志界别为info
func SetInfoLevel() {
	logger, err := seelog.LoggerFromConfigAsString(infoConfig)
	if err != nil {
		log.Fatalf("parse seelog config error, %v", err)
	}
	seelog.ReplaceLogger(logger)
}

// SetWarnLevel 设置日志界别为warn
func SetWarnLevel() {
	logger, err := seelog.LoggerFromConfigAsString(warnConfig)
	if err != nil {
		log.Fatalf("parse seelog config error, %v", err)
	}
	seelog.ReplaceLogger(logger)
}

// SetErrorLevel 设置日志界别为error
func SetErrorLevel() {
	logger, err := seelog.LoggerFromConfigAsString(errorConfig)
	if err != nil {
		log.Fatalf("parse seelog config error, %v", err)
	}
	seelog.ReplaceLogger(logger)
}

const (
	// LogLevelError error级别
	LogLevelError = "ERROR"
	// LogLevelInfo info级别
	LogLevelInfo = "INFO"
	// LogLevelDebug debug级别
	LogLevelDebug = "DEBUG"
	// LogLevelWarn warn级别
	LogLevelWarn = "WARN"
)

func createAppNameFormatter(params string) seelog.FormatterFunc {
	return func(message string, level seelog.LogLevel, context seelog.LogContextInterface) interface{} {
		serviceName := "base-backend-lite"
		return serviceName
	}
}

// Errorf error级别格式化日志
func Errorf(format string, params ...interface{}) error {
	prefix := getPrefix(LogLevelError)
	return seelog.Errorf(prefix+format, params...)
}

// Error error级别日志
func Error(params ...interface{}) error {
	prefix := getPrefix(LogLevelError)
	var newParams []interface{}
	newParams = append(newParams, prefix)
	for _, param := range params {
		newParams = append(newParams, param)
	}
	return seelog.Error(prefix, newParams)
}

// Infof info级别格式化日志
func Infof(format string, params ...interface{}) {
	seelog.Infof(getPrefix(LogLevelInfo)+format, params...)
}

// Info info级别日志
func Info(params ...interface{}) {
	prefix := getPrefix(LogLevelInfo)
	var newParams []interface{}
	newParams = append(newParams, prefix)
	for _, param := range params {
		newParams = append(newParams, param)
	}
	seelog.Info(newParams...)
}

// Debugf debug级别格式化日志
func Debugf(format string, params ...interface{}) {
	seelog.Debugf(getPrefix(LogLevelDebug)+format, params...)
}

// Debug debug级别日志
func Debug(params ...interface{}) {
	prefix := getPrefix(LogLevelDebug)
	var newParams []interface{}
	newParams = append(newParams, prefix)
	for _, param := range params {
		newParams = append(newParams, param)
	}
	seelog.Debug(newParams...)
}

// Warnf warn级别格式化日志
func Warnf(format string, params ...interface{}) {
	seelog.Warnf(getPrefix(LogLevelWarn)+format, params...)
}

// Warn warn级别日志
func Warn(params ...interface{}) {
	prefix := getPrefix(LogLevelWarn)
	var newParams []interface{}
	newParams = append(newParams, prefix)
	for _, param := range params {
		newParams = append(newParams, param)
	}
	seelog.Warn(newParams...)
}

// Flush 刷新日志到文件
func Flush() {
	seelog.Flush()
}

func getPrefix(level string) string {
	callerInfo := getCallerName()
	return fmt.Sprintf("%s [%d] %s: ", level, goid.Get(), callerInfo)
}

func getCallerName() string {
	_, file, line, _ := runtime.Caller(3)
	return fmt.Sprintf("%s.%d", filepath.Base(file), line)
}

var debugConfig string = `
<seelog minlevel="debug">
	<outputs formatid="fmt_info">
		<filter levels="debug,info,warn,error">
			<rollingfile formatid="fmt_info" type="size" filename="./logs/application.log" maxsize="33554432" maxrolls="64"/>
		</filter>
		<console/>
	</outputs>
	<formats>
		<format id="fmt_info" format="%Date(2006-01-02 15:04:05.999) %ServiceName %Msg%n" />
		<format id="fmt_err" format="%Date(2006-01-02 15:04:05.999) %ServiceName %Msg%n" />
	</formats>
</seelog>
`

var infoConfig string = `
<seelog minlevel="info">
	<outputs formatid="fmt_info">
		<filter levels="info,warn,error">
			<rollingfile formatid="fmt_info" type="size" filename="./logs/application.log" maxsize="33554432" maxrolls="64"/>
		</filter>
		<console/>
	</outputs>
	<formats>
		<format id="fmt_info" format="%Date(2006-01-02 15:04:05.999) %ServiceName %Msg%n" />
		<format id="fmt_err" format="%Date(2006-01-02 15:04:05.999) %ServiceName %Msg%n" />
	</formats>
</seelog>
`

var warnConfig string = `
<seelog minlevel="warn">
	<outputs formatid="fmt_info">
		<filter levels="warn,error">
			<rollingfile formatid="fmt_info" type="size" filename="./logs/application.log" maxsize="33554432" maxrolls="64"/>
		</filter>
		<console/>
	</outputs>
	<formats>
		<format id="fmt_info" format="%Date(2006-01-02 15:04:05.999) %ServiceName %Msg%n" />
		<format id="fmt_err" format="%Date(2006-01-02 15:04:05.999) %ServiceName %Msg%n" />
	</formats>
</seelog>
`

var errorConfig string = `
<seelog minlevel="error">
	<outputs formatid="fmt_info">
		<filter levels="error">
			<rollingfile formatid="fmt_info" type="size" filename="./logs/application.log" maxsize="33554432" maxrolls="64"/>
		</filter>
		<console/>
	</outputs>
	<formats>
		<format id="fmt_info" format="%Date(2006-01-02 15:04:05.999) %ServiceName %Msg%n" />
		<format id="fmt_err" format="%Date(2006-01-02 15:04:05.999) %ServiceName %Msg%n" />
	</formats>
</seelog>
`

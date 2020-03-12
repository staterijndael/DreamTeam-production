package logwrap

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
)

func InitializePanic(text string) {
	logrus.Panic(text)
	os.Exit(1)
}

func InvalidArgValue(argName, value, format string, details ...interface{}) {
	myFormat := fmt.Sprintf("Invalid arg value; Arg: %s. Value: %s.", argName, value)
	logrus.Errorf(myFormat+format, details...)
}

func Info(format string, details ...interface{}) {
	logrus.SetReportCaller(false)
	logrus.Infof(format, details...)
}

func Debug(format string, details ...interface{}) {
	logrus.SetReportCaller(true)
	logrus.Debugf(format, details...)
}

func Error(format string, details ...interface{}) {
	logrus.SetReportCaller(true)
	logrus.Errorf(format, details...)
}

func NetworkError(operationName, url, format string, details ...interface{}) {
	myFormat := fmt.Sprintf("NetworkError; Operation: %s; Url: %s; ", operationName, url)
	logrus.Errorf(myFormat+format, details...)
}

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.DebugLevel)
}

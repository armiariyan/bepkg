package mongo

import (
	"path"
	"runtime"
	"strconv"
	"time"

	"github.com/armiariyan/bepkg/logger"
	"go.uber.org/zap"
)

func (conn *connection) SetLogger(l logger.Logger) {
	conn.logger = l
}

func (conn *connection) Debug(d bool) {
	conn.debug = d
}

func (conn *connection) logInfo(startTime time.Time, message ...interface{}) {
	var zapMsg []zap.Field
	var callerFunc string
	timeNow := time.Now()
	pc, _, _, ok := runtime.Caller(1)
	d := runtime.FuncForPC(pc)
	if ok && d != nil {
		callerFunc = path.Base(d.Name())
	}

	zapMsg = append(zapMsg, logger.ToField("caller", callerFunc))
	zapMsg = append(zapMsg, logger.ToField("mongo_rt", timeNow.Sub(startTime).Nanoseconds()))
	for k, v := range message {
		zapMsg = append(zapMsg, logger.ToField(strconv.Itoa(k), v))
	}
	if conn.logger != nil && conn.debug != false {
		conn.logger.Info(
			"MongoDB",
			zapMsg...,
		)
	}
}

func (conn *connection) logError(startTime time.Time, message ...interface{}) {
	var zapMsg []zap.Field
	var callerFunc string

	pc, _, _, ok := runtime.Caller(1)
	d := runtime.FuncForPC(pc)
	if ok && d != nil {
		callerFunc = path.Base(d.Name())
	}

	zapMsg = append(zapMsg, logger.ToField("caller", callerFunc))
	zapMsg = append(zapMsg, logger.ToField("mongo_rt", time.Now().Sub(startTime).Nanoseconds()))
	for k, v := range message {
		zapMsg = append(zapMsg, logger.ToField(strconv.Itoa(k), v))
	}
	if conn.logger != nil {
		conn.logger.Error(
			"MongoDB",
			zapMsg...,
		)
	}
}

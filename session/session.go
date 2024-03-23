package session

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	Logger "github.com/armiariyan/bepkg/logger"
	JsonIter "github.com/json-iterator/go"
	Map "github.com/orcaman/concurrent-map"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type Session struct {
	Map                     Map.ConcurrentMap
	Logger                  Logger.Logger
	RequestTime             time.Time
	ThreadID                string
	AppName, AppVersion, IP string
	Port                    int
	SrcIP, URL, Method      string
	Header, Request         interface{}
	ErrorMessage            string
}

func New(logger Logger.Logger) *Session {
	return &Session{
		RequestTime: time.Now(),
		Logger:      logger,
		Map:         Map.New(),
	}
}

func (session *Session) SetThreadID(threadID string) *Session {
	session.ThreadID = threadID
	return session
}

func (session *Session) SetMethod(method string) *Session {
	session.Method = method
	return session
}

func (session *Session) SetAppName(appName string) *Session {
	session.AppName = appName
	return session
}

func (session *Session) SetAppVersion(appVersion string) *Session {
	session.AppVersion = appVersion
	return session
}

func (session *Session) SetURL(url string) *Session {
	session.URL = url
	return session
}

func (session *Session) SetIP(ip string) *Session {
	session.IP = ip
	return session
}

func (session *Session) SetPort(port int) *Session {
	session.Port = port
	return session
}

func (session *Session) SetSrcIP(srcIp string) *Session {
	session.SrcIP = srcIp
	return session
}

func (session *Session) SetHeader(header interface{}) *Session {
	session.Header = header
	return session
}

func (session *Session) SetRequest(request interface{}) *Session {
	session.Request = request
	return session
}

func (session *Session) SetErrorMessage(errorMessage string) *Session {
	session.ErrorMessage = errorMessage
	return session
}

func (session *Session) Get(key string) (data interface{}, err error) {
	data, ok := session.Map.Get(key)
	if !ok {
		err = errors.New("not found")
	}
	return
}

func (session *Session) Put(key string, data interface{}) {
	session.Map.Set(key, data)
}

func (session *Session) T1(message ...interface{}) {
	session.Logger.Info("|",
		zap.String("_app_tag", "T1"),
		zap.String("_app_thread_id", session.ThreadID),
		zap.String("_app_method", session.Method),
		zap.String("_app_uri", session.URL),
		zap.String("_message", formatResponse(message...)),
	)
}

func (session *Session) T2(message ...interface{}) time.Time {
	session.Logger.Info("|",
		zap.String("_app_tag", "T2"),
		zap.String("_app_thread_id", session.ThreadID),
		zap.String("_app_method", session.Method),
		zap.String("_app_uri", session.URL),
		zap.String("_message", formatResponse(message...)),
	)

	return time.Now()
}

func (session *Session) T3(startProcessTime time.Time, message ...interface{}) {
	stop := time.Now()

	session.Logger.Info("|",
		zap.String("_app_tag", "T3"),
		zap.String("_app_thread_id", session.ThreadID),
		zap.String("_app_method", session.Method),
		zap.String("_app_uri", session.URL),
		zap.String("_message", formatResponse(message...)),
		zap.String("_process_time", fmt.Sprintf("%d ms", stop.Sub(startProcessTime).Nanoseconds()/1000000)),
	)
}

func (session *Session) T4(message ...interface{}) {
	stop := time.Now()
	rt := stop.Sub(session.RequestTime).Nanoseconds() / 1000000

	session.Logger.Info("|",
		zap.String("_app_tag", "T4"),
		zap.String("_app_thread_id", session.ThreadID),
		zap.String("_app_method", session.Method),
		zap.String("_app_uri", session.URL),
		zap.String("_message", formatResponse(message...)),
		zap.String("_response_time", fmt.Sprintf("%d ms", rt)),
	)

	session.Logger.TDR(Logger.LogTdrModel{
		AppName:        session.AppName,
		AppVersion:     session.AppVersion,
		IP:             session.IP,
		Port:           session.Port,
		SrcIP:          session.SrcIP,
		RespTime:       rt,
		Path:           session.URL,
		Header:         session.Header,
		Request:        session.Request,
		Response:       formatResponse(message...),
		Error:          session.ErrorMessage,
		ThreadID:       session.ThreadID,
		AdditionalData: session.Map,
	})
}

func (session *Session) Info(message ...interface{}) {
	session.Logger.Info("|",
		zap.String("_app_tag", "INFO"),
		zap.String("_app_thread_id", session.ThreadID),
		zap.String("_app_method", session.Method),
		zap.String("_app_uri", session.URL),
		zap.String("_message", formatResponse(message...)),
	)
}

func (session *Session) Error(message ...interface{}) {
	session.Logger.Error("|",
		zap.String("_app_tag", "ERROR"),
		zap.String("_app_thread_id", session.ThreadID),
		zap.String("_app_method", session.Method),
		zap.String("_app_uri", session.URL),
		zap.String("_message", formatResponse(message...)),
	)
}

var json = JsonIter.ConfigCompatibleWithStandardLibrary

func formatResponse(message ...interface{}) string {
	sb := strings.Builder{}

	for _, msg := range message {
		var m []byte
		if reflect.ValueOf(msg).Kind().String() == "string" {
			m = []byte(msg.(string))
		} else {
			m, _ = json.Marshal(msg)
		}

		sb.Write(m)
	}

	return sb.String()
}

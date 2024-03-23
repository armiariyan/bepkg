package logger

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type Coordinate struct {
	Latitude  float32
	Longitude float32
}

type Location struct {
	Address    string
	Coordinate Coordinate
}

type Name struct {
	First string
	Last  string
}

type Data struct {
	Say      string
	Name     Name
	Location Location
}

type Response struct {
	Success bool
	Message string
	Data    Data
}

type CoordinateJSON struct {
	Latitude  float32 `json:"lat"`
	Longitude float32 `json:"lng"`
}

type LocationJSON struct {
	Address    string         `json:"string"`
	Coordinate CoordinateJSON `json:"coordinates"`
}

type NameJSON struct {
	First string `json:"first"`
	Last  string `json:"last"`
}

type DataJSON struct {
	Say      string       `json:"say"`
	Name     NameJSON     `json:"name"`
	Location LocationJSON `json:"location"`
}

type ResponseJSON struct {
	Success bool     `json:"success"`
	Message string   `json:"message"`
	Data    DataJSON `json:"data"`
}

func newOptions(fileLocation, fileTdrLocation string, fileMaxAge time.Duration, stdout bool) (options Options) {
	options = Options{
		FileLocation:    fileLocation,
		FileTdrLocation: fileTdrLocation,
		FileMaxAge:      fileMaxAge,
		Stdout:          stdout,
	}

	return
}

func newLogger(options Options) (logger Logger) {
	logger = New(options)
	return
}

func newTDR(appName, appVersion, ip, srcIP, path string, port int, respTime int64, header, request, response interface{}) (logTdr LogTdrModel) {
	logTdr = LogTdrModel{
		AppName:    appName,
		AppVersion: appVersion,
		IP:         ip,
		Port:       port,
		RespTime:   respTime,
		Path:       path,
		Header:     header,
		Request:    request,
		Response:   response,
		SrcIP:      srcIP,
	}

	return
}

func newTDRLogger(fileLocation, fileTdrLocation string, fileMaxAge time.Duration, stdout bool,
	appName, appVersion, ip, srcIP, path string, port int, respTime int64, header, request, response interface{}) (options Options, logTdr LogTdrModel, logger Logger) {
	options = newOptions(fileLocation, fileTdrLocation, fileMaxAge, stdout)
	logger = newLogger(options)
	logTdr = newTDR(appName, appVersion, ip, srcIP, path, port, respTime, header, request, response)
	return
}

func TestTDRLogger(t *testing.T) {
	assert := assert.New(t)

	// logger options
	fileLocation := "log.log"
	fileTdrLocation := "tdr.log"
	fileMaxAge := time.Minute * 5
	stdout := true

	// tdr values

	// set headers using map
	headers := make(map[string]string)
	headers["a"] = "b"
	headers["c"] = "d"

	appName := "Testing"
	appVersion := "v0.0.0"
	ip := "127.0.0.1"
	srcIP := "0.0.0.0"
	port := 80
	respTime := int64(17)
	path := "/v1/check/health"
	request := `{"action": "hello"}`
	response := `{"success": true, "message": null, "data": {"say": "Hello, World!", "name": {"first": "Bias", "last": "Tegaralaga"}, "location": {"address": "Jakarta", "coordinates": {"lat": 0.0, "lng": 0.1}}}}`

	responseStruct := Response{}

	responseStructJSON := ResponseJSON{}

	options, logTdr, logger := newTDRLogger(fileLocation, fileTdrLocation, fileMaxAge, stdout,
		appName, appVersion, ip, srcIP, path, port, respTime, headers, request, response)

	logTdrStruct := newTDR(appName, appVersion, ip, srcIP, path, port, respTime, headers, request, responseStruct)

	logTdrStructJSON := newTDR(appName, appVersion, ip, srcIP, path, port, respTime, headers, request, responseStructJSON)

	assert.Equal(fileLocation, options.FileLocation)
	assert.Equal(fileTdrLocation, options.FileTdrLocation)
	assert.Equal(fileMaxAge, options.FileMaxAge)
	assert.Equal(stdout, options.Stdout)

	assert.Equal(appName, logTdr.AppName)
	assert.Equal(appVersion, logTdr.AppVersion)
	assert.Equal(ip, logTdr.IP)
	assert.Equal(port, logTdr.Port)
	assert.Equal(respTime, logTdr.RespTime)
	assert.Equal(path, logTdr.Path)
	assert.Equal(headers, logTdr.Header)
	assert.Equal(request, logTdr.Request)
	assert.Equal(response, logTdr.Response)

	logger.TDR(logTdr)

	logger.TDR(logTdrStruct)

	logger.TDR(logTdrStructJSON)

}

func BenchmarkTDRLogger(b *testing.B) {

	fileLocation := "log.log"
	fileTdrLocation := "tdr.log"
	fileMaxAge := time.Minute * 5
	stdout := false

	logger := newLogger(newOptions(fileLocation, fileTdrLocation, fileMaxAge, stdout))

	for i := 0; i < b.N; i++ {

		// tdr values

		// set headers using map
		headers := make(map[string]string)
		headers["a"] = "b"
		headers["c"] = "d"

		appName := "Testing"
		appVersion := "v0.0.0"
		ip := "127.0.0.1"
		srcIP := "0.0.0.0"
		port := 80
		respTime := int64(17)
		path := "/v1/check/health"
		request := `{"action": "hello"}`
		response := `{"success": true, "message": null, "data": {"say": "Hello, World!"}}`

		logger.TDR(newTDR(appName, appVersion, ip, srcIP, path, port, respTime, headers, request, response))
	}
}

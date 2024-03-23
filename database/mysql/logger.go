package mysql

import (
	"database/sql/driver"
	"fmt"
	"log"
	"os"
	"reflect"
	"regexp"
	"runtime"
	"time"
	"unicode"
)

var (
	defaultLogger            = Logger{log.New(os.Stdout, "\r\n", 0)}
	sqlRegexp                = regexp.MustCompile(`\?`)
	numericPlaceHolderRegexp = regexp.MustCompile(`\$\d+`)
)

var goSrcRegexp = regexp.MustCompile(`source/.*.go`)
var goTestRegexp = regexp.MustCompile(`source/.*test.go`)

// NowFunc : time func
var NowFunc = func() time.Time {
	return time.Now()
}

func isPrintable(s string) bool {
	for _, r := range s {
		if !unicode.IsPrint(r) {
			return false
		}
	}
	return true
}

// LogFormatter : log formatter
var LogFormatter = func(values ...interface{}) (messages []interface{}) {
	if len(values) > 1 {
		var (
			sql             string
			formattedValues []string
			level           = values[0]
			currentTime     = "\n\033[33m[" + NowFunc().Format("2006-01-02 15:04:05") + "]\033[0m"
			source          = fmt.Sprintf("\033[35m(%v)\033[0m", values[1])
		)

		messages = []interface{}{source, currentTime}

		if level != "sql" {
			messages = append(messages, "\033[31;1m")
			messages = append(messages, values[2:]...)
			messages = append(messages, "\033[0m")

			return
		}

		for _, value := range values[3].([]interface{}) {
			indirectValue := reflect.Indirect(reflect.ValueOf(value))
			value = indirectValue.Interface()

			if !indirectValue.IsValid() {
				formattedValues = append(formattedValues, "NULL")
			}

			switch value.(type) {
			case time.Time:
				formattedValues = append(formattedValues, fmt.Sprintf("'%v'", value.(time.Time).Format("2006-01-02 15:04:05")))
			case []byte:
				if str := string(value.([]byte)); isPrintable(str) {
					formattedValues = append(formattedValues, fmt.Sprintf("'%v'", str))
				} else {
					formattedValues = append(formattedValues, "'<binary>'")
				}
			case driver.Valuer:
				if value, err := value.(driver.Valuer).Value(); err == nil && value != nil {
					formattedValues = append(formattedValues, fmt.Sprintf("'%v'", value))
				} else {
					formattedValues = append(formattedValues, "NULL")
				}
			default:
				formattedValues = append(formattedValues, fmt.Sprintf("'%v'", value))
			}
		}

		// differentiate between $n placeholders or else treat like ?
		if numericPlaceHolderRegexp.MatchString(values[2].(string)) {
			sql = values[3].(string)
			for index, value := range formattedValues {
				placeholder := fmt.Sprintf(`\$%d([^\d]|$)`, index+1)
				sql = regexp.MustCompile(placeholder).ReplaceAllString(sql, value+"$1")
			}
		} else {
			formattedValuesLength := len(formattedValues)
			for index, value := range sqlRegexp.Split(values[2].(string), -1) {
				sql += value
				if index < formattedValuesLength {
					sql += formattedValues[index]
				}
			}
		}

		messages = append(messages, sql)
	}

	return
}

type logger interface {
	Print(v ...interface{})
}

// LogWriter log writer interface
type LogWriter interface {
	Println(v ...interface{})
}

// Logger default logger
type Logger struct {
	LogWriter
}

// Print format & print log
func (logger Logger) Print(values ...interface{}) {
	logger.Println(LogFormatter(values...)...)
}

func fileWithLineNum() string {
	for i := 2; i < 15; i++ {
		_, file, line, ok := runtime.Caller(i)
		if ok && (!goSrcRegexp.MatchString(file) || goTestRegexp.MatchString(file)) {
			return fmt.Sprintf("%v:%v", file, line)
		}
	}
	return ""
}

// Log request with defaultLogger
func Log(logMode bool, query string, args ...interface{}) {
	if logMode {
		defaultLogger.Print("sql", fileWithLineNum(), query, args)
	}
}

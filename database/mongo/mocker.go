package mongo

import (
	"context"
	"path"
	"reflect"
	"runtime"
	"strings"
	"time"

	"gitlab.com/gobang/bepkg/logger"
)

//Mock mocking mongo client
type Mock struct {
	stubs []mockStub
}

//MockResult mock result
type MockResult struct {
	Error         error
	InsertedIDs   []string
	InsertedID    string
	DocumentCount int64
	UpsertedID    string
	DeletedCount  int64
	ModifiedCount int64
	Document      interface{}
	Documents     []interface{}
	OutputVal     interface{}
}

type mockStub struct {
	Name        string
	FilterMatch interface{}
	MockResult
}

//SetLogger implementation
func (conn *Mock) SetLogger(l logger.Logger) {}

//Debug implementation
func (conn *Mock) Debug(d bool) {}

//DB implementation
func (conn *Mock) DB(dbName string) *connection { return nil }

//Collection implementation
func (conn *Mock) Collection(collectionName string) (err error) { return nil }

//Disconnect implementation
func (conn *Mock) Disconnect() (err error) { return nil }

//GetContext implementation
func (conn *Mock) GetContext() context.Context { return context.TODO() }

//SetContext implementation
func (conn *Mock) SetContext(c context.Context) {}

//WithTimeout implementation
func (conn *Mock) WithTimeout(timeSec time.Duration) context.CancelFunc { return func() {} }

func (conn *Mock) getCaller(level int) string {
	var callerFunc string
	pc, _, _, ok := runtime.Caller(level)
	d := runtime.FuncForPC(pc)

	if ok && d != nil {
		callerFunc = path.Base(d.Name())
	}

	if callerFunc != "" {
		f := strings.Split(callerFunc, ".")
		callerFunc = f[len(f)-1]
	}

	return callerFunc
}

func (conn *Mock) stubMatchCaller(v mockStub) bool {
	c := conn.getCaller(2)
	if v.Name == c {
		return true
	}
	return false
}

func (conn *Mock) setOutputVal(outputVal interface{}, setVal interface{}) (err error) {
	value := reflect.ValueOf(outputVal)
	if value.Kind() != reflect.Ptr {
		err = ErrOutputValNotPointer
		return
	}
	direct := reflect.Indirect(value)
	direct.Set(reflect.ValueOf(setVal))
	return
}

//Stub stub any func name
func (conn *Mock) Stub(funcName string, filterMatch interface{}, result MockResult) {

	conn.stubs = append(conn.stubs, mockStub{
		Name:        funcName,
		FilterMatch: filterMatch,
		MockResult:  result,
	})
}

//Find mock func
func (conn *Mock) Find(filter interface{}, outputVal interface{}, opts ...*Options) (err error) {
	for _, v := range conn.stubs {
		if conn.stubMatchCaller(v) && reflect.DeepEqual(v.FilterMatch, filter) {
			conn.setOutputVal(outputVal, v.OutputVal)
			return v.Error
		}
	}
	panic("Unsetted")
}

//FindOne mock func
func (conn *Mock) FindOne(filter interface{}, outputVal interface{}, opts ...*Options) error {
	for _, v := range conn.stubs {
		if conn.stubMatchCaller(v) && reflect.DeepEqual(v.FilterMatch, filter) {
			conn.setOutputVal(outputVal, v.OutputVal)
			return v.Error
		}
	}
	panic("Unsetted")
}

//FindOneAndDelete mock func
func (conn *Mock) FindOneAndDelete(filter interface{}, outputVal interface{}, opts ...*Options) error {
	for _, v := range conn.stubs {
		if conn.stubMatchCaller(v) && reflect.DeepEqual(v.FilterMatch, filter) {
			outputVal = v.OutputVal
			return v.Error
		}
	}
	panic("Unsetted")
}

//FindOneAndReplace mock func
func (conn *Mock) FindOneAndReplace(filter, replacement interface{}, outputVal interface{}, opts ...*Options) error {
	for _, v := range conn.stubs {
		if conn.stubMatchCaller(v) && reflect.DeepEqual(v.FilterMatch, filter) {
			outputVal = v.OutputVal
			return v.Error
		}
	}
	panic("Unsetted")
}

//FindOneAndUpdate mock func
func (conn *Mock) FindOneAndUpdate(filter, update interface{}, outputVal interface{}, opts ...*Options) error {
	for _, v := range conn.stubs {
		if conn.stubMatchCaller(v) && reflect.DeepEqual(v.FilterMatch, filter) {
			outputVal = v.OutputVal
			return v.Error
		}
	}
	panic("Unsetted")
}

//InsertMany mock func
func (conn *Mock) InsertMany(documents []interface{}) (insertedIDs []string, err error) {
	for _, v := range conn.stubs {
		if conn.stubMatchCaller(v) {
			docMatch := true
			for i := 0; i < len(documents); i++ {
				if !reflect.DeepEqual(v.Documents[i], documents[i]) {
					docMatch = false
					break
				}
			}
			if docMatch {
				return v.InsertedIDs, v.Error
			}
		}
	}
	panic("Unsetted")
}

//InsertOne mock func
func (conn *Mock) InsertOne(document interface{}) (insertedID string, err error) {
	for _, v := range conn.stubs {
		if conn.stubMatchCaller(v) && reflect.DeepEqual(v.Document, document) {
			return v.InsertedID, v.Error
		}
	}
	panic("Unsetted")
}

//UpdateMany mock func
func (conn *Mock) UpdateMany(filter, update interface{}) (modifiedCount int64, err error) {
	for _, v := range conn.stubs {
		if conn.stubMatchCaller(v) && reflect.DeepEqual(v.FilterMatch, filter) {
			return v.ModifiedCount, v.Error
		}
	}
	panic("unsetted")
}

//UpdateOne mock func
func (conn *Mock) UpdateOne(filter, update interface{}, opts ...*Options) (upsertedID string, modifiedCount int64, err error) {
	for _, v := range conn.stubs {
		if conn.stubMatchCaller(v) && reflect.DeepEqual(v.FilterMatch, filter) {
			return v.UpsertedID, v.ModifiedCount, v.Error
		}
	}
	panic("unsetted")
}

//ReplaceOne mock func
func (conn *Mock) ReplaceOne(filter, replacement interface{}, opts ...*Options) (upsertedID string, modifiedCount int64, err error) {
	for _, v := range conn.stubs {
		if conn.stubMatchCaller(v) && reflect.DeepEqual(v.FilterMatch, filter) {
			return v.UpsertedID, v.ModifiedCount, v.Error
		}
	}
	panic("unsetted")
}

//CountDocuments mock func
func (conn *Mock) CountDocuments(filter interface{}, opts ...*Options) (total int64, err error) {
	for _, v := range conn.stubs {
		if conn.stubMatchCaller(v) && reflect.DeepEqual(v.FilterMatch, filter) {
			return v.DocumentCount, v.Error
		}
	}
	panic("unsetted")
}

//DeleteOne mock func
func (conn *Mock) DeleteOne(filter interface{}) (DeletedCount int64, err error) {
	for _, v := range conn.stubs {
		if conn.stubMatchCaller(v) && reflect.DeepEqual(v.FilterMatch, filter) {
			return v.DeletedCount, v.Error
		}
	}
	panic("unsetted")
}

//DeleteMany mock func
func (conn *Mock) DeleteMany(filter interface{}) (DeletedCount int64, err error) {
	for _, v := range conn.stubs {
		if conn.stubMatchCaller(v) && reflect.DeepEqual(v.FilterMatch, filter) {
			return v.DeletedCount, v.Error
		}
	}
	panic("unsetted")
}

//Aggregate mock func
func (conn *Mock) Aggregate(pipeline interface{}, outputVal interface{}) (err error) {
	for _, v := range conn.stubs {
		if conn.stubMatchCaller(v) && reflect.DeepEqual(v.FilterMatch, pipeline) {
			outputVal = v.OutputVal
			return v.Error
		}
	}
	panic("unsetted")
}

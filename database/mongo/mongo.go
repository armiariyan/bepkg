package mongo

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	mo "go.mongodb.org/mongo-driver/mongo"
	moOpts "go.mongodb.org/mongo-driver/mongo/options"

	"gitlab.com/gobang/bepkg/logger"
)

type (
	//M map filter type
	M map[string]interface{}
	//Client mongo interface for mongo client connection
	Client interface {
		Disconnect() error
		DB(dbName string) *connection
		Collection(collectionName string) (err error)
		SetLogger(l logger.Logger)
		Debug(d bool)
		GetContext() context.Context
		SetContext(c context.Context)
		WithTimeout(timeSec time.Duration) context.CancelFunc
		StartSession(c context.Context) (mo.Session, error)

		//Find(filter interface{}, outputVal interface{}, opts ...*Options) (results []interface{}, err error)
		Find(filter interface{}, outputVal interface{}, opts ...*Options) (err error)
		FindOne(filter interface{}, outputVal interface{}, opts ...*Options) error
		FindOneAndDelete(filter interface{}, outputVal interface{}, opts ...*Options) error
		FindOneAndReplace(filter, replacement interface{}, outputVal interface{}, opts ...*Options) error
		FindOneAndUpdate(filter, update interface{}, outputVal interface{}, opts ...*Options) error

		InsertMany(documents []interface{}) (insertedIDs []string, err error)
		InsertOne(document interface{}) (insertedID string, err error)

		UpdateMany(filter, update interface{}) (modifiedCount int64, err error)
		UpdateOne(filter, update interface{}, opts ...*Options) (upsertedID string, modifiedCount int64, err error)
		ReplaceOne(filter, replacement interface{}, opts ...*Options) (upsertedID string, modifiedCount int64, err error)

		CountDocuments(filter interface{}, opts ...*Options) (total int64, err error)

		DeleteOne(filter interface{}) (DeletedCount int64, err error)
		DeleteMany(filter interface{}) (DeletedCount int64, err error)

		Aggregate(pipeline interface{}, outputVal interface{}) (err error)
	}

	connection struct {
		ctx        context.Context
		client     *mo.Client
		collection *mo.Collection
		dbName     string
		logger     logger.Logger
		debug      bool
	}
)

func (conn *connection) DB(dbName string) *connection {
	conn.dbName = dbName
	return conn
}

func (conn *connection) Collection(collectionName string) (err error) {
	if conn.dbName == "" {
		err = errors.New("DB required")
		return
	}
	conn.collection = conn.client.Database(conn.dbName).Collection(collectionName)
	return
}

// func (conn *connection) ObjectID(oid string) primitive.ObjectID {
// 	v, _ := primitive.ObjectIDFromHex(oid)
// 	return v
// }

//Disconnect terminate connection with mongo client
func (conn *connection) Disconnect() (err error) {
	err = conn.client.Disconnect(conn.ctx)
	return
}

//GetContext get connection context
func (conn *connection) GetContext() context.Context {
	return conn.ctx
}

//SetContext set connection context
func (conn *connection) SetContext(c context.Context) {
	conn.ctx = c
}

//WithTimeout set timeout based on context
func (conn *connection) WithTimeout(timeSec time.Duration) context.CancelFunc {
	ctx, cancel := context.WithTimeout(conn.ctx, timeSec*time.Second)
	conn.ctx = ctx
	return cancel
}

//WithTimeout set timeout based on context
func (conn *connection) StartSession(c context.Context) (mo.Session, error) {
	return conn.client.StartSession()
}

//Connect start a connection with mongodb based on uri
//ex URI: mongodb://localhost:27017
func Connect(ctx context.Context, URI string, opts ...ClientOptions) (mongoClient Client, err error) {

	clientOptions := moOpts.Client().ApplyURI(URI)

	if len(opts) > 0 {
		if opts[0].MaxConnIdleTime > 0 {
			clientOptions.SetMaxConnIdleTime(opts[0].MaxConnIdleTime)
		}
		if opts[0].MaxPoolSize > 0 {
			clientOptions.SetMaxPoolSize(opts[0].MaxPoolSize)
		}
		if opts[0].MinPoolSize > 0 {
			clientOptions.SetMinPoolSize(opts[0].MinPoolSize)
		}
		if opts[0].Auth.Username != "" {
			credential := moOpts.Credential{
				AuthMechanism: opts[0].Auth.AuthMechanism,
				Username:      opts[0].Auth.Username, // your mongodb user
				Password:      opts[0].Auth.Password, // ...and mongodb
			}
			clientOptions.SetAuth(credential)
		}
	}

	// Connect to MongoDB
	moClient, err := mo.Connect(ctx, clientOptions)
	if err != nil {
		return
	}

	//Check the connection
	err = moClient.Ping(ctx, nil)
	if err != nil {
		return
	}

	mongoClient = &connection{
		ctx:    ctx,
		client: moClient,
	}

	return
}

//ToTime convert mongo primitive datetime to native time.Time, return zero time if invalid
//you may check zero time with Time.IsZero() function
func ToTime(t interface{}) (time.Time, error) {
	if s, ok := t.(primitive.DateTime); ok {
		return s.Time(), nil
	}

	return time.Time{}, ErrInvalidTimeType
}

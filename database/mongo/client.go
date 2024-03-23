package mongo

import (
	"reflect"
	"time"

	"github.com/jinzhu/copier"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mo "go.mongodb.org/mongo-driver/mongo"
	moOpts "go.mongodb.org/mongo-driver/mongo/options"
)

func deref(t reflect.Type) reflect.Type {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t
}

//Find finds the documents matching a model.
func (conn *connection) Find(filter interface{}, outputVal interface{}, opts ...*Options) (err error) {
	var Opts *moOpts.FindOptions
	var cur *mo.Cursor

	startTime := time.Now()

	value := reflect.ValueOf(outputVal)
	if value.Kind() != reflect.Ptr {
		err = ErrOutputValNotPointer
		return
	}
	direct := reflect.Indirect(value)
	slice := deref(value.Type())
	if slice.Kind() != reflect.Slice {
		err = ErrOutputValNotSlicePointer
		return
	}
	isPtr := slice.Elem().Kind() == reflect.Ptr
	base := deref(slice.Elem())

	if opts != nil {
		Opts = new(moOpts.FindOptions)
		copier.Copy(Opts, opts[0])
		cur, err = conn.collection.Find(conn.ctx, filter, Opts)
	} else {
		cur, err = conn.collection.Find(conn.ctx, filter)
	}

	if err != nil {
		conn.logError(startTime, filter, opts, err)
		return
	}

	for cur.Next(conn.ctx) {
		//Note that reflect.New() will create a new value which is initialized
		//to its zero value (so it will not be a copy of the original)
		//outputVal must be a pointer
		//if not a pointer --> reflect.New(reflect.TypeOf(outpuVal)).Elem().Interface()
		//result := reflect.New(reflect.ValueOf(outputVal).Elem().Type()).Interface()

		vp := reflect.New(base)

		// Create a value into which the single document can be decoded
		err := cur.Decode(vp.Interface())
		if err != nil {
			conn.logError(startTime, filter, opts, err)
			return err
		}

		if isPtr {
			direct.Set(reflect.Append(direct, vp))
		} else {
			direct.Set(reflect.Append(direct, reflect.Indirect(vp)))
		}

	}

	if err := cur.Err(); err != nil {
		cur.Close(conn.ctx)
		conn.logError(startTime, filter, opts, err)
		return err
	}

	// Close the cursor once finished
	cur.Close(conn.ctx)

	conn.logInfo(startTime, filter, outputVal, opts)

	return
}

//FindOne returns up to one document that matches the model.
func (conn *connection) FindOne(filter interface{}, outputVal interface{}, opts ...*Options) error {
	var Opts *moOpts.FindOneOptions
	var res *mo.SingleResult
	startTime := time.Now()
	if opts != nil {
		Opts = new(moOpts.FindOneOptions)
		copier.Copy(Opts, opts[0])
		res = conn.collection.FindOne(conn.ctx, filter, Opts)
	} else {
		res = conn.collection.FindOne(conn.ctx, filter)
	}
	if res.Err() != nil {
		conn.logError(startTime, filter, opts, res.Err())
		return res.Err()
	}

	if err := res.Decode(outputVal); err != nil {
		conn.logError(startTime, filter, opts, err)
		return err
	}

	conn.logInfo(startTime, filter, opts, outputVal)

	return nil
}

//FindOneAndDelete find a single document and deletes it, returning the original in result.
func (conn *connection) FindOneAndDelete(filter interface{}, outputVal interface{}, opts ...*Options) error {
	var Opts *moOpts.FindOneAndDeleteOptions
	var res *mo.SingleResult
	startTime := time.Now()
	if opts != nil {
		Opts = new(moOpts.FindOneAndDeleteOptions)
		copier.Copy(Opts, opts[0])
		res = conn.collection.FindOneAndDelete(conn.ctx, filter, Opts)
	} else {
		res = conn.collection.FindOneAndDelete(conn.ctx, filter)
	}

	if res.Err() != nil {
		conn.logError(startTime, filter, opts, res.Err)
		return res.Err()
	}

	if err := res.Decode(outputVal); err != nil {
		conn.logError(startTime, filter, opts, err)
		return err
	}

	conn.logInfo(startTime, filter, opts, outputVal)

	return nil

}

//FindOneAndReplace finds a single document and replaces it, returning either the original or the replaced document.
func (conn *connection) FindOneAndReplace(filter, replacement interface{}, outputVal interface{}, opts ...*Options) error {
	var Opts *moOpts.FindOneAndReplaceOptions
	var res *mo.SingleResult
	startTime := time.Now()
	if opts != nil {
		Opts = new(moOpts.FindOneAndReplaceOptions)
		copier.Copy(Opts, opts[0])
		res = conn.collection.FindOneAndReplace(conn.ctx, filter, replacement, Opts)
	} else {
		res = conn.collection.FindOneAndReplace(conn.ctx, filter, replacement)
	}

	if res.Err() != nil {
		conn.logError(startTime, filter, replacement, opts, res.Err())
		return res.Err()
	}

	if err := res.Decode(outputVal); err != nil {
		conn.logError(startTime, filter, replacement, opts, err)
		return err
	}

	conn.logInfo(startTime, filter, replacement, opts, outputVal)

	return nil

}

//FindOneAndUpdate finds a single document and updates it, returning either the original or the updated.
func (conn *connection) FindOneAndUpdate(filter, update interface{}, outputVal interface{}, opts ...*Options) error {
	var Opts *moOpts.FindOneAndUpdateOptions
	var res *mo.SingleResult
	startTime := time.Now()
	if opts != nil {
		Opts = new(moOpts.FindOneAndUpdateOptions)
		copier.Copy(Opts, opts[0])
		res = conn.collection.FindOneAndUpdate(conn.ctx, filter, update, Opts)
	} else {
		res = conn.collection.FindOneAndUpdate(conn.ctx, filter, update)
	}

	if res.Err() != nil {
		conn.logError(startTime, filter, update, opts, res.Err())
		return res.Err()
	}

	if err := res.Decode(outputVal); err != nil {
		return err
	}

	conn.logInfo(startTime, filter, update, opts, outputVal)
	return nil

}

func (conn *connection) InsertMany(documents []interface{}) (insertedIDs []string, err error) {
	var result *mo.InsertManyResult
	startTime := time.Now()
	result, err = conn.collection.InsertMany(conn.ctx, documents)
	if err != nil {
		conn.logError(startTime, documents, err)
		return
	}
	if len(result.InsertedIDs) > 0 {
		for _, v := range result.InsertedIDs {
			insertedIDs = append(insertedIDs, v.(primitive.ObjectID).Hex())
		}
	}

	conn.logInfo(startTime, documents, insertedIDs)

	return

}

//InsertOne inserts a single document into the collection
func (conn *connection) InsertOne(document interface{}) (insertedID string, err error) {
	var result *mo.InsertOneResult
	startTime := time.Now()
	result, err = conn.collection.InsertOne(conn.ctx, document)
	if err != nil {
		conn.logError(startTime, document, err)
		return
	}

	insertedID = result.InsertedID.(primitive.ObjectID).Hex()

	conn.logInfo(startTime, document, insertedID)
	return
}

//UpdateMany updates multiple documents in the collection.
func (conn *connection) UpdateMany(filter, update interface{}) (modifiedCount int64, err error) {
	var result *mo.UpdateResult
	startTime := time.Now()
	result, err = conn.collection.UpdateMany(conn.ctx, filter, update)
	if err != nil {
		conn.logError(startTime, filter, update, err)
		return
	}
	modifiedCount = result.ModifiedCount

	conn.logInfo(startTime, filter, update, modifiedCount)
	return
}

//UpdateOne updates a single document in the collection.
func (conn *connection) UpdateOne(filter, update interface{}, opts ...*Options) (upsertedID string, modifiedCount int64, err error) {
	var result *mo.UpdateResult

	var Opts *moOpts.UpdateOptions
	startTime := time.Now()
	if opts != nil {
		Opts = new(moOpts.UpdateOptions)
		copier.Copy(Opts, opts[0])
		result, err = conn.collection.UpdateOne(conn.ctx, filter, update, Opts)
	} else {
		result, err = conn.collection.UpdateOne(conn.ctx, filter, update)
	}

	if err != nil {
		conn.logError(startTime, filter, update, err)
		return
	}
	modifiedCount = result.ModifiedCount
	if result.UpsertedID != nil {
		upsertedID = result.UpsertedID.(primitive.ObjectID).Hex()
	}
	conn.logInfo(startTime, filter, update, upsertedID, modifiedCount)
	return
}

//ReplaceOne replaces a single document in the collection.
func (conn *connection) ReplaceOne(filter, replacement interface{}, opts ...*Options) (upsertedID string, modifiedCount int64, err error) {
	var result *mo.UpdateResult
	var Opts *moOpts.ReplaceOptions
	startTime := time.Now()
	if opts != nil {
		Opts = new(moOpts.ReplaceOptions)
		copier.Copy(Opts, opts[0])
		result, err = conn.collection.ReplaceOne(conn.ctx, filter, replacement, Opts)
	} else {
		result, err = conn.collection.ReplaceOne(conn.ctx, filter, replacement)
	}

	if err != nil {
		conn.logError(startTime, filter, replacement, err)
		return
	}
	modifiedCount = result.ModifiedCount
	if result.UpsertedID != nil {
		upsertedID = result.UpsertedID.(primitive.ObjectID).Hex()
	}
	conn.logInfo(startTime, filter, replacement, upsertedID, modifiedCount)

	return

}

//CountDocuments gets the number of documents matching the filter.
func (conn *connection) CountDocuments(filter interface{}, opts ...*Options) (total int64, err error) {

	var Opts *moOpts.CountOptions
	startTime := time.Now()
	if opts != nil {
		Opts = new(moOpts.CountOptions)
		copier.Copy(Opts, opts[0])
		total, err = conn.collection.CountDocuments(conn.ctx, filter, Opts)
	} else {
		total, err = conn.collection.CountDocuments(conn.ctx, filter)
	}

	if err != nil {
		conn.logError(startTime, filter, opts, err)
	}

	conn.logInfo(startTime, filter, opts, total)

	return
}

//DeleteOne deletes a single document from the collection.
func (conn *connection) DeleteOne(filter interface{}) (DeletedCount int64, err error) {
	var result *mo.DeleteResult
	startTime := time.Now()
	result, err = conn.collection.DeleteOne(conn.ctx, filter)
	if err != nil {
		conn.logError(startTime, filter, err)
		return
	}
	DeletedCount = result.DeletedCount

	conn.logInfo(startTime, filter, DeletedCount)
	return
}

//DeleteMany deletes multiple documents from the collection.
func (conn *connection) DeleteMany(filter interface{}) (DeletedCount int64, err error) {
	var result *mo.DeleteResult
	startTime := time.Now()
	result, err = conn.collection.DeleteMany(conn.ctx, filter)
	if err != nil {
		conn.logError(startTime, filter, err)
		return
	}
	DeletedCount = result.DeletedCount
	conn.logInfo(startTime, filter, DeletedCount)
	return

}

//Aggregate runs an aggregation framework pipeline.
func (conn *connection) Aggregate(pipeline interface{}, outputVal interface{}) (err error) {
	var cur *mo.Cursor
	startTime := time.Now()
	value := reflect.ValueOf(outputVal)
	if value.Kind() != reflect.Ptr {
		err = ErrOutputValNotPointer
		return
	}
	direct := reflect.Indirect(value)
	slice := deref(value.Type())
	if slice.Kind() != reflect.Slice {
		err = ErrOutputValNotSlicePointer
		return
	}
	isPtr := slice.Elem().Kind() == reflect.Ptr
	base := deref(slice.Elem())

	cur, err = conn.collection.Aggregate(conn.ctx, pipeline)
	if err != nil {
		conn.logError(startTime, pipeline, err)
		return
	}

	for cur.Next(conn.ctx) {
		vp := reflect.New(base)

		// Create a value into which the single document can be decoded
		err := cur.Decode(vp.Interface())
		if err != nil {
			conn.logError(startTime, pipeline, err)
			return err
		}

		if isPtr {
			direct.Set(reflect.Append(direct, vp))
		} else {
			direct.Set(reflect.Append(direct, reflect.Indirect(vp)))
		}
	}

	if err := cur.Err(); err != nil {
		cur.Close(conn.ctx)
		conn.logError(startTime, pipeline, err)
		return err
	}

	// Close the cursor once finished
	cur.Close(conn.ctx)
	conn.logInfo(startTime, pipeline, outputVal)
	return

}

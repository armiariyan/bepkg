package mongo

import (
	"context"
	"fmt"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"gitlab.com/gobang/bepkg/logger"
)

type Trainer struct {
	Name string `bson:"name"`
	Age  int    `bson:"age"`
	City string `bson:"city"`
}

type Pelatih struct {
	Nama string
	Umur int
	Kota string
}

func connMock(t *testing.T) (client Client) {
	var err error
	client, err = Connect(context.TODO(), "mongodb://0.0.0.0:27001,0.0.0.0:27002,0.0.0.0:27003/?replicaSet=mongo-rs", ClientOptions{
		MaxPoolSize: 5000,
	})
	if err != nil {
		t.Errorf("Something happen, %s", err)
	}
	client.DB("test").Collection("trainers")

	logr := logger.New(logger.Options{
		Stdout: true,
	})
	client.Debug(true)

	client.SetLogger(logr)
	return
}

func TestConnect(t *testing.T) {
	connMock(t)
}

func TestInsertOne(t *testing.T) {
	client := connMock(t)

	ash := Trainer{"Ash", 10, "Pallet Town"}
	misty := Trainer{"Misty", 10, "Cerulean City"}
	brock := Trainer{"Brock", 15, "Pewter City"}

	id, err := client.InsertOne(ash)
	fmt.Println(id)
	if err != nil {
		t.Errorf("Something happen, %s", err)
	}

	id, err = client.InsertOne(misty)
	fmt.Println(id)
	if err != nil {
		t.Errorf("Something happen, %s", err)
	}

	id, err = client.InsertOne(brock)
	fmt.Println(id)
	if err != nil {
		t.Errorf("Something happen, %s", err)
	}
}

func TestInsertMany(t *testing.T) {
	client := connMock(t)

	ash := Trainer{"AshMany", 10, "Pallet Town"}
	misty := Trainer{"MistyMany", 10, "Cerulean City"}
	brock := Trainer{"BrockMany", 15, "Pewter City"}

	ids, err := client.InsertMany([]interface{}{ash, misty, brock})
	fmt.Println(ids)
	if err != nil {
		t.Errorf("Something happen, %s", err)
	}

}

func TestUpdateOne(t *testing.T) {
	client := connMock(t)

	_, count, err := client.UpdateOne(M{
		"name": "BrockMany",
	}, M{
		"$set": M{
			"name": "BrockManyUpdatedOne",
			"city": "Updated One Pallet Town",
		},
	})
	fmt.Printf("Update One %d\n", count)
	if err != nil {
		t.Errorf("Something happen, %s", err)
	}

}

func TestReplaceOne(t *testing.T) {
	client := connMock(t)

	_, count, err := client.ReplaceOne(M{
		"name": "AshMany",
	}, M{
		"nama": "AshMany replaced",
		"umur": 23,
		"kota": "replaced One Pallet Town",
	})
	fmt.Printf("Replaced One %d\n", count)
	if err != nil {
		t.Errorf("Something happen, %s", err)
	}

}

func TestUpdateMany(t *testing.T) {
	client := connMock(t)

	count, err := client.UpdateMany(M{
		"name": "MistyMany",
	}, M{
		"$set": M{
			"name": "MistyManyUpdatedMany",
			"city": "Updated Many Pallet Town",
		},
	})
	fmt.Printf("Update Many %d\n", count)
	if err != nil {
		t.Errorf("Something happen, %s", err)
	}

}

func TestFindAll(t *testing.T) {
	client := connMock(t)
	//filter := bson.D{{}}
	filter := M{
		"name": M{
			"$in": []string{
				"Brock",
			},
		},
	}
	filter = M{}
	//{"name":{"$in":["Brock"]}}
	// opts := &Options{
	// 	Limit: 5,
	// }

	//var tr Trainer
	var trs []Trainer
	err := client.Find(filter, &trs)
	if err != nil {
		t.Errorf("Something happen, %s", err)
	}
	for _, v := range trs {
		fmt.Printf("%#v\n", v)
	}
}

func TestFindOne(t *testing.T) {
	client := connMock(t)
	//oid, _ := primitive.ObjectIDFromHex("5d6b553f55f5d7f409383e9d")
	//x, err := client.FindOne(M{"_id": oid})
	var tr Trainer
	err := client.FindOne(M{"name": "Brock"}, &tr)
	if err != nil {
		t.Errorf("Something happen, %s", err)
	}
	fmt.Printf("Results %#v\n", tr)
}

func TestDeleteOne(t *testing.T) {
	client := connMock(t)

	cnt, err := client.DeleteOne(M{
		"name": "Brock",
	})
	if err != nil {
		t.Errorf("Something happen, %s", err)
	}
	fmt.Printf("Results %#v\n", cnt)

}

func TestDeleteMany(t *testing.T) {
	client := connMock(t)

	cnt, err := client.DeleteMany(M{
		"name": "MistyManyUpdatedMany",
	})
	if err != nil {
		t.Errorf("Something happen, %s", err)
	}
	fmt.Printf("Results %#v\n", cnt)

}

func TestReplaceOneUpsert(t *testing.T) {
	client := connMock(t)
	upsertedID, modCnt, err := client.ReplaceOne(M{
		"name": "BrookManyUpsertasdf",
	}, Trainer{
		Name: "BrookManyUpsertasdf",
		City: "jambi",
		Age:  100,
	}, &Options{
		Upsert: true,
	})
	if err != nil {
		t.Errorf("Something happen, %s", err)
	}

	fmt.Printf("Upserted %#v\n", upsertedID)
	fmt.Printf("Modified %#v\n", modCnt)
}

func TestUpdateOneUpsert(t *testing.T) {
	client := connMock(t)

	defer client.WithTimeout(20)()

	_, count, err := client.UpdateOne(M{
		"name": "BrockManyUpdatedUpsert",
	}, M{
		"$set": M{
			"name": "BrockManyUpdatedUpsert",
			"city": "Updated One Pallet Town",
		},
	}, &Options{
		Upsert: true,
	})
	fmt.Printf("Update One %d\n", count)
	if err != nil {
		t.Errorf("Something happen, %s", err)
	}

}

func TestToTime(t *testing.T) {
	timeNow := time.Now()

	tn := primitive.NewDateTimeFromTime(timeNow)
	x, err := ToTime(tn)
	if err != nil {
		t.Errorf("Something happen, %s", err)
	}
	if timeNow.Unix() != x.Unix() {
		t.Errorf("Timestamp missmatch, %d vs %d", timeNow.Unix(), x.Unix())
	}
	fmt.Printf("%#v\n", x)
}

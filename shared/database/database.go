package database

import (
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"server-monitoring/domain/nodes"
	loogers "server-monitoring/utils/looger"
	"sync"

	//"gopkg.in/mgo.v2"
	"log"
	"time"

	"github.com/boltdb/bolt"
	_ "github.com/go-sql-driver/mysql" // MySQL driver
	"github.com/jmoiron/sqlx"
	//"gopkg.in/mgo.v2"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	// BoltDB wrapper
	BoltDB *bolt.DB
	// Mongo wrapper
	Mongo *mongo.Client
	// SQL wrapper
	SQL *sqlx.DB
	DB  *mongo.Database
	// Database info
	databases      Info
	mutax          sync.RWMutex
	nodeCollection *mongo.Collection
)

func init() {
	mutax = sync.RWMutex{}
}

// Type is the type of database from a Type* constant
type Type string

const (
	// TypeBolt is BoltDB
	TypeBolt Type = "Bolt"
	// TypeMongoDB is MongoDB
	TypeMongoDB Type = "MongoDB"
	// TypeMySQL is MySQL
	TypeMySQL Type = "MySQL"
)

// Info contains the database configurations
type Info struct {
	// Database type
	Type Type
	// MySQL info if used
	MySQL MySQLInfo
	// Bolt info if used
	Bolt BoltInfo
	// MongoDB info if used
	MongoDB MongoDBInfo
}

// MySQLInfo is the details for the database connection
type MySQLInfo struct {
	Username  string
	Password  string
	Name      string
	Hostname  string
	Port      int
	Parameter string
}

// BoltInfo is the details for the database connection
type BoltInfo struct {
	Path string
}

// MongoDBInfo is the details for the database connection
type MongoDBInfo struct {
	URL      string
	Database string
}

// DSN returns the Data Source Name
func DSN(ci MySQLInfo) string {
	// Example: root:@tcp(localhost:3306)/test
	return ci.Username +
		":" +
		ci.Password +
		"@tcp(" +
		ci.Hostname +
		":" +
		fmt.Sprintf("%d", ci.Port) +
		")/" +
		ci.Name + ci.Parameter
}

// Connect to the database
func Connect(d Info) {
	var err error

	// Store the config
	databases = d

	switch d.Type {
	case TypeMySQL:
		// Connect to MySQL
		if SQL, err = sqlx.Connect("mysql", DSN(d.MySQL)); err != nil {
			log.Println("SQL Driver Error", err)
		}

		// Check if is alive
		if err = SQL.Ping(); err != nil {
			log.Println("Database Error", err)
		}
	case TypeBolt:
		// Connect to Bolt
		if BoltDB, err = bolt.Open(d.Bolt.Path, 0600, nil); err != nil {
			log.Println("Bolt Driver Error", err)
		}
	case TypeMongoDB:

		clientOptions := options.Client().
			ApplyURI(databases.MongoDB.URL)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		Mongo, err = mongo.Connect(ctx, clientOptions)
		if err != nil {
			log.Fatal(err)
		}
		db := Mongo.Database("monitoring")

		if err := db.CreateCollection(ctx, "users"); err != nil {
			fmt.Println(err)
		}
		if err := db.CreateCollection(ctx, "requests"); err != nil {
			fmt.Println(err)
		}
		if err := db.CreateCollection(ctx, "nodes"); err != nil {
			fmt.Println(err)
		}

		//defer func(Mongo *mongo.Client, ctx context.Context) {
		//	err := Mongo.Disconnect(ctx)
		//	if err != nil {
		//
		//	}
		//}(Mongo, ctx)
		// Connect to MongoDB
		//if Mongo, err = mgo.DialWithTimeout(d.MongoDB.URL, 5*time.Second); err != nil {
		//	log.Println("MongoDB Driver Error", err)
		//	return
		//}
		//
		//// Prevents these errors: read tcp 127.0.0.1:27017: i/o timeout
		//Mongo.SetSocketTimeout(1 * time.Second)

		// Check if is alive
		if err = Mongo.Ping(ctx, readpref.Primary()); err != nil {
			log.Println("Database Error", err)
		}

		//DB =
		//println(DB.Name())
		//if DB == nil {
		//	log.Fatal("mango not connected")
		//}
		//nodeCollection = DB.Collection("nodes")
		//
		//item := nodes.Node{
		//	ID:        primitive.ObjectID{},
		//	SrcPort:   "2",
		//	DstPort:   "2",
		//	SrcIp:     "2",
		//	DstIp:     "2",
		//	SrcMac:    "2",
		//	DstMac:    "2",
		//	Protocol:  "22",
		//	Timestamp: 0,
		//}
		//nod := bson.D{
		//	{"src_ip", item.SrcIp},
		//	{"src_mac", item.SrcMac},
		//	{"dst_ip", item.DstIp},
		//	{"dst_mac", item.DstMac},
		//	{"protocol", item.Protocol},
		//	{"timestamp", item.Timestamp}}
		//_, err := nodeCollection.InsertOne(ctx, nod)
		//
		//podcastResult, err = podcastsCollection.InsertOne(ctx, bson.D{
		//	{Key: "title", Value: "The Polyglot Developer Podcast"},
		//	{Key: "author", Value: "Nic Raboy"},
		//	{Key: "tags", Value: bson.A{"development", "programming", "coding"}},
		//})
		//episodeResult, err := episodesCollection.InsertMany(ctx, []interface{}{
		//	bson.D{
		//		{"podcast", podcastResult.InsertedID},
		//		{"title", "GraphQL for API Development"},
		//		{"description", "Learn about GraphQL from the co-creator of GraphQL, Lee Byron."},
		//		{"duration", 25},
		//	},
		//	bson.D{
		//		{"podcast", podcastResult.InsertedID},
		//		{"title", "Progressive Web Application Development"},
		//		{"description", "Learn about PWA development with Tara Manicsic."},
		//		{"duration", 32},
		//	},
		//})
		if err != nil {
			log.Fatal(err)
		}
		//fmt.Printf("Inserted %v documents into episode collection!\n", userResult.InsertedID)
	default:
		log.Println("No registered database in config")
	}
}

func InsertLogs(items []nodes.Node) error {

	if Mongo != nil {

		nodeCollection = Mongo.Database("monitoring").Collection("nodes")

		//var nodes []interface{}

		for _, item := range items {
			//test1 := map[string]interface{}{
			//	"src_ip":    item.SrcIp,
			//	"src_mac":   item.SrcMac,
			//	"dst_ip":    item.DstIp,
			//	"dst_mac":   item.DstMac,
			//	"protocol":  item.Protocol,
			//	"timestamp": item.Timestamp,
			//}
			nod := bson.D{
				{"src_ip", (item.SrcIp)},
				{"src_mac", item.SrcMac},
				{"dst_ip", item.DstIp},
				{"dst_mac", item.DstMac},
				{"protocol", item.Protocol},
				{"timestamp", item.Timestamp}}
			_, err := nodeCollection.InsertOne(context.TODO(), nod)
			if err != nil {
				loogers.Error("insert node", err)
				//return err
			}
			//nod := bson.D{
			//	{"src_ip", item.SrcIp},
			//	{"src_mac", item.SrcMac},
			//	{"dst_ip", item.DstIp},
			//	{"dst_mac", item.DstMac},
			//	{"protocol", item.Protocol},
			//	{"timestamp", item.Timestamp}}
			//nodes = append(nodes, test1)
		}

		//_, err := nodeCollection.InsertMany(context.TODO(), nodes)
		//if err != nil {
		//	loogers.Error("insert node", err)
		//	return err
		//}

	}
	return nil
}

// Update makes a modification to Bolt
func Update(bucketName string, key string, dataStruct interface{}) error {
	err := BoltDB.Update(func(tx *bolt.Tx) error {
		// Create the bucket
		bucket, e := tx.CreateBucketIfNotExists([]byte(bucketName))
		if e != nil {
			return e
		}

		// Encode the record
		encodedRecord, e := json.Marshal(dataStruct)
		if e != nil {
			return e
		}

		// Store the record
		if e = bucket.Put([]byte(key), encodedRecord); e != nil {
			return e
		}
		return nil
	})
	return err
}

// View retrieves a record in Bolt
func View(bucketName string, key string, dataStruct interface{}) error {
	err := BoltDB.View(func(tx *bolt.Tx) error {
		// Get the bucket
		b := tx.Bucket([]byte(bucketName))
		if b == nil {
			return bolt.ErrBucketNotFound
		}

		// Retrieve the record
		v := b.Get([]byte(key))
		if len(v) < 1 {
			return bolt.ErrInvalid
		}

		// Decode the record
		e := json.Unmarshal(v, &dataStruct)
		if e != nil {
			return e
		}

		return nil
	})

	return err
}

// Delete removes a record from Bolt
func Delete(bucketName string, key string) error {
	err := BoltDB.Update(func(tx *bolt.Tx) error {
		// Get the bucket
		b := tx.Bucket([]byte(bucketName))
		if b == nil {
			return bolt.ErrBucketNotFound
		}

		return b.Delete([]byte(key))
	})
	return err
}

// CheckConnection returns true if MongoDB is available
func CheckConnection() bool {
	if Mongo == nil {
		Connect(databases)
	}

	if Mongo != nil {
		return true
	}

	return false
}

// ReadConfig returns the database information
func ReadConfig() Info {
	return databases
}

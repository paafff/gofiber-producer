// package database

// import (
// 	"fmt"
// 	"log"
// 	"time"

// 	"gofiber-producer/internal/config"
// 	"gofiber-producer/internal/domain/models"

// 	"gorm.io/driver/postgres"
// 	"gorm.io/gorm"
// )

// var DB *gorm.DB

// func InitDatabase() {
// 	// Load database configuration from config.AppConfig
// 	dbConfig := config.AppConfig.Database

// 	fmt.Println(dbConfig)

// 	// Create DSN
// 	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
// 		dbConfig.Host, dbConfig.Username, dbConfig.Password, dbConfig.Name, dbConfig.Port, dbConfig.SSLMode, dbConfig.TimeZone)

// 	var err error

// 	// Logika retry untuk koneksi database
// 	for i := 0; i < 10; i++ {
// 		DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
// 		if err == nil {
// 			break
// 		}
// 		log.Printf("failed to connect to database (attempt %d): %v", i+1, err)
// 		time.Sleep(2 * time.Second)
// 	}

// 	if err != nil {
// 		log.Fatal("could not connect to the database: ", err)
// 	}

// 	// AutoMigrate will create the table based on the User struct
// 	err = DB.AutoMigrate(&models.User{})
// 	if err != nil {
// 		log.Fatal("failed to migrate database: ", err)
// 	}

// 	log.Println("Database migrated successfully")
// }

package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"gofiber-producer/internal/config"
	"gofiber-producer/internal/domain/models"

	"github.com/go-redis/redis/v8"
	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DB       *gorm.DB
	MongoDB  *mongo.Client
	RedisDB  *redis.Client
	RabbitMQ *amqp.Connection
)

func InitDatabase() {
	// PostgreSQL Initialization
	initPostgres()

	// MongoDB Initialization
	initMongoDB()

	// Redis Initialization
	initRedis()

	// RabbitMQ Initialization
	initRabbitMQ()
}

// func createDatabase(dbConfig config.AppConfig.PostgreSQL) error {
// 	dsn := fmt.Sprintf("host=%s user=%s password=%s port=%d sslmode=%s TimeZone=%s",
// 			dbConfig.Host, dbConfig.Username, dbConfig.Password, dbConfig.Port, dbConfig.SSLMode, dbConfig.TimeZone)

// 	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
// 	if err != nil {
// 			return err
// 	}

// 	createDBQuery := fmt.Sprintf("CREATE DATABASE %s", dbConfig.Name)
// 	err = db.Exec(createDBQuery).Error
// 	if err != nil {
// 			// Ignore error if database already exists
// 			if !gorm.ErrRecordNotFound.Is(err) {
// 					return err
// 			}
// 	}

// 	sqlDB, err := db.DB()
// 	if err != nil {
// 			return err
// 	}
// 	defer sqlDB.Close()

// 	return nil
// }

func initPostgres() {
	dbConfig := config.AppConfig.PostgreSQL

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
		dbConfig.Host, dbConfig.Username, dbConfig.Password, dbConfig.Name, dbConfig.Port, dbConfig.SSLMode, dbConfig.TimeZone)

	var err error
	for i := 0; i < 10; i++ {
		DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			break
		}
		log.Printf("failed to connect to database (attempt %d): %v", i+1, err)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		log.Fatal("could not connect to the database: ", err)
	}

	err = DB.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatal("failed to migrate database: ", err)
	}

	log.Println("Database migrated successfully")
}

func initMongoDB() {
	mongoConfig := config.AppConfig.MongoDB
	clientOptions := options.Client().ApplyURI(mongoConfig.URI)

	var err error
	MongoDB, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal("could not connect to MongoDB: ", err)
	}

	err = MongoDB.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal("could not ping MongoDB: ", err)
	}

	log.Println("Connected to MongoDB successfully")
}

func SaveLoginActivity(activity models.Activity) error {
	collection := MongoDB.Database("db_gofiber").Collection("activities")
	_, err := collection.InsertOne(context.Background(), activity)
	return err
}

// func InitMongoDB(uri string) error {
// 	clientOptions := options.Client().ApplyURI(uri)
// 	client, err := mongo.Connect(context.Background(), clientOptions)
// 	if err != nil {
// 			return err
// 	}

// 	err = client.Ping(context.Background(), nil)
// 	if err != nil {
// 			return err
// 	}

// 	MongoDB = client
// 	return nil
// }

func initRedis() {
	redisConfig := config.AppConfig.Redis
	RedisDB = redis.NewClient(&redis.Options{
		Addr:     redisConfig.Addr,
		Password: redisConfig.Password,
		DB:       redisConfig.DB,
	})

	_, err := RedisDB.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal("could not connect to Redis: ", err)
	}

	log.Println("Connected to Redis successfully")
}

var Conn *amqp.Connection
var Channel *amqp.Channel

func initRabbitMQ() {
	rabbitConfig := config.AppConfig.RabbitMQ
	// var err error

	// RabbitMQ, err = amqp.Dial(rabbitConfig.URI)
	// if err != nil {
	// 	log.Fatal("could not connect to RabbitMQ: ", err)
	// }

	var err error
	Conn, err = amqp.Dial(rabbitConfig.URI)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}

	Channel, err = Conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}

	_, err = Channel.QueueDeclare(
		"task_queue", // name
		true,         // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}

	log.Println("Connected to RabbitMQ successfully")
}

// MongoDB Example
func SaveToMongoDB(collection string, document interface{}) error {
	coll := MongoDB.Database("exampleDB").Collection(collection)
	_, err := coll.InsertOne(context.TODO(), document)
	return err
}

func GetFromMongoDB(collection string, filter interface{}) (*mongo.SingleResult, error) {
	coll := MongoDB.Database("exampleDB").Collection(collection)
	result := coll.FindOne(context.TODO(), filter)
	return result, result.Err()
}

// Redis Example
func SaveSessionToRedis(sessionID string, data interface{}, expiration time.Duration) error {
	return RedisDB.Set(context.Background(), sessionID, data, expiration).Err()
}

func GetSessionFromRedis(sessionID string) (string, error) {
	return RedisDB.Get(context.Background(), sessionID).Result()
}

// RabbitMQ Pub
func PublishMessage(body string) error {
	err := Channel.Publish(
		"",           // exchange
		"task_queue", // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	return err
}

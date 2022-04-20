package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/go-sql-driver/mysql"
	"github.com/segmentio/kafka-go"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var goRedisClient = redis.NewClient(&redis.Options{
	Addr:     "35.193.240.14:6379",
	Password: "",
	DB:       0,
})

type Game struct {
	Game_ID   int    `json:"game_id"`
	Players   int    `json:"players"`
	Game_Name string `json:"game_name"`
	Winner    int    `json:"winner"`
	Date      string `json:"date"`
	Queue     string `json:"queue"`
}
type Gamert struct {
	Game_ID   int    `json:"game_id"`
	Players   int    `json:"players"`
	Game_Name string `json:"game_name"`
	Winner    int    `json:"winner"`
}

const (
	topic         = "my-topic"
	brokerAddress = "kafka"
)

func consume(ctx context.Context) {
	l := log.New(os.Stdout, "Kafka escuchando: ", 0)
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{brokerAddress},
		Topic:   topic,
		GroupID: "my-group",
		Logger:  l,
	})

	for {
		log.Printf("---------------------------------------------------")
		msg, err := r.ReadMessage(ctx)
		log.Printf("---------------------------------------------------")
		if err != nil {
			panic("No se pudo leer el mensaje: " + err.Error())
		}
		var mensaje string
		mensaje = string(msg.Value)
		fmt.Println("Mensaje recibido: ", string(msg.Value))
		go insertRedis(string(msg.Value))
		go insertMongo(ctx, string(msg.Value))
		go insertTibd(ctx, mensaje)
	}
}

func insertRedis(body string) {
	var ctx = context.Background()

	val, err := goRedisClient.Get(ctx, "games").Result()
	if err == redis.Nil {
		//si no  hay datos los crea
		erro := goRedisClient.Set(ctx, "games", body, 10*time.Hour).Err()
		if erro != nil {
			panic(erro)
		}
	} else if err != nil {
		panic(err)
	} else {
		// si ya hay datos se concatena
		erro := goRedisClient.Set(ctx, "games", val+","+body, 10*time.Hour).Err()
		if erro != nil {
			panic(erro)
		}
	}
}

func insertMongo(ctx context.Context, body string) {
	var game Game
	t := time.Now()
	fecha := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())
	game.Date = fecha
	game.Queue = "Kafka"
	json.Unmarshal([]byte(body), &game)
	var collection = GetCollection("datas")
	var err error
	_, err = collection.InsertOne(ctx, game)

	if err != nil {
		fmt.Println(err.Error())
	}
}

func insertTibd(ctx context.Context, body string) {
	var game Gamert
	cfg := mysql.Config{
		User:                 "root",
		Passwd:               "",
		Net:                  "tcp",
		Addr:                 "34.72.26.121:4000",
		DBName:               "F2tidb",
		AllowNativePasswords: true,
	}

	json.Unmarshal([]byte(body), &game)

	fmt.Println("Go conectandose a Mysql...")

	// Open up our database connection.
	// I've set up a database on my local machine using phpmyadmin.
	// The database is called testDb
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		panic(err.Error())
	}
	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}

	fmt.Println("Connected!")
	defer db.Close()
	insertidb, err := db.Query("INSERT INTO Game(game_id,players,game_name,winner) VALUES (" + strconv.Itoa(game.Game_ID) + "," + strconv.Itoa(game.Players) + ",'" + game.Game_Name + "'," + strconv.Itoa(game.Winner) + ")")
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("Juego en TIDB")
	defer insertidb.Close()
}

var MONGO = "mongodb+srv://AdminUser:123@cluster0.eelrn.mongodb.net/proyectof2?retryWrites=true&w=majority"

func GetCollection(collection string) *mongo.Collection {
	client, err := mongo.NewClient(options.Client().ApplyURI(MONGO))
	if err != nil {
		panic(err.Error())
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		panic(err.Error())
	}

	return client.Database("proyectof2").Collection(collection)
}

func main() {
	ctx := context.Background()
	consume(ctx)
}

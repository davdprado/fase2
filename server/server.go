package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"math/rand"
	"net"
	"strconv"
	"time"

	"github.com/segmentio/kafka-go"

	pb "fase.com/grpc/gen/proto"

	"google.golang.org/grpc"
)

const (
	topic         = "my-topic"
	brokerAddress = "localhost:9092"
)

func produce(ctx context.Context, req *pb.Game) {

	winner := int(juegos(req.GameId, req.Players))
	name := names(req.GameId)

	conn, err := kafka.DialLeader(context.Background(), "tcp", brokerAddress, topic, 0)
	conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	_, err = conn.WriteMessages(
		kafka.Message{
			Key: []byte(strconv.Itoa(int(req.GameId))),
			Value: []byte("{\"game_id\":" + strconv.Itoa(int(req.GameId)) +
				", \"Players\": " +
				strconv.Itoa(int(req.Players)) +
				", \"game_name\": \"" + name +
				"\", \"winner\":" +
				strconv.Itoa(winner) +
				"}"),
			//
		},
	)

	if err != nil {
		fmt.Println("failed to write messages:")
		log.Fatal("failed to write messages:", err)
	}

	if err := conn.Close(); err != nil {
		fmt.Println("failed to close writer:")
		log.Fatal("failed to close writer:", err)
	}
	log.Printf("write successfully")
	fmt.Println("write successfully")
}

type TestApiServer struct {
	pb.UnimplementedTestApiServer
}

func (s *TestApiServer) Sendgame(ctx context.Context, req *pb.Game) (*pb.Response, error) {
	ganador := juegos(req.GameId, req.Players)
	log.Println(req.GameId, req.Players, ganador)
	msg := pb.Response{Msg: "Creado"}
	produce(ctx, req)
	return &msg, nil
}

func juegos(id int32, jugadores int32) int32 {
	if id == 1 {
		return 1
	} else if id == 2 {
		return jugadores
	} else if id == 3 {
		return jugadores / 2
	} else if id == 4 {
		return rand.Int31n(jugadores)
	} else if id == 5 {
		valor := int32(math.Sqrt(float64(jugadores)))
		return valor
	}
	return 0
}

func names(id int32) string {
	if id == 1 {
		return "First"
	} else if id == 2 {
		return "Last"
	} else if id == 3 {
		return "Middle"
	} else if id == 4 {
		return "Random"
	} else if id == 5 {
		return "Sqrt"
	}
	return ""
}

func main() {
	listner, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		log.Fatalln(err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterTestApiServer(grpcServer, &TestApiServer{})
	log.Println("Server on port 8080")
	err = grpcServer.Serve(listner)
	if err != nil {
		log.Fatal(err)
	}

}

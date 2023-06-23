package main

import (
	"context"
	"encoding/json"
	"log"
	"net"
	"strconv"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/szewczukk/score-service/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type UserScoreModel struct {
	UserId int
	Score  int
}

type ScoreServiceServer struct {
	proto.UnimplementedScoreServiceServer
	Db *gorm.DB
}

type CompleteResolutionPayload struct {
	UserId       int32 `json:"userId"`
	ResolutionId int32 `json:"resolutionId"`
}

func NewScoreServiceServer(db *gorm.DB) *ScoreServiceServer {
	return &ScoreServiceServer{
		Db: db,
	}
}

func main() {
	db, err := gorm.Open(sqlite.Open("score.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&UserScoreModel{})

	listener, err := net.Listen("tcp", ":3003")
	if err != nil {
		panic(err)
	}

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	defer ch.Close()

	queue, err := ch.QueueDeclare(
		"resolutionCompleted",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		panic(err)
	}

	userCreatedQueue, err := ch.QueueDeclare(
		"userCreated",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		panic(err)
	}

	msgs, err := ch.Consume(
		queue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		panic(err)
	}

	go func() {
		var forever chan struct{}
		for msg := range msgs {
			payload := new(CompleteResolutionPayload)
			err := json.Unmarshal(msg.Body, payload)
			if err != nil {
				panic(err)
			}

			userScore := UserScoreModel{}
			err = db.First(&userScore, &UserScoreModel{
				UserId: int(payload.UserId),
			}).Error

			if err != nil {
				panic(err)
			}

			db.Model(
				&UserScoreModel{},
			).Where(
				"user_id = ?", userScore.UserId,
			).Update(
				"score", userScore.Score+1,
			)
		}
		<-forever
	}()

	userCreatedMsgs, err := ch.Consume(
		userCreatedQueue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		panic(err)
	}

	go func() {
		var forever chan struct{}
		for msg := range userCreatedMsgs {
			userIdStr := string(msg.Body)
			userIdInt64, err := strconv.ParseInt(userIdStr, 10, 32)
			if err != nil {
				panic(err)
			}

			db.Create(&UserScoreModel{
				UserId: int(userIdInt64),
				Score:  0,
			})
		}
		<-forever
	}()

	grpcServer := grpc.NewServer()
	scoreServiceServer := NewScoreServiceServer(db)
	proto.RegisterScoreServiceServer(grpcServer, scoreServiceServer)
	reflection.Register(grpcServer)

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (s *ScoreServiceServer) GetAllUserScores(
	ctx context.Context,
	request *proto.GetAllUserScoresRequest,
) (*proto.GetAllUserScoresResponse, error) {
	var userScoreModels []UserScoreModel
	s.Db.Find(&userScoreModels)

	var protoUserScores []*proto.UserScore

	for _, userScore := range userScoreModels {
		protoUserScores = append(protoUserScores, &proto.UserScore{
			UserId: int32(userScore.UserId),
			Score:  int32(userScore.Score),
		})
	}

	return &proto.GetAllUserScoresResponse{Scores: protoUserScores}, nil
}

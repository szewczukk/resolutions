package main

import (
	"context"
	"log"
	"net"

	"github.com/szewczukk/score-service/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type UserScoreModel struct {
	UserId int `gorm:"primarykey"`
	Score  int
}

type ScoreServiceServer struct {
	proto.UnimplementedScoreServiceServer
	Db *gorm.DB
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

	listener, err := net.Listen("tcp", ":3001")
	if err != nil {
		panic(err)
	}

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

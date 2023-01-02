package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	db "go-grpc-postgresql-crud/db"
	"go-grpc-postgresql-crud/model"
	pb "go-grpc-postgresql-crud/proto"
	"log"
	"net"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

func init() {
	DB = db.DatabaseConnection()
}

var (
	err  error
	port = flag.Int("port", 50051, "gRPC server port")
	DB   *gorm.DB
)

type server struct {
	pb.UnimplementedMovieServiceServer
}

func main() {
	log.Println("database connected successfully ...")

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	log.Println("gRPC server running ...")

	s := grpc.NewServer()

	pb.RegisterMovieServiceServer(s, &server{})
	log.Printf("Server listening at %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve : %v", err)
	}
}

func (*server) CreateMovie(ctx context.Context, req *pb.CreateMovieRequest) (*pb.CreateMovieResponse, error) {
	log.Println("Create Movie ...")
	movie := req.GetMovie()
	movie.Id = uuid.New().String()

	data := model.Movie{
		ID:          movie.GetId(),
		Title:       movie.GetTitle(),
		Description: movie.GetDescription(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	res := DB.Create(&data)
	if res.RowsAffected == 0 {
		return nil, errors.New("Failed to create new movie ...")
	}
	return &pb.CreateMovieResponse{
		Movie: &pb.Movie{
			Id:          movie.GetId(),
			Title:       movie.GetTitle(),
			Description: movie.GetDescription(),
		},
	}, nil
}

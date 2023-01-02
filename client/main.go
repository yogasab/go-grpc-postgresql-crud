package main

import (
	"flag"
	"go-grpc-postgresql-crud/model"
	"log"
	"net/http"

	pb "go-grpc-postgresql-crud/proto"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	addr = flag.String("addr", "localhost:50051", "the address GRPC to connect to")
	err  error
)

func main() {
	flag.Parse()
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()
	client := pb.NewMovieServiceClient(conn)
	r := gin.Default()
	r.POST("/movies", func(ctx *gin.Context) {
		var (
			movie model.Movie
		)

		err := ctx.ShouldBind(&movie)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}

		data := &pb.Movie{
			Title:       movie.Title,
			Description: movie.Description,
		}
		res, err := client.CreateMovie(ctx, &pb.CreateMovieRequest{
			Movie: data,
		})
		if err != nil {
			log.Println(err)
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}
		ctx.JSON(http.StatusCreated, gin.H{
			"movie": res.Movie,
		})
	})
	r.Run(":5000")
}

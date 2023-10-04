package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	quizpb "github.com/wada1355/quiz_grpc/pkg/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type QuizServer struct {
	quizpb.UnimplementedQuizServiceServer
}

func NewQuizServer() *QuizServer {
	return &QuizServer{}
}

func (s *QuizServer) Hello(ctx context.Context, req *quizpb.HelloRequest) (*quizpb.HelloResponse, error) {
	return &quizpb.HelloResponse{
		Message: fmt.Sprintf("Hello, %s!", req.GetName()),
	}, nil
}

func main() {
	port := 8080
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		panic(err)
	}

	s := grpc.NewServer()

	quizpb.RegisterQuizServiceServer(s, NewQuizServer())

	reflection.Register(s)

	go func() {
		log.Printf("start gRPC server port: %v", port)
		s.Serve(listener)
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("stopping gRPC server...")
	s.GracefulStop()
}

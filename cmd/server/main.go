package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"github.com/wada1355/quiz_grpc/cmd/server/quiz"
	quizpb "github.com/wada1355/quiz_grpc/pkg/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	port := 8080
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		panic(err)
	}

	s := grpc.NewServer()

	quizpb.RegisterQuizServiceServer(s, quiz.NewQuizServer())

	reflection.Register(s)

	go func() {
		log.Printf("start gRPC server port: %v", port)
		s.Serve(listener)
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Printf("stop gRPC server port: %v", port)
	s.GracefulStop()
}

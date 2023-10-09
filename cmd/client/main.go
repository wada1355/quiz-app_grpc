package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/wada1355/quiz_grpc/cmd/quiz"
	quizpb "github.com/wada1355/quiz_grpc/pkg/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	address := "localhost:8080"
	conn, err := grpc.Dial(
		address,

		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		log.Fatal("Failed to connect gRPC")
		return
	}
	defer conn.Close()

	client := quizpb.NewQuizServiceClient(conn)
	scanner := bufio.NewScanner(os.Stdin)
	q := quiz.NewQuizService(client, scanner)
	if err := q.Quiz(); err != nil {
		fmt.Print(err)
	}
}

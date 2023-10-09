package main

import (
	// (一部抜粋)
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

	quizpb "github.com/wada1355/quiz_grpc/pkg/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	scanner *bufio.Scanner
	client  quizpb.QuizServiceClient
)

func main() {
	scanner = bufio.NewScanner(os.Stdin)

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

	client = quizpb.NewQuizServiceClient(conn)

	if err := Quiz(); err != nil {
		fmt.Print(err)
	}
}

func Quiz() error {
	stream, err := client.Quiz(context.Background())
	if err != nil {
		return err
	}
	questionNum, err := decideQuestionsNum(stream)
	sendNum, err := strconv.Atoi(questionNum)
	if err != nil {
		return err
	}
	sendCount := 0
	var sendEnd, recvEnd bool

	for !(sendEnd && recvEnd) {
		if !recvEnd {
			if res, err := stream.Recv(); err != nil {
				if errors.Is(err, io.EOF) {
					recvEnd = true
				} else {
					return err
				}
			} else {
				fmt.Println(res.GetMessage())
			}
		}
		if !sendEnd {
			scanner.Scan()
			myAnswer := scanner.Text()
			sendCount++
			if err := stream.Send(&quizpb.QuizRequest{
				Answer: myAnswer,
			}); err != nil {
				return err
			}
			if sendCount == sendNum {
				sendEnd = true
				if err := stream.CloseSend(); err != nil {
					return err
				}
			}
		}

	}
	return nil
}

func decideQuestionsNum(stream quizpb.QuizService_QuizClient) (string, error) {
	if err := receivePrint(stream); err != nil {
		return "", err
	}

	scanner.Scan()
	questionNum := scanner.Text()
	if err := stream.Send(&quizpb.QuizRequest{
		Answer: questionNum,
	}); err != nil {
		return "", err
	}

	if err := receivePrint(stream); err != nil {
		return "", err
	}

	return questionNum, nil
}

func receivePrint(stream quizpb.QuizService_QuizClient) error {
	msg, err := stream.Recv()
	if err != nil {
		return err
	}
	fmt.Println(msg.GetMessage())
	return nil
}

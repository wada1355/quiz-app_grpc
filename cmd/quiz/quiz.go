package quiz

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"

	quizpb "github.com/wada1355/quiz_grpc/pkg/grpc"
)

type QuizService struct {
	client  quizpb.QuizServiceClient
	scanner *bufio.Scanner
}

func NewQuizService(client quizpb.QuizServiceClient, scanner *bufio.Scanner) *QuizService {
	return &QuizService{
		client:  client,
		scanner: scanner,
	}
}

func (s *QuizService) Quiz() error {
	scanner := bufio.NewScanner(os.Stdin)

	stream, err := s.client.PlayQuiz(context.Background())
	if err != nil {
		return err
	}
	questionNum, err := s.decideQuestionsNum(stream)
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
			if err := stream.Send(&quizpb.QuizReq{
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

func (s *QuizService) decideQuestionsNum(stream quizpb.QuizService_PlayQuizClient) (string, error) {
	if err := receivePrint(stream); err != nil {
		return "", err
	}

	s.scanner.Scan()
	questionNum := s.scanner.Text()
	if err := stream.Send(&quizpb.QuizReq{
		Answer: questionNum,
	}); err != nil {
		return "", err
	}

	if err := receivePrint(stream); err != nil {
		return "", err
	}

	return questionNum, nil
}

func receivePrint(stream quizpb.QuizService_PlayQuizClient) error {
	msg, err := stream.Recv()
	if err != nil {
		return err
	}
	fmt.Println(msg.GetMessage())
	return nil
}

package quiz

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
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

func (s *QuizService) PlayQuiz() error {
	stream, err := s.client.PlayQuiz(context.Background())
	if err != nil {
		return err
	}
	questionNum, err := s.decideQuestionNum(stream)
	sendCount := 0
	var sendEnd, recvEnd bool

	if err := receivePrint(stream); err != nil {
		return err
	}

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
			s.scanner.Scan()
			myAnswer := s.scanner.Text()
			sendCount++
			if err := stream.Send(&quizpb.QuizReq{
				Answer: myAnswer,
			}); err != nil {
				return err
			}
			if sendCount == questionNum {
				sendEnd = true
				if err := stream.CloseSend(); err != nil {
					return err
				}
			}
		}

	}
	return nil
}

func (s *QuizService) decideQuestionNum(stream quizpb.QuizService_PlayQuizClient) (int, error) {
	if err := receivePrint(stream); err != nil {
		return 0, err
	}

	s.scanner.Scan()
	inputNum := s.scanner.Text()
	if err := stream.Send(&quizpb.QuizReq{
		Answer: inputNum,
	}); err != nil {
		return 0, err
	}

	questionNum, err := strconv.Atoi(inputNum)
	if err != nil {
		return 0, err
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

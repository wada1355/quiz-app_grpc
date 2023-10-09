package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"

	"github.com/wada1355/quiz_grpc/cmd/server/quizset"
	quizpb "github.com/wada1355/quiz_grpc/pkg/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type QuizServer struct {
	quizpb.UnimplementedQuizServiceServer
}

type QuizSession struct {
	quizSet     []quizset.QuizSet
	userAnswers []string
	stream      quizpb.QuizService_QuizServer
}

func NewQuizServer() *QuizServer {
	return &QuizServer{}
}

func (s *QuizServer) Hello(ctx context.Context, req *quizpb.HelloRequest) (*quizpb.HelloResponse, error) {
	return &quizpb.HelloResponse{
		Message: fmt.Sprintf("Hello, %s!", req.GetName()),
	}, nil
}

func (s *QuizServer) Quiz(stream quizpb.QuizService_QuizServer) error {
	session := &QuizSession{
		stream: stream,
	}

	quizSet, err := session.prepareQuizSet()
	if err != nil {
		return err
	}
	userAnswers := make([]string, 0)
	for i, q := range quizSet {
		if err := stream.Send(&quizpb.QuizResponse{
			Message: fmt.Sprintf("問題%d: %s", i+1, q.Question),
		}); err != nil {
			return err
		}
		req, err := stream.Recv()
		if err != nil {
			return err
		}
		userAnswers = append(userAnswers, req.GetAnswer())
	}
	session.quizSet = quizSet
	session.userAnswers = userAnswers
	if err := session.sendResult(); err != nil {
		return err
	}
	return nil
}

func (s *QuizSession) prepareQuizSet() ([]quizset.QuizSet, error) {
	if err := s.stream.Send(&quizpb.QuizResponse{
		Message: fmt.Sprintf("何問出題しますか？"),
	}); err != nil {
		return nil, err
	}
	req, err := s.stream.Recv()
	if err != nil {
		return nil, err
	}
	quizNumStr := req.GetAnswer()
	quizNum, err := strconv.Atoi(quizNumStr)
	if err != nil {
		return nil, err
	}
	if err := s.stream.Send(&quizpb.QuizResponse{
		Message: fmt.Sprintf("\n------%d問出題します-----\n", quizNum),
	}); err != nil {
		return nil, err
	}
	return quizset.GetRandomQuizSet(quizNum), nil
}

func (s *QuizSession) sendResult() error {
	if err := s.stream.Send(&quizpb.QuizResponse{
		Message: fmt.Sprint("\n------結果発表します-----\n"),
	}); err != nil {
		return err
	}
	correctNum := 0
	for i, q := range s.quizSet {
		userAnswer := s.userAnswers[i]
		if q.Answer == userAnswer {
			if err := s.stream.Send(&quizpb.QuizResponse{
				Message: fmt.Sprintf("問題%d: %s✅", i+1, userAnswer),
			}); err != nil {
				return err
			}
			correctNum++
		} else {
			if err := s.stream.Send(&quizpb.QuizResponse{
				Message: fmt.Sprintf("問題%d: %s❌ → 正解は%s", i+1, userAnswer, q.Answer),
			}); err != nil {
				return err
			}
		}
	}
	accuracy := int(float64(correctNum) / float64(len(s.quizSet)) * 100.0)
	if err := s.stream.Send(&quizpb.QuizResponse{
		Message: fmt.Sprintf("問題数: %d, 正解数: %d, 正解数: %d％", len(s.quizSet), correctNum, accuracy),
	}); err != nil {
		return err
	}
	return nil
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

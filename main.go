package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"sync"

	pb "github.com/dsrvlabs/vatz-plugin-terra/plugin"
	"github.com/dsrvlabs/vatz-plugin-terra/policy"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
)

var (
	executor policy.Executor
)

func init() {
	executor = policy.NewExecutor()
}

type pluginServer struct {
	pb.UnimplementedManagerPluginServer
}

func (s *pluginServer) Init(context.Context, *emptypb.Empty) (*pb.PluginInfo, error) {
	// TODO: TBD
	return nil, nil
}

func (s *pluginServer) Verify(context.Context, *emptypb.Empty) (*pb.VerifyInfo, error) {
	// TODO: TBD
	return nil, nil
}

func (s *pluginServer) Execute(ctx context.Context, req *pb.ExecuteRequest) (*pb.ExecuteResponse, error) {
	log.Println("pluginServer.Execute")

	resp := &pb.ExecuteResponse{
		State:   pb.ExecuteResponse_SUCCESS,
		Message: "OK",
	}

	fmt.Printf("ExecuteInfo %+v\n", req.ExecuteInfo)
	fmt.Printf("Fields %+v\n", req.ExecuteInfo.Fields)

	val, ok := req.ExecuteInfo.Fields["function"]
	if !ok {
		resp.State = pb.ExecuteResponse_FAILURE
		resp.Message = "no valid function"
		return resp, nil
	}

	funcName := val.GetStringValue()

	fmt.Println("Function is ", funcName)
	if funcName == "" {
		resp.State = pb.ExecuteResponse_FAILURE
		resp.Message = "no valid function"
		return resp, nil
	}

	switch funcName {
	case "UpTerrad":
		isUp, err := executor.UpTerra()
		if err != nil {
			return nil, err
		}

		if !isUp {
			resp.Message = "dead"
		}
		//	case "UpOracle":
		//		isUp, err := executor.IsHeimdallUp()
		//		if err != nil {
		//			return nil, err
		//		}
		//
		//		if !isUp {
		//			resp.Message = "dead"
		//		}
	default:
		log.Println("No selection")
		resp.Message = "No function"
	}

	return resp, nil
}

func main() {
	ch := make(chan os.Signal, 1)
	startServer(ch)
}

func startServer(ch <-chan os.Signal) {
	log.Println("Start vatz-plugin-terra")

	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	c, err := net.Listen("tcp", "0.0.0.0:9091")
	if err != nil {
		log.Println(err)
	}

	s := grpc.NewServer()

	serv := pluginServer{}
	pb.RegisterManagerPluginServer(s, &serv)

	reflection.Register(s)

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		_ = <-ch
		cancel()
		s.GracefulStop()
		wg.Done()
	}()

	if err := s.Serve(c); err != nil {
		log.Panic(err)
	}

	wg.Wait()
}

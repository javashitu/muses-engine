package micro

import (
	"fmt"
	"log"
	"muses-engine/config"
	"muses-engine/internal/facade"
	"muses-engine/pkg/consul"
	"muses-engine/pkg/proto_generated/pb"
	"net"
	"time"

	"google.golang.org/grpc"
)

func RegistConsul() {
	go func() {
		regestryConsul()
	}()
}

func regestryConsul() {
	log.Println("begin connect to consul")
	consulAddr := config.MicroConfig.Consul.GenConsulAddr()
	log.Println("this consulAddr is ", consulAddr)
	//1. 启动consul，注入consul服务地址
	consulService, _ := consul.NewService(consulAddr)
	log.Println("begin register service, the grpc register config ", config.MicroConfig.Grpc)

	//2. 注册服务
	consulErr := consulService.RegisterService(consul.RegisterService{
		ServiceName:   config.MicroConfig.Grpc.ServiceName,
		Address:       config.MicroConfig.Grpc.Address,
		Port:          config.MicroConfig.Grpc.Port,
		HeathCheckTTL: time.Duration(config.MicroConfig.Grpc.HealthCheckSeconds) * time.Second,
	})
	log.Println("register service success")

	if consulErr != nil {
		log.Fatalln(consulErr)
		return
	}

	//3. 启动grpc服务器
	grpcServer := grpc.NewServer()

	//4. 把hello和grpc服务器绑定
	pb.RegisterMyHelloServer(grpcServer, &facade.MyHello{})
	pb.RegisterFileStoreServiceServer(grpcServer, &facade.FileStoreService{})
	// registGrpcServer(grpcServer)

	//5. 启动端口监听
	listener, err := net.Listen("tcp", config.MicroConfig.Grpc.GenServiceAddr())
	if err != nil {
		fmt.Println("listener error:", err)
		return
	}
	defer listener.Close()

	log.Println("Service starting successfully")

	//6. grpc和端口监听绑定
	errGrpc := grpcServer.Serve(listener)
	if errGrpc != nil {
		log.Panic("grpc server cant't start error:", errGrpc)
		return
	}
	log.Println("regist finished")

}

func registGrpcServer(grpcServer *grpc.Server) {
	pb.RegisterMyHelloServer(grpcServer, &facade.MyHello{})
	pb.RegisterFileStoreServiceServer(grpcServer, &facade.FileStoreService{})

}

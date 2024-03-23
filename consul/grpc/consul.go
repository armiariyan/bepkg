package grpc

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"time"
)

type ConsulConnect struct {
	Target  string
	Timeout time.Duration
}

func ConnectWithRoundRobin(cfg *ConsulConnect) (conn *grpc.ClientConn, err error) {
	SetupResolver()

	ctx, _ := context.WithTimeout(context.Background(), cfg.Timeout)
	conn, err = grpc.DialContext(ctx, cfg.Target, grpc.WithInsecure(), grpc.WithBalancerName("round_robin"))
	return
}

func Close(conn *grpc.ClientConn) {
	err := conn.Close()
	if err != nil {
		fmt.Println(err)
	}
}

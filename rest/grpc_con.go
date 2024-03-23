package rest

import (
	"context"
	"time"

	"google.golang.org/grpc"

	Session "github.com/armiariyan/bepkg/session"
)

const sessionKey = "session_key"

type RpcConnection struct {
	options Options
	Conn    *grpc.ClientConn
}

func (rpc *RpcConnection) CreateContext(parent context.Context, session *Session.Session) (ctx context.Context) {
	ctx, _ = context.WithTimeout(parent, rpc.options.Timeout*time.Second)
	ctx = context.WithValue(ctx, sessionKey, session)
	return
}

func NewGRpcConnection(options Options) *RpcConnection {
	// todo still always insecure
	conn, err := grpc.Dial(options.Address, grpc.WithInsecure(), withClientUnaryInterceptor())
	if err != nil {
		panic(err)
	}

	return &RpcConnection{
		Conn:    conn,
		options: options,
	}
}

func NewGRpcConnectionE(options Options) (rpc *RpcConnection, err error) {
	// todo still always insecure
	conn, err := grpc.Dial(options.Address, grpc.WithInsecure(), withClientUnaryInterceptor())
	if err != nil {
		return
	}

	rpc = &RpcConnection{
		Conn:    conn,
		options: options,
	}
	return
}

func clientInterceptor(
	ctx context.Context,
	method string,
	req interface{},
	reply interface{},
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption,
) error {
	session := ctx.Value(sessionKey).(*Session.Session)
	processTime := session.T2("[request][", method, "] ---> ", req)
	err := invoker(ctx, method, req, reply, cc, opts...)
	session.T3(processTime, "[response][", method, "] ---> ", reply)
	return err
}

func withClientUnaryInterceptor() grpc.DialOption {
	return grpc.WithUnaryInterceptor(clientInterceptor)
}

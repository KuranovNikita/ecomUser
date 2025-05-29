package suite

import (
	"context"
	"ecomUser/internal/config"
	"fmt"
	"net"
	"strconv"
	"testing"

	user1 "github.com/KuranovNikita/ecomProto/gen/go/user"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Suite struct {
	*testing.T
	Cfg        *config.Config
	AuthClient user1.UserServiceClient
}

const (
	grpcHost = "localhost"
)

func New(t *testing.T) (context.Context, *Suite) {
	t.Helper()
	t.Parallel()

	godotenv.Load(".env")

	fmt.Println("start suite")
	cfg := config.MustLoad()
	fmt.Println("end config suite")

	grpcRequestTimeout := cfg.GRPCTimeout

	ctx, cancelCtx := context.WithTimeout(context.Background(), grpcRequestTimeout)

	t.Cleanup(func() {
		t.Helper()
		cancelCtx()
	})

	cc, err := grpc.NewClient(
		grpcAddress(cfg),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatalf("grpc server connection failed: %v", err)
	}

	t.Cleanup(func() {
		if cc != nil {
			cc.Close()
		}
	})

	return ctx, &Suite{
		T:          t,
		Cfg:        cfg,
		AuthClient: user1.NewUserServiceClient(cc),
	}

}

func grpcAddress(cfg *config.Config) string {
	return net.JoinHostPort(grpcHost, strconv.Itoa(cfg.GRPCPort))
}

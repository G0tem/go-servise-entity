package grpc

import (
	"context"
	"time"

	"github.com/G0tem/go-servise-entity/internal/config"
	"github.com/G0tem/go-servise-entity/proto/auth"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// AuthClient обертка для gRPC клиента авторизации
type AuthClient struct {
	conn   *grpc.ClientConn
	client auth.AuthServiceClient
}

// NewAuthClient создает новый gRPC клиент для авторизации
func NewAuthClient(cfg *config.Config) (*AuthClient, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(
		ctx,
		cfg.AuthGrpcAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, err
	}

	client := auth.NewAuthServiceClient(conn)

	log.Info().
		Str("address", cfg.AuthGrpcAddress).
		Msg("Connected to auth gRPC service")

	return &AuthClient{
		conn:   conn,
		client: client,
	}, nil
}

// GetTestData получает тестовые данные от auth сервиса
func (c *AuthClient) GetTestData(ctx context.Context, message string) (*auth.GetTestDataResponse, error) {
	req := &auth.GetTestDataRequest{
		Message: message,
	}

	return c.client.GetTestData(ctx, req)
}

// GetUserInfo получает информацию о пользователе от auth сервиса
func (c *AuthClient) GetUserInfo(ctx context.Context, userID string) (*auth.GetUserInfoResponse, error) {
	req := &auth.GetUserInfoRequest{
		UserId: userID,
	}

	return c.client.GetUserInfo(ctx, req)
}

// Close закрывает соединение с gRPC сервером
func (c *AuthClient) Close() error {
	return c.conn.Close()
}

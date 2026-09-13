package mongox

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// Config описывает подключение к MongoDB.
type Config struct {
	Host     string
	Port     string
	Database string
	User     string
	Password string
	// AuthSource используется только при наличии User/Password.
	AuthSource string
	// ObjectIDAsHexString включает декодирование ObjectID в hex-строки.
	// Сервисы, работающие с bson.ObjectID напрямую, должны оставить false.
	ObjectIDAsHexString bool
}

// Connect подключается к MongoDB и проверяет соединение, возвращая базу данных.
func Connect(ctx context.Context, cfg Config) (*mongo.Database, error) {
	clientOptions := options.Client().ApplyURI(buildURI(cfg))
	if cfg.ObjectIDAsHexString {
		clientOptions.SetBSONOptions(&options.BSONOptions{
			ObjectIDAsHexString: true,
		})
	}

	client, err := mongo.Connect(clientOptions)
	if err != nil {
		return nil, fmt.Errorf("не смог подключиться к MongoDB: %w", err)
	}

	if err = client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("не смог проверить подключение к MongoDB: %w", err)
	}

	return client.Database(cfg.Database), nil
}

func buildURI(cfg Config) string {
	port := cfg.Port
	if port == "" {
		port = "27017"
	}

	if cfg.User == "" && cfg.Password == "" {
		return fmt.Sprintf("mongodb://%s:%s/%s", cfg.Host, port, cfg.Database)
	}

	uri := fmt.Sprintf("mongodb://%s:%s@%s:%s/%s", cfg.User, cfg.Password, cfg.Host, port, cfg.Database)
	if cfg.AuthSource != "" {
		uri += "?authSource=" + cfg.AuthSource
	}
	return uri
}

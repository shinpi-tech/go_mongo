package mongox

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"strings"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// Config описывает подключение к MongoDB.
type Config struct {
	// URI — полная строка подключения. Если задана, поля Host/Port/User/Password
	// игнорируются (Database по-прежнему берётся из cfg.Database).
	URI      string
	Host     string
	Port     string
	Database string
	User     string
	Password string
	// AuthSource используется только при наличии User/Password.
	AuthSource string
	// DirectConnection включает directConnection=true (подключение к одному узлу
	// без обнаружения реплик).
	DirectConnection bool
	// ObjectIDAsHexString включает декодирование ObjectID в hex-строки.
	// Сервисы, работающие с bson.ObjectID напрямую, должны оставить false.
	ObjectIDAsHexString bool
}

// Connect подключается к MongoDB и проверяет соединение, возвращая базу данных.
func Connect(ctx context.Context, cfg Config) (*mongo.Database, error) {
	clientOptions := options.Client().ApplyURI(resolveURI(cfg))
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

// resolveURI выбирает источник строки подключения:
// Config.URI → переменная окружения MONGO_URL → сборка из Host/Port/User/Password.
func resolveURI(cfg Config) string {
	if cfg.URI != "" {
		return cfg.URI
	}
	if env := os.Getenv("MONGO_URL"); env != "" {
		return env
	}
	return buildURI(cfg)
}

func buildURI(cfg Config) string {
	port := cfg.Port
	if port == "" {
		port = "27017"
	}

	var b strings.Builder
	b.WriteString("mongodb://")
	if cfg.User != "" || cfg.Password != "" {
		// url.UserPassword корректно экранирует спецсимволы в логине и пароле.
		b.WriteString(url.UserPassword(cfg.User, cfg.Password).String())
		b.WriteString("@")
	}
	fmt.Fprintf(&b, "%s:%s/%s", cfg.Host, port, cfg.Database)

	var params []string
	if cfg.AuthSource != "" && (cfg.User != "" || cfg.Password != "") {
		params = append(params, "authSource="+cfg.AuthSource)
	}
	if cfg.DirectConnection {
		params = append(params, "directConnection=true")
	}
	if len(params) > 0 {
		b.WriteString("?" + strings.Join(params, "&"))
	}

	return b.String()
}

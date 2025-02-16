package postgres

import (
	"context"
	"net"
	"net/url"
	"os"

	"github.com/cockroachdb/errors"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	defaultUser     = "postgres"
	defaultPassword = "postgres"
	defaultHost     = "localhost"
	defaultPort     = "5432"
	defaultDBName   = "shop"
)

func Pool(ctx context.Context, connString string) (*pgxpool.Pool, error) {
	if connString == "" {
		u := url.URL{
			Scheme:   "postgres",
			User:     url.UserPassword(os.Getenv("DATABASE_USER"), os.Getenv("DATABASE_PASSWORD")),
			Host:     net.JoinHostPort(os.Getenv("DATABASE_HOST"), os.Getenv("DATABASE_PORT")),
			Path:     os.Getenv("DATABASE_NAME"),
			RawQuery: "sslmode=disable",
		}

		p, _ := u.User.Password()
		if u.User.Username() == "" || p == "" {
			u.User = url.UserPassword(defaultUser, defaultPassword)
		}

		if u.Hostname() == "" || u.Port() == "" {
			u.Host = net.JoinHostPort(defaultHost, defaultPort)
		}

		if u.Path == "" {
			u.Path = defaultDBName
		}

		connString = u.String()
	}

	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	err = pool.Ping(ctx)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return pool, nil
}

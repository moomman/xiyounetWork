package db

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	db "github.com/moomman/xiyounetWork/internal/dao/db/sqlc"
)

type DB struct {
	db.Store
}

func Init(dataSourceName string) *db.SqlStore {
	pool, err := pgxpool.Connect(context.Background(), dataSourceName)
	if err != nil {
		panic(err)
	}
	return &db.SqlStore{
		Queries: db.New(pool),
		DB: pool ,
	}
}



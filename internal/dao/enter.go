package dao

import "github.com/moomman/xiyounetWork/internal/dao/db"

type group struct {
	DB db.DB
}

var Group = new(group)
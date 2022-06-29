////////////////////////////////////////////////////////////////////////////////
// wraith - the wraith game engine and server
// Copyright (c) 2022 Michael D. Henderson
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published
// by the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.
////////////////////////////////////////////////////////////////////////////////

package models

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/mdhender/wraith/storage/config"
	"log"
	"time"
)

type Store struct {
	db          *sql.DB
	version     string
	dateFormat  string
	endOfTime   time.Time
	ctx         context.Context
	unitsByCode map[string]*Unit
	unitsById   map[int]*Unit
}

func Open(cfg *config.Global) (*Store, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(localhost:3306)/%s?multiStatements=true&parseTime=true", cfg.User, cfg.Password, cfg.Schema)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	} else if err := db.Ping(); err != nil {
		return nil, err
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	maxConns := 10
	db.SetMaxOpenConns(maxConns)
	db.SetMaxIdleConns(maxConns)

	return &Store{
		db:        db,
		version:   "0.1.0",
		endOfTime: time.Date(2099, 12, 31, 23, 59, 59, 0, time.UTC),
		ctx:       context.Background(),
	}, nil
}

func (s *Store) Close() {
	if s.db == nil {
		return
	}
	if err := s.db.Close(); err != nil {
		log.Printf("%+v\n", err)
	}
	s.db = nil
}

func (s *Store) Ping() error {
	return s.db.Ping()
}

func (s *Store) Version() string {
	return s.version
}

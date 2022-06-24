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

package engine

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/mdhender/wraith/storage/config"
	"log"
	"time"
)

// Open returns an initialized engine with the base configuration and the games store loaded.
func Open(cfg *config.Global, ctx context.Context) (e *Engine, err error) {
	if cfg == nil {
		return nil, errors.New("missing base config")
	}

	e = &Engine{ctx: ctx}
	e.config.base = cfg

	dsn := fmt.Sprintf("%s:%s@tcp(localhost:3306)/%s?multiStatements=true&parseTime=true", cfg.User, cfg.Password, cfg.Schema)
	e.db, err = sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	} else if err := e.db.Ping(); err != nil {
		return nil, err
	}
	e.db.SetConnMaxLifetime(time.Minute * 3)
	maxConns := 10
	e.db.SetMaxOpenConns(maxConns)
	e.db.SetMaxIdleConns(maxConns)

	return e, nil
}

func (e *Engine) createGame() error {
	return fmt.Errorf("engine.createGame: %w", ErrNotImplemented)
}

func (e *Engine) deleteGame(id int) error {
	_, err := e.db.Exec("delete from games where id = ?", id)
	if err != nil {
		return err
	}
	return nil
}

func (e *Engine) deleteGameByName(shortName string) error {
	_, err := e.db.Exec("delete from games where short_name = ?", shortName)
	if err != nil {
		return err
	}
	return nil
}

func (e *Engine) lookupGame(id int) *Game {
	var g Game
	row := e.db.QueryRow("select id, name, descr, current_turn from games where id = ?", id)
	err := row.Scan(&g.Id, &g.Name, &g.Descr, &g.Turn)
	if err != nil {
		return nil
	}
	return &g
}

func (e *Engine) lookupGameByName(shortName string) *Game {
	var g Game
	row := e.db.QueryRow("select id, name, descr, current_turn from games where short_name= ?", shortName)
	err := row.Scan(&g.Id, &g.Name, &g.Descr, &g.Turn)
	if err != nil {
		return nil
	}
	return &g
}

func (e *Engine) saveGame() error {
	// get a transaction with a deferred rollback in case things fail
	tx, err := e.db.BeginTx(e.ctx, nil)
	if err != nil {
		return fmt.Errorf("engine.saveGame: beginTx: %w", err)
	}
	defer tx.Rollback()

	r, err := tx.ExecContext(e.ctx, "insert into games (short_name, name, descr, current_turn) values (?, ?, ?, ?)",
		e.game.ShortName, e.game.Name, e.game.Descr, e.game.Turn)
	if err != nil {
		return fmt.Errorf("engine.saveGame: insert: 107: %w", err)
	}
	id, err := r.LastInsertId()
	if err != nil {
		return fmt.Errorf("engine.saveGame: fetchId: %w", err)
	}
	e.game.Id = int(id)
	log.Printf("created game %3d %s\n", int(id), e.game.ShortName)

	for _, turn := range e.game.Turns {
		_, err = tx.ExecContext(e.ctx, "insert into turns (game_id, turn, start_dt, end_dt) values (?, ?, ?, ?)",
			e.game.Id, turn.No, turn.EffDt, turn.EndDt)
		if err != nil {
			return fmt.Errorf("engine.saveGame: insert: 120: %w", err)
		}
	}

	for _, nation := range e.game.Nations {
		r, err := tx.ExecContext(e.ctx, "insert into nations (game_id, nation_no, speciality, descr) values (?, ?, ?, ?)",
			e.game.Id, nation.Id, nation.Speciality, nation.Description)
		if err != nil {
			return fmt.Errorf("engine.saveGame: insert: 128: %w", err)
		}
		id, err := r.LastInsertId()
		if err != nil {
			return fmt.Errorf("engine.saveGame: %w", err)
		}
		log.Printf("created nation %3d %8d\n", nation.Id, int(id))
	}

	return tx.Commit()
}

// Load retrieves a game from the store
func (e *Engine) Load(id string) error {
	if e == nil {
		return ErrNoEngine
	} else if e.db == nil {
		return ErrNoStore
	}
	e.reset()

	return fmt.Errorf("engine.Load: %w", ErrNotImplemented)
}

func (e *Engine) Save() error {
	if e == nil {
		return ErrNoEngine
	} else if e.db == nil {
		return ErrNoStore
	} else if e.game == nil {
		return ErrNoGame
	}

	return fmt.Errorf("engine.Save: %w", ErrNotImplemented)
}

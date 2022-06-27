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
	"github.com/google/uuid"
	"github.com/mdhender/wraith/storage/config"
	"log"
	"os"
)

// Bootstrap creates a new store.
func Bootstrap(cfg *config.Global) (*Store, error) {
	s, err := Open(cfg)
	if err != nil {
		return nil, err
	}

	// load and run the schema generation script
	b, err := os.ReadFile(cfg.SchemaFile)
	if err != nil {
		return nil, err
	}
	log.Printf("loaded schema file %q\n", cfg.SchemaFile)
	if _, err = s.db.Exec(string(b)); err != nil {
		return nil, err
	}
	log.Printf("executed schema file %q\n", cfg.SchemaFile)

	// create the default users required by the engine
	for _, user := range []string{"nobody", "sysop", "batch"} {
		err := s.CreateUser(user, user, user, uuid.New().String())
		if err != nil {
			return nil, err
		}
		log.Printf("created user %q\n", user)
	}

	// create the default set of units used by the engine
	type unit struct {
		code  string
		name  string
		descr string
	}
	for _, u := range []unit{
		{"ANM", "anti-missile", "anti-missile"},
		{"ASC", "assault-craft", "assault-craft"},
		{"ASW", "assault-weapon", "assault-weapon"},
		{"AUT", "automation", "automation"},
		{"CNGD", "consumer-goods", "consumer-goods"},
		{"ESH", "energy-shield", "energy-shield"},
		{"EWP", "energy-weapon", "energy-weapon"},
		{"FCT", "factory", "factory"},
		{"FOOD", "food", "food"},
		{"FRM", "farm", "farm"},
		{"FUEL", "fuel", "fuel"},
		{"GOLD", "gold", "gold"},
		{"HDR", "hyper-drive", "hyper-drive"},
		{"LSP", "life-support", "life-support"},
		{"LTSU", "light-structural", "light-structural"},
		{"MIN", "mine", "mine"},
		{"MLR", "military-robots", "military-robots"},
		{"MLSP", "military-supplies", "military-supplies"},
		{"MSS", "missile", "missile"},
		{"MSL", "missile-launcher", "missile-launcher"},
		{"MTLS", "metallics", "metallics"},
		{"NMTS", "non-metallics", "non-metallics"},
		{"SDR", "space-drive", "space-drive"},
		{"SNR", "sensor", "sensor"},
		{"SLSU", "super-light-structural", "super-light-structural"},
		{"STUN", "structural", "structural"},
		{"TPT", "transport", "transport"},
	} {
		if err := s.CreateUnit(u.code, u.name, u.descr, false); err != nil {
			return nil, err
		}
	}

	return s, nil
}

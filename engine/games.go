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

//// Games configuration
//type Games struct {
//	Store string       `json:"store"` // default path to store data
//	Index []GamesIndex `json:"index"`
//}
//
//type GamesIndex struct {
//	Id    string `json:"id"`    // unique identifier for game
//	Store string `json:"store"` // path to the game store file
//}
//
//// ReadGames loads a store from a JSON file.
//// It returns any errors.
//func (e *Engine) ReadGames() error {
//	b, err := ioutil.ReadFile(filepath.Join(e.stores.games.Store, "store.json"))
//	if err != nil {
//		return err
//	}
//	return json.Unmarshal(b, e.stores.games)
//}
//
//// WriteGames writes a store to a JSON file.
//// It returns any errors.
//func (e *Engine) WriteGames() error {
//	b, err := json.MarshalIndent(e.stores.games, "", "  ")
//	if err != nil {
//		return err
//	}
//	return ioutil.WriteFile(filepath.Join(e.stores.games.Store, "store.json"), b, 0600)
//}

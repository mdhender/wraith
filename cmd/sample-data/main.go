/*
 * Wraith Game Engine
 * Copyright (c) 2022 Michael D. Henderson
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published
 * by the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

package main

import (
	"github.com/mdhender/wraith/internal/seeder"
	//"github.com/mdhender/wraith/internal/storage/json"
	//"github.com/mdhender/wraith/internal/storage/words"
	"log"
	"math/rand"
	"time"
)

func main() {
	started := time.Now()

	// default log format to UTC
	log.SetFlags(log.Ldate | log.Ltime | log.LUTC)

	// use the crand package to seed the default PRNG source.
	if seed, err := seeder.Seed(); err != nil {
		log.Fatalln(err)
	} else {
		rand.Seed(seed)
	}

	if err := run(); err != nil {
		log.Println(err)
	}

	elapsed := time.Now().Sub(started)
	log.Printf("total time: %+v\n", elapsed)
}

func run() error {
	//psg, err := words.New(" ")
	//if err != nil {
	//	return err
	//}
	//s, err := json.NewAccountsDriver("accounts.json")
	//if err != nil {
	//	return err
	//}
	//v := s.Read(func(a json.Account) bool {
	//	return a.Handle == "sysop"
	//})
	//if len(v) == 0 {
	//	if err := s.Create(json.Account{
	//		Id:     "",
	//		Handle: "sysop",
	//		Salt:   []byte("salt"),
	//		Secret: []byte(psg.Generate(5)),
	//		Roles:  []string{"sysop"},
	//		Games:  nil,
	//	}); err != nil {
	//		log.Fatalln(err)
	//	}
	//} else {
	//	fmt.Printf("sysop: found %+v\n", v)
	//}
	//for _, a := range s.ReadAll() {
	//	fmt.Printf("readall: found %+v\n", a)
	//	fmt.Printf("readall: salt %q secret %q\n", string(a.Salt), string(a.Secret))
	//}
	//for _, a := range s.ReadAll() {
	//	if a.Handle != "sysop" {
	//		_ = s.DeleteById(a.Id)
	//	}
	//}
	return nil
}

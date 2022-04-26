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

package cmd

import (
	"fmt"
	"github.com/mdhender/wraith/internal/config"
	"github.com/spf13/cobra"
	"log"
	"math/rand"
)

var globalInit struct {
	ConfigFile string
	Secrets    struct {
		Signing string
		Sysop   string
	}
}

var cmdInit = &cobra.Command{
	Use:   "init",
	Short: "create a new global configuration file",
	Long:  `Create a new global configuration file.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		globalInit.ConfigFile = globalBase.ConfigFile
		if globalInit.ConfigFile == "" {
			return fmt.Errorf("missing config file name")
		}
		if len(globalInit.Secrets.Sysop) < 12 {
			return fmt.Errorf("sysop-password must be at least 12 characters long")
		}
		if len(globalInit.Secrets.Signing) == 0 {
			const bag = "0123456789_ABCDEFGHIJKLMNOPQRSTUVWXYZ-abcdefghijklmnopqrstuvwxyz.:,"
			bytes := make([]byte, len(bag))
			for i := 0; i < len(bag); i++ {
				bytes[i] = byte(bag[rand.Intn(len(bag))])
			}
			globalInit.Secrets.Signing = string(bytes)
		}
		if len(globalInit.Secrets.Signing) < 12 {
			return fmt.Errorf("signing-key must be at least 12 characters long")
		}
		cfg := config.Config{
			ConfigFile: globalInit.ConfigFile,
			Secrets: struct {
				Signing string `json:"signing-secret"`
				Sysop   string `json:"sysop-password"`
			}{
				Signing: globalInit.Secrets.Signing,
				Sysop:   globalInit.Secrets.Sysop,
			},
		}
		if err := config.Write(globalInit.ConfigFile, &cfg); err != nil {
			return err
		}
		log.Printf("[init] created %q\n", globalInit.ConfigFile)
		return nil
	},
}

func init() {
	cmdInit.Flags().StringVar(&globalInit.Secrets.Signing, "signing-key", "", "key for signing tokens")
	cmdInit.Flags().StringVar(&globalInit.Secrets.Sysop, "sysop-password", "", "password for sysop account")
	_ = cmdInit.MarkFlagRequired("sysop-password")

	cmdBase.AddCommand(cmdInit)
}

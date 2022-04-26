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
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
)

var rootGlobals struct {
	TestFlag    bool
	VerboseFlag bool
	ConfigFile  string // configuration file from command line flag

	envPrefix  string // value to prepend when converting flags to env variables
	cfgName    string // default configuration file name
	homeFolder string // derived path to home directory
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "wraith",
	Short: "Wraith game engine",
	Long: `wraith is the game engine for Wraith.
This application provides an API to the game engine.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// set the env and config since this hook always runs first
		rootGlobals.envPrefix, rootGlobals.cfgName = "WRAITH", ".wraith"
		// find home directory
		var err error
		if rootGlobals.homeFolder, err = homedir.Dir(); err != nil {
			return err
		}
		// now bind viper and cobra configuration
		return bindConfig(cmd)
	},
	Run: func(cmd *cobra.Command, args []string) {
		log.Printf("env: %-30s == %q\n", "HOME", rootGlobals.homeFolder)
		log.Printf("env: %-30s == %q\n", "WRAITH_CONFIG", viper.ConfigFileUsed())
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	rootCmd.PersistentFlags().StringVar(&rootGlobals.ConfigFile, "config", "", "Config file (default is $HOME/.wraith)")
	rootCmd.PersistentFlags().BoolVar(&rootGlobals.TestFlag, "test", false, "Test mode")
	rootCmd.PersistentFlags().BoolVar(&rootGlobals.VerboseFlag, "verbose", false, "Verbose mode")

	// Cobra also supports local flags, which will only run when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

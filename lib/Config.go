/*
 * Copyright (c) 2022 Aisha Tammy <aisha@bsd.ac>
 *
 * Permission to use, copy, modify, and distribute this software for any
 * purpose with or without fee is hereby granted, provided that the above
 * copyright notice and this permission notice appear in all copies.
 *
 * THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
 * WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
 * MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
 * ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
 * WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
 * ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
 * OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
 *
 */
package lib

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/ini.v1"
)

type ConfigData struct {
	Username    string `ini:"Username"`
	AccessToken string `ini:"AccessToken"`
	Homeserver  string `ini:"Homeserver"`
	DebugLevel  int    `ini:"DebugLevel"`
	Autojoin    bool   `ini:"Autojoin"`
	DBDirectory string `ini:"DBDirectory"`
	DBtype      string `ini:"DBtype"`
	DBdsn       string `ini:"DBdsn"`
}

func ReadConfig(ConfigFile string) (*ConfigData, error) {
	cfg := new(ConfigData)
	cfg.Username = ""
	cfg.AccessToken = ""
	cfg.Homeserver = ""
	cfg.DebugLevel = 0
	cfg.Autojoin = false
	cfg.DBDirectory = "/var/db/karma-bot"

	// valid SQL driver name: sqlite3, mysql, postgresql
	cfg.DBtype = "sqlite3"
	// Data Source Name (DSN) examples:
	// sqlite   - file:/var/db/karma-bot/sqlite3/data.sqlite3
	// postgres -
	//            postgres://username:password@localhost:5432/dbName
	//            postgres://username:password@%2Fvar%2Frun%2Fpostgresql/dbName
	//
	// mysql    -
	//            username:password@tcp(localhost:3306)/dbName
	//            username:password@unix(/tmp/mysql.sock)/dbName
	//
	cfg.DBdsn = ""

	err := ini.MapTo(cfg, ConfigFile)
	if err != nil {
		return nil, fmt.Errorf("Failed to read config file '%s': %v", ConfigFile, err)
	}
	if cfg.Username == "" {
		return nil, fmt.Errorf("Config file does not have 'Homeserver'")
	}
	if cfg.AccessToken == "" {
		return nil, fmt.Errorf("Config file does not have 'Username'")
	}
	if cfg.Homeserver == "" {
		return nil, fmt.Errorf("Config file does not have 'AccessToken'")
	}
	absDBDir, err := filepath.Abs(cfg.DBDirectory)
	if err != nil {
		return nil, fmt.Errorf("Could not get absolute path of DBDirectory (%s): %v", cfg.DBDirectory, err)
	}
	cfg.DBDirectory = absDBDir
	dbDirStat, err := os.Stat(cfg.DBDirectory)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("Database directory '%v' does not exist", cfg.DBDirectory)
	}
	if !dbDirStat.IsDir() {
		return nil, fmt.Errorf("Database directory '%v' exists but is not a directory", cfg.DBDirectory)
	}

	if cfg.DBtype == "sqlite3" {
		if cfg.DBdsn == "" {
			cfg.DBdsn = "file:" + filepath.Join(cfg.DBDirectory, "sqlite3/data.sqlite3")
		}
		dataDir := filepath.Join(cfg.DBDirectory, "sqlite3")
		err = os.MkdirAll(dataDir, os.ModePerm)
		if err != nil {
			return nil, fmt.Errorf("Could not create sqlite3 database directory '%s': %v", dataDir, err)
		}
	}
	bdbDir := filepath.Join(cfg.DBDirectory, "badger")
	err = os.MkdirAll(bdbDir, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("Could not create badger database directory '%s': %v", bdbDir, err)
	}

	return cfg, nil
}

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
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/ini.v1"
)

var KarmaStartTime = time.Now()

type ConfigData struct {
	Username    string `ini:"Username"`
	AccessToken string `ini:"AccessToken"`
	Homeserver  string `ini:"Homeserver"`
	DebugLevel  int    `ini:"DebugLevel"`
	Autojoin    bool   `ini:"Autojoin"`
	Database    string `ini:"Database"`
}

func ReadConfig(ConfigFile string) (*ConfigData, *sql.DB) {
	cfg := new(ConfigData)
	cfg.Username = ""
	cfg.AccessToken = ""
	cfg.Homeserver = ""
	cfg.DebugLevel = 0
	cfg.Autojoin = false
	cfg.Database = "/var/db/karma-bot/sqlite3.db"

	err := ini.MapTo(cfg, ConfigFile)
	if err != nil {
		fmt.Printf("Failed to read config file '%v': %v\n", ConfigFile, err)
		os.Exit(1)
	}
	if cfg.Username == "" {
		fmt.Printf("Config file does not have 'Homeserver'\n")
		os.Exit(1)
	}
	if cfg.AccessToken == "" {
		fmt.Printf("Config file does not have 'Username'\n")
		os.Exit(1)
	}
	if cfg.Homeserver == "" {
		fmt.Printf("Config file does not have 'AccessToken'\n")
		os.Exit(1)
	}
	dbDir := filepath.Dir(cfg.Database)
	dbDirStat, err := os.Stat(dbDir)
	if os.IsNotExist(err) {
		fmt.Printf("Database directory '%v' does not exist\n", dbDir)
		os.Exit(1)
	}
	if !dbDirStat.IsDir() {
		fmt.Printf("Database directory '%v' exists but is not a directory\n", dbDir)
		os.Exit(1)
	}

	db := InitDB(cfg.Database)

	return cfg, db
}

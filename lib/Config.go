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
	DBDirectory string `ini:"DBDirectory"`
}

func ReadConfig(ConfigFile string) (*ConfigData, *sql.DB, *BDBStore) {
	cfg := new(ConfigData)
	cfg.Username = ""
	cfg.AccessToken = ""
	cfg.Homeserver = ""
	cfg.DebugLevel = 0
	cfg.Autojoin = false
	cfg.DBDirectory = "/var/db/karma-bot"

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
	dbDirStat, err := os.Stat(cfg.DBDirectory)
	if os.IsNotExist(err) {
		fmt.Printf("Database directory '%v' does not exist\n", cfg.DBDirectory)
		os.Exit(1)
	}
	if !dbDirStat.IsDir() {
		fmt.Printf("Database directory '%v' exists but is not a directory\n", cfg.DBDirectory)
		os.Exit(1)
	}
	voteDir := filepath.Join(cfg.DBDirectory, "votes")
	err = os.MkdirAll(voteDir, os.ModePerm)
	if err != nil {
		fmt.Printf("Could not create votes database directory '%s': %v", voteDir, err)
		os.Exit(1)
	}
	bdbDir := filepath.Join(cfg.DBDirectory, "badger")
	err = os.MkdirAll(bdbDir, os.ModePerm)
	if err != nil {
		fmt.Printf("Could not create badger database directory '%s': %v", bdbDir, err)
		os.Exit(1)
	}

	votesDB := InitDB(filepath.Join(voteDir, "votes.db"))
	bdbStore := NewBDBStore(bdbDir)
	return cfg, votesDB, bdbStore
}

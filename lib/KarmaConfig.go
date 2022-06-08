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
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/ini.v1"
)

type UnveilInfo struct {
	Dir   string
	Perms string
}

type KarmaConfig struct {
	Username       string   `ini:"Username"`
	AccessToken    string   `ini:"AccessToken"`
	Homeserver     string   `ini:"Homeserver"`
	Autojoin       bool     `ini:"Autojoin"`
	DBDirectory    string   `ini:"DBDirectory"`
	DBtype         string   `ini:"DBtype"`
	DBdsn          string   `ini:"DBdsn"`
	ResponseFreq   int64    `ini:"ResponseFreq"`
	PositiveEmojis string   `ini:"PositiveEmojis"`
	NegativeEmojis string   `ini:"NegativeEmojis"`
	UnveilDirs     []string `init:"UnveilDirs"`
	UnveilInfo     []UnveilInfo
}

func ReadConfig(ConfigFile string) (*KarmaConfig, error) {
	var err error
	var absDBDir string
	var bdbDir string
	var dataDir string
	var dbDirStat os.FileInfo
	var i int
	var uinfo string
	var udir []string

	cfg := new(KarmaConfig)
	cfg.Username = ""
	cfg.AccessToken = ""
	cfg.Homeserver = ""
	cfg.Autojoin = false
	cfg.DBDirectory = "/var/db/karma-bot"
	cfg.ResponseFreq = 5000000 // 5 seconds
	cfg.PositiveEmojis = "❤️,👍️,💯,🍌,🎉,💞,💗,💓,💖,💘,💝,💕,😻,😍,❤️‍🔥"
	cfg.NegativeEmojis = "👎️,💔,😠,👿,🙁,☹️,🤬,☠️,💀"
	cfg.UnveilDirs = []string{}

	// valid SQL driver name: sqlite3, mysql, pgx
	cfg.DBtype = "sqlite3"
	cfg.DBdsn = ""

	err = ini.MapTo(cfg, ConfigFile)
	if err != nil {
		err = fmt.Errorf("Failed to read config file '%s': %v", ConfigFile, err)
		goto failed
	}
	if cfg.Username == "" {
		err = fmt.Errorf("Config file does not have 'Homeserver'")
		goto failed
	}
	if cfg.AccessToken == "" {
		err = fmt.Errorf("Config file does not have 'Username'")
		goto failed
	}
	if cfg.Homeserver == "" {
		err = fmt.Errorf("Config file does not have 'AccessToken'")
		goto failed
	}
	absDBDir, err = filepath.Abs(cfg.DBDirectory)
	if err != nil {
		err = fmt.Errorf("Could not get absolute path of DBDirectory (%s): %v", cfg.DBDirectory, err)
		goto failed
	}
	cfg.DBDirectory = absDBDir
	dbDirStat, err = os.Stat(cfg.DBDirectory)
	if os.IsNotExist(err) {
		err = fmt.Errorf("Database directory '%s' does not exist", cfg.DBDirectory)
		goto failed
	}
	if !dbDirStat.IsDir() {
		err = fmt.Errorf("Path '%s' exists but is not a directory", cfg.DBDirectory)
		goto failed
	}

	if cfg.DBtype != "sqlite3" && cfg.DBtype != "pgx" && cfg.DBtype != "mysql" {
		err = fmt.Errorf("Unknown database type %q - accepted values are \"mysql\", \"postgresql\", \"sqlite3\"", cfg.DBtype)
		goto failed
	}
	if cfg.DBtype == "sqlite3" {
		cfg.DBdsn = "file:" + filepath.Join(cfg.DBDirectory, "sqlite3", "data.sqlite3")
		dataDir = filepath.Join(cfg.DBDirectory, "sqlite3")
		if _, err = os.Stat(dataDir); errors.Is(err, os.ErrNotExist) {
			err = os.Mkdir(dataDir, os.ModePerm)
			if err != nil {
				err = fmt.Errorf("Could not create sqlite3 database directory '%s': %v", dataDir, err)
				goto failed
			}
		}
	}
	bdbDir = filepath.Join(cfg.DBDirectory, "badger")
	if _, err = os.Stat(bdbDir); errors.Is(err, os.ErrNotExist) {
		err = os.Mkdir(bdbDir, os.ModePerm)
		if err != nil {
			err = fmt.Errorf("Could not create badger database directory '%s': %v", bdbDir, err)
			goto failed
		}
	}

	i = len(cfg.UnveilDirs)
	cfg.UnveilInfo = make([]UnveilInfo, i, i)
	for i, uinfo = range cfg.UnveilDirs {
		udir = strings.SplitN(uinfo, ":", 2)
		if len(udir) < 2 || udir[0] == "" || udir[1] == "" {
			err = fmt.Errorf("Could not get unveil information from '%s'", uinfo)
			goto failed
		}
		cfg.UnveilInfo[i].Perms = udir[0]
		cfg.UnveilInfo[i].Dir = udir[1]
	}

	return cfg, nil

failed:
	return nil, err
}

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
	"sort"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/v4/stdlib"
	_ "github.com/mattn/go-sqlite3"
)

type SQLStore struct {
	DB     *sql.DB
	DBtype string
	Logger *BotLogger
}

func NewSQLStore(DBtype, DBdsn string, b *BotLogger) (*SQLStore, error) {
	var err error
	var sqlStore *SQLStore

	sqlStore = new(SQLStore)
	sqlStore.DBtype = DBtype
	sqlDB, err := sql.Open(DBtype, DBdsn)
	if err != nil {
		return nil, err
	}
	err = sqlDB.Ping()
	if err != nil {
		return nil, err
	}
	sqlStore.DB = sqlDB
	sqlStore.Logger = b
	return sqlStore, nil
}

func (s *SQLStore) Close() {
	s.DB.Close()
}

func (s *SQLStore) GetVersion() (BotVersion, error) {
	var cver BotVersion
	err := s.DB.QueryRow(`SELECT major,minor,patch FROM version;`).Scan(&cver.Major, &cver.Minor, &cver.Patch)
	if err != nil {
		s.Logger.Warnf("Failed to get version: %v", err)
		s.Logger.Warnf("Using default version 0.0.0")
		return BotVersion{0, 0, 0, nil}, err
	}
	return cver, nil
}

func (s *SQLStore) UpdateDB(dbPatches BotVersionArr) error {
	cver, _ := s.GetVersion()
	s.Logger.Infof("Current database version is %v.%v.%v", cver.Major, cver.Minor, cver.Patch)
	s.Logger.Infof("Upgrading database to latest version...")
	sort.Sort(dbPatches)
	for _, kver := range dbPatches {
		if !BVLess(cver, kver) {
			continue
		}
		s.Logger.Infof("Applying patch version %v.%v.%v", kver.Major, kver.Minor, kver.Patch)
		err := kver.SQLPatch(s.DB, s.DBtype)
		if err != nil {
			s.Logger.Errorf("Failed to patch: %v", err)
			s.Logger.Errorf("Aborting update of the database")
			return err
		}
	}
	s.Logger.Infof("Update finished")
	return nil
}

var SQLKarmaPatches = []BotVersion{SQLPatchv_1_0_0}

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

	"go.uber.org/zap"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/v4/stdlib"
	_ "github.com/mattn/go-sqlite3"

	sqlp "bsd.ac/karma-bot/lib/sqlp"
)

type SQLStore struct {
	DB *sql.DB
}

func (s *SQLStore) SGet(key string) (string, error) {
	var val string
	err := s.DB.QueryRow(`SELECT data FROM stateStore WHERE key = '$1';`, key).Scan(&val)
	if err != nil {
		zap.S().Errorf("SGet failed to fetch values for '%s': %v", key, err)
		return "", err
	}
	return val, err
}

func (s *SQLStore) SSet(key, val string) error {
	res, err := s.DB.Exec(`UPDATE dataStore SET val = '$2' WHERE key = '$1';`, key, val)
	nRows, _ := res.RowsAffected()
	if err != nil && nRows == 0 {
		res, err = s.DB.Exec(`INSERT INTO dataStore(key, val) values('$1', '$2');`, key, val)
		if err != nil {
			zap.S().Errorf("SSet could not update or insert: %v", err)
		}
	}
	return err
}

func (s *SQLStore) Get(key []byte) ([]byte, error) {
	str := string(key)
	val, err := s.SGet(str)
	return []byte(val), err
}

func (s *SQLStore) Set(key, val []byte) error {
	skey := string(key)
	sval := string(val)
	return s.SSet(skey, sval)
}

func (s *SQLStore) GetVersion() (BotVersion, error) {
	var cver BotVersion
	err := s.DB.QueryRow(`SELECT major,minor,patch FROM version;`).Scan(&cver.Major, &cver.Minor, &cver.Patch)
	if err != nil {
		zap.S().Warnf("Failed to get version: %v", err)
		zap.S().Warnf("Using default version 0.0.0")
		return BotVersion{0, 0, 0, nil}, err
	}
	return cver, nil
}

func (s *SQLStore) UpdateDB(dbPatches BotVersionArr) error {
	cver, _ := s.GetVersion()
	zap.S().Infof("Current database version is %v.%v.%v", cver.Major, cver.Minor, cver.Patch)
	zap.S().Infof("Upgrading database to latest version...")
	sort.Sort(dbPatches)
	driverType := s.DB.Driver()
	driverName, err := SQLDriverName(driverType)
	if err != nil {
		zap.S().Errorf("Failed to find database driver name: %v", err)
		return err
	}
	zap.S().Infof("Current database driver in use: %s", driverName)
	for _, kver := range dbPatches {
		if !BVLess(cver, kver) {
			continue
		}
		zap.S().Infof("Applying patch version %v.%v.%v", kver.Major, kver.Minor, kver.Patch)
		err := kver.SQLPatch(s.DB, driverName)
		if err != nil {
			zap.S().Errorf("Failed to patch: %v", err)
			zap.S().Errorf("Aborting update of the database")
			return err
		}
	}
	zap.S().Infof("Update finished")
	return nil
}

func NewSQLStore(DBtype, DBdsn string) (*SQLStore, error) {
	zap.S().Debugf("Opening SQL store of type '%s' with DSN: %s", DBtype, DBdsn)
	sqlStore := new(SQLStore)
	sqlDB, err := sql.Open(DBtype, DBdsn)
	if err != nil {
		zap.S().Errorf("Could not open the database: %v", err)
		return nil, err
	}
	err = sqlDB.Ping()
	if err != nil {
		zap.S().Errorf("Could not ping the database: %v", err)
		return nil, err
	}
	sqlStore.DB = sqlDB
	dbPatches := BotVersionArr{
		BotVersion{1, 0, 0, sqlp.SQLpatchv_1_0_0},
	}
	err = sqlStore.UpdateDB(dbPatches)
	return sqlStore, err
}

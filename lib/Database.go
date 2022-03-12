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
	"sort"

	_ "github.com/mattn/go-sqlite3"
	"maunium.net/go/mautrix/id"
)

func GetVersion(db *sql.DB) KarmaVersion {
	var cver KarmaVersion
	err := db.QueryRow("SELECT major,minor,patch FROM version;").Scan(&cver.Major, &cver.Minor, &cver.Patch)
	if err != nil {
		return KarmaVersion{0, 0, 0, nil}
	}
	return cver
}

func InitRoom(db *sql.DB, rID id.RoomID) bool {
	_, err := db.Query("SELECT * FROM room_tables WHERE room_name = " + rID.String())
	if err != nil {
		rID_table_name := RandomString(15)
		_, err = db.Exec(fmt.Sprintf("INSERT INTO room_tables(room_name, room_table_name) values('%v', '%v');", rID.String(), rID_table_name))
		if err != nil {
			fmt.Printf("Error while inserting room_table_name for '%v'\n", rID.String())
			return false
		}
		_, err = db.Exec(fmt.Sprintf("CREATE TABLE IF NOT EXISTS %v LIKE room_table_example;"))
		if err != nil {
			fmt.Printf("Error while cloning the default table")
			return false
		}
	}
	return true
}

func GetRoomTable(db *sql.DB, room_name string) string {
	var room_table_name string
	err := db.QueryRow("SELECT room_table_name FROM room_tables WHERE room_name = " + room_name).Scan(&room_table_name)
	if err != nil {
		return ""
	}
	return room_table_name
}

func InitDB(database string) *sql.DB {
	db, err := sql.Open("sqlite3", database)
	if err != nil {
		fmt.Printf("Error while opening the database '%v': '%v'\n", database, err)
		os.Exit(1)
	}

	cver := GetVersion(db)
	fmt.Printf("Current database version is %v.%v.%v\n", cver.Major, cver.Minor, cver.Patch)
	fmt.Printf("Upgrading database to latest version...\n")
	sort.Sort(KVPatches)
	for _, kver := range KVPatches {
		if !KVLess(cver, kver) {
			continue
		}
		fmt.Printf("Applying patch version %v.%v.%v\n", kver.Major, kver.Minor, kver.Patch)
		kver.SQLPatch(db)
	}

	return db
}

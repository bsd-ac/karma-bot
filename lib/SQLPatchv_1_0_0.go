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

	_ "github.com/mattn/go-sqlite3"
)

func sqlpatch(db *sql.DB) bool {
	sqlquery := `CREATE TABLE IF NOT EXISTS version (present BOOL PRIMARY KEY DEFAULT TRUE, major INTEGER NOT NULL, minor INTEGER NOT NULL, patch INTEGER NOT NULL, CONSTRAINT present_uniq CHECK (present));
INSERT INTO version(present, major, minor, patch) values(1, 1, 0, 0);
CREATE TABLE IF NOT EXISTS room_tables (room_name STRING PRIMARY KEY NOT NULL, room_table_name);
CREATE TABLE IF NOT EXISTS room_table_example (senderID STRING NOT NULL, targetID STRING NOT NULL, eventID STRING NOT NULL, vote INTEGER NOT NULL, PRIMARY KEY(senderID, targetID, eventID));
`
	_, err := db.Exec(sqlquery)
	if err != nil {
		fmt.Printf("Error while applying patch 1.0.0: %v\n", err)
		return false
	}
	return true
}

var SQLPatchv_1_0_0 = KarmaVersion{1, 0, 0, sqlpatch}

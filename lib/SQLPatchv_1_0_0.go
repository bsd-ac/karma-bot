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
)

func SQLPatchv_1_0_0_(db *sql.DB, dbType string) error {
	sqlquery := "CREATE TABLE version (present BOOL PRIMARY KEY DEFAULT TRUE, major INTEGER NOT NULL, minor INTEGER NOT NULL, patch INTEGER NOT NULL, CONSTRAINT present_uniq CHECK (present)); CREATE TABLE events (senderID VARCHAR NOT NULL, targetID VARCHAR NOT NULL, eventID VARCHAR, roomID VARCHAR, vote INTEGER NOT NULL, PRIMARY KEY(eventID, roomID)); CREATE TABLE optout (uidHash VARCHAR NOT NULL, PRIMARY KEY(uidHash)); INSERT INTO version(present, major, minor, patch) values(TRUE, 1, 0, 0);"
	_, err := db.Exec(sqlquery)
	return err
}

var SQLPatchv_1_0_0 = BotVersion{1, 0, 0, SQLPatchv_1_0_0_}

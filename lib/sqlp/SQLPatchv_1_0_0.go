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

	"go.uber.org/zap"
)

func SQLpatchv_1_0_0(db *sql.DB, driverName string) error {
	sqlquery := "CREATE TABLE IF NOT EXISTS version (present BOOL PRIMARY KEY DEFAULT TRUE, major INTEGER NOT NULL, minor INTEGER NOT NULL, patch INTEGER NOT NULL, CONSTRAINT present_uniq CHECK (present)); CREATE TABLE IF NOT EXISTS stateStore (key BYTEA PRIMARY KEY, val BYTEA); CREATE TABLE IF NOT EXISTS votes (senderID VARCHAR NOT NULL, targetID VARCHAR NOT NULL, eventID VARCHAR, roomID VARCHAR, vote INTEGER NOT NULL, PRIMARY KEY(senderID, targetID, eventID, roomID)); INSERT INTO version(present, major, minor, patch) values(TRUE, 1, 0, 0);"
	_, err := db.Exec(sqlquery)
	if err != nil {
		zap.S().Errorf("Error while applying patch 1.0.0: %v", err)
		return err
	}
	return nil
}

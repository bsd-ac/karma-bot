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
	"path/filepath"
	"os"
	"testing"
)

func TestKarmaUtils(t *testing.T) {
	dbDir, err := os.MkdirTemp("", "*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dbDir)

	bLogger := NewBotLogger()
	sqlStore, err := NewSQLStore("sqlite3", "file:" + filepath.Join(dbDir, "data.sqlite3"), bLogger)
	if err != nil {
		t.Fatal(err)
	}
	err = sqlStore.UpdateDB(SQLKarmaPatches)
	if err != nil {
		t.Fatal(err)
	}
	kBot := new(KarmaBot)
	kBot.sqlDB = sqlStore
	kBot.logger = bLogger

	t.Log("Starting the test now")

	////// t1
}

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
	t.Logf("Using %q as tempdir", dbDir)
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

	userA := "@banana-bot:matrix.org"
	userB := "@jane-doe:matrix.org"
	event := "some-random-event-matrix.org"
	roomA := "some-cool-room-matrix.org"
	roomB := "other-cool-room-matrix.org"
	var vote int64
	vote = 1

	////// t1
	kBot.KarmaAdd(userA, userB, event, roomA, vote)
	if kBot.GetKarmaTotal(userB) != 1 {
		t.Errorf("t1.1 failure")
	}
	if kBot.GetKarmaTotal(userA) != 0 {
		t.Errorf("t1.2 failure")
	}

	////// t2
	kBot.OptOut(userA)
	if kBot.GetKarmaTotal(userB) != 0 {
		t.Errorf("t2.1 failure")
	}
	if kBot.GetKarmaTotal(userA) != 0 {
		t.Errorf("t2.2 failure")
	}

	////// t3
	kBot.KarmaAdd(userA, userB, event, roomA, vote)
	if kBot.GetKarmaTotal(userB) != 0 {
		t.Errorf("t3.1 failure")
	}
	if kBot.GetKarmaTotal(userA) != 0 {
		t.Errorf("t3.2 failure")
	}

	////// t4
	kBot.OptIn(userA)
	kBot.KarmaAdd(userA, userB, event, roomA, vote)
	kBot.KarmaAdd(userA, userB, event, roomA, vote)
	if kBot.GetKarmaTotal(userB) != 1 {
		t.Errorf("t4.1 failure")
	}
	if kBot.GetKarmaTotal(userA) != 0 {
		t.Errorf("t4.2 failure")
	}

	////// t5
	kBot.KarmaAdd(userA, userB, event, roomB, vote)
	if kBot.GetKarmaTotal(userB) != 2 {
		t.Errorf("t5.1 failure")
	}
	if kBot.GetKarma(userB, roomA) != 1 {
		t.Errorf("t5.2 failure")
	}
	if kBot.GetKarma(userB, roomB) != 1 {
		t.Errorf("t5.3 failure")
	}
	if kBot.GetKarmaTotal(userA) != 0 {
		t.Errorf("t5.4 failure")
	}

}

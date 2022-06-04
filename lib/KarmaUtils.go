/*
 * Copyright (c) 2022 Aisha Tammy <aisha@bsd.ac>
 * Copyright (c) 2021 Aaron Bieber <aaron@bolddaemon.com>
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
	"golang.org/x/crypto/blake2b"
)

func uidHash(userID string) []byte {
	optOut := []byte("optOutUID_")
	uid := []byte(userID)
	blake2b_uid := blake2b.Sum512(uid)
	uid_key := append(optOut[:], blake2b_uid[:]...)
	return uid_key
}

func (kBot *KarmaBot) IsOptOut(userID string) bool {
	uhash := uidHash(userID)
	query := `SELECT COUNT(*) FROM optout WHERE uidHash = ?`
	var ucount int64
	err := kBot.sqlDB.DB.QueryRow(query, uhash).Scan(&ucount)
	if err == nil && ucount > 0 {
		return true
	}
	if err != nil {
		kBot.logger.Warnf("Error in IsOptOut for user %q: %v", userID, err)
	}
	return false
}

func (kBot *KarmaBot) OptOut(userID string) {
	query := `DELETE FROM events WHERE senderID = ? OR targetID = ?`
	_, err := kBot.sqlDB.DB.Exec(query, userID, userID)
	if err != nil {
		kBot.logger.Warnf("Error in OptOut while deleting user from events %q: %v", userID, err)
	}
	uhash := uidHash(userID)
	query = `INSERT INTO optout (uidHash) VALUES (?)`
	_, err = kBot.sqlDB.DB.Exec(query, uhash)
	if err != nil {
		kBot.logger.Warnf("Error in OptOut while inserting hash for user %q: %v", userID, err)
	}
}

func (kBot *KarmaBot) OptIn(userID string) {
	uhash := uidHash(userID)
	query := `DELETE FROM optout WHERE uidHash = ?`
	_, err := kBot.sqlDB.DB.Exec(query, uhash)
	if err != nil {
		kBot.logger.Warnf("Error in OptIn for user %q: %v", userID, err)
	}
}

func (kBot *KarmaBot) GetKarma(userID, roomID string) int64 {
	query := `SELECT SUM(vote) FROM events WHERE targetID = ? AND roomID = ?`
	var karma int64
	err := kBot.sqlDB.DB.QueryRow(query, userID, roomID).Scan(&karma)
	if err != nil {
		kBot.logger.Warnf("Error in GetKarm for user %q: %v", userID, err)
		karma = 0
	}
	return karma
}

func (kBot *KarmaBot) GetKarmaTotal(userID string) int64 {
	query := `SELECT SUM(vote) FROM events WHERE targetID = ?`
	var karma int64
	err := kBot.sqlDB.DB.QueryRow(query, userID).Scan(&karma)
	if err != nil {
		kBot.logger.Warnf("Error in GetKarmaTotal for user %q: %v", userID, err)
		karma = 0
	}

	return karma
}

func (kBot *KarmaBot) KarmaAdd(senderID, targetID, eventID, roomID string, vote int64) {
	if kBot.IsOptOut(senderID) || kBot.IsOptOut(targetID) {
		return
	}
	query := `INSERT INTO events (senderID, targetID, eventID, roomID, vote) VALUES (?, ?, ?, ?, ?)`
	_, err := kBot.sqlDB.DB.Exec(query, senderID, targetID, eventID, roomID, vote)
	if err != nil {
		kBot.logger.Warnf("Error in KarmaAdd for (%s, %s, %s, %s, %d): %v", senderID, targetID, eventID, roomID, vote, err)
	}
}

func (kBot *KarmaBot) KarmaDelete(eventID, roomID string) {
	query := `DELETE FROM events WHERE eventID = ? AND roomID = ?`
	_, err := kBot.sqlDB.DB.Exec(query, eventID, roomID)
	if err != nil {
		kBot.logger.Warnf("Error in KarmaDelete for (%s, %s): %v", eventID, roomID, err)
	}
}

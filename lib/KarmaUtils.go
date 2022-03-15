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
	"encoding/binary"

	"golang.org/x/crypto/blake2b"
)

/*
"go.uber.org/zap"
"maunium.net/go/mautrix"
"maunium.net/go/mautrix/id"
*/

func uid2OptOutKey(userID string) []byte {
	optOut := []byte("optOutUID_")
	uid := []byte(userID)
	blake2b_uid := blake2b.Sum512(uid)
	uid_key := append(optOut[:], blake2b_uid[:]...)
	return uid_key
}

func uid2KarmaKey(userID string) []byte {
	return []byte("karmaUID_" + userID)
}

func KarmaIsOptOut(userID string, s *BDBStore) (bool, error) {
	uid_key := uid2OptOutKey(userID)
	_, err := s.Get(uid_key)
	if err != nil {
		return false, err
	}
	return true, nil
}

func KarmaOptOut(userID string, s *BDBStore) error {
	uid_key := uid2OptOutKey(userID)
	val := []byte{0}
	return s.Set(uid_key, val)
}

func KarmaOptIn(userID string, s *BDBStore) error {
	uid_key := uid2OptOutKey(userID)
	return s.Delete(uid_key)
}

func GetKarma(userID string, s *BDBStore) (int64, error) {
	val, err := s.Get(uid2KarmaKey(userID))
	if err != nil {
		return 0, err
	}
	karma, _ := binary.Varint(val)
	return karma, err
}

func SetKarma(userID string, karma int64, s *BDBStore) error {
	buf := make([]byte, binary.MaxVarintLen64)
	binary.PutVarint(buf, karma)
	return s.Set(uid2KarmaKey(userID), buf)
}

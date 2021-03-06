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
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/id"

	badger "github.com/dgraph-io/badger/v3"
)

/*
   Main functions for usage

   Set(key, val []byte) error       - set value of 'key' to 'val' (as []byte)
   Get(key []byte) ([]byte, error)  - get value of 'key'          (as []byte)
   SSet(key, val string) error      - set value of 'key' to 'val' (as string)
   SGet(key string) ([]byte, error) - get value of 'key'          (as string)
*/
type BDBStore struct {
	DB     *badger.DB
	Logger *BotLogger
}

func NewBDBStore(dbPath string, b *BotLogger) (*BDBStore, error) {
	opts := badger.DefaultOptions(dbPath)
	opts.Logger = b
	bdb, err := badger.Open(opts)
	if err != nil {
		return nil, err
	}
	bdbStore := new(BDBStore)
	bdbStore.DB = bdb
	bdbStore.Logger = b
	return bdbStore, nil
}

func (s *BDBStore) Close() {
	s.DB.Close()
}

func (s *BDBStore) Get(key []byte) ([]byte, error) {
	var val []byte
	err := s.DB.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			val = nil
			return err
		}
		item.Value(func(iVal []byte) error {
			val = append([]byte{}, iVal...)
			return nil
		})
		return nil
	})
	return val, err
}

func (s *BDBStore) SGet(key string) (string, error) {
	byt, err := s.Get([]byte(key))
	return string(byt), err
}

func (s *BDBStore) Set(key, val []byte) error {
	var err error
	for i := 0; i < 3; i++ {
		err = s.DB.Update(func(txn *badger.Txn) error {
			var err error
			err = txn.Set(key, val)
			if err != nil {
				if err == badger.ErrTxnTooBig {
					txn.Commit()
				}
			}
			return err
		})
		if err == nil || (err != badger.ErrTxnTooBig && err != badger.ErrConflict) {
			return err
		}
	}
	return err
}

func (s *BDBStore) SSet(key, val string) error {
	return s.Set([]byte(key), []byte(val))
}

func (s *BDBStore) Delete(key []byte) error {
	err := s.DB.Update(func(txn *badger.Txn) error {
		err := txn.Delete(key)
		return err
	})
	return err
}

func (s *BDBStore) SDelete(key string) error {
	return s.Delete([]byte(key))
}

func (s *BDBStore) SaveFilterID(userID id.UserID, filterID string) {
	uid := "userid_filter_" + userID.String()
	err := s.Set([]byte(uid), []byte(filterID))
	if err != nil {
		s.Logger.Errorf("Error in SaveFilterID(%s, %s): %v", userID.String(), filterID, err)
	}
}

func (s *BDBStore) LoadFilterID(userID id.UserID) string {
	uid := "userid_filter_" + userID.String()
	filter, err := s.Get([]byte(uid))
	if err != nil {
		s.Logger.Errorf("Error in LoadFilterID(%s): %v", userID.String(), err)
		return ""
	}
	return string(filter[:])
}

func (s *BDBStore) SaveNextBatch(userID id.UserID, nextBatchToken string) {
	uid := "userid_batch_" + userID.String()
	err := s.Set([]byte(uid), []byte(nextBatchToken))
	if err != nil {
		s.Logger.Errorf("Error in SaveNextBatch(%s, %s): %v", userID.String(), nextBatchToken, err)
	}
}

func (s *BDBStore) LoadNextBatch(userID id.UserID) string {
	uid := "userid_batch_" + userID.String()
	batch, err := s.Get([]byte(uid))
	if err != nil {
		s.Logger.Errorf("Error in LoadNextBatch(%s): %v", userID.String(), err)
		return ""
	}
	return string(batch[:])
}

func (s *BDBStore) SaveRoom(room *mautrix.Room) {
	rdata, _ := EncodeRoom(room)
	rid := "roomid_" + room.ID.String()
	err := s.Set([]byte(rid), rdata)
	if err != nil {
		s.Logger.Errorf("Error in SaveRoom(%d): %v", room.ID.String(), err)
	}
}

func (s *BDBStore) LoadRoom(roomID id.RoomID) *mautrix.Room {
	rid := "roomid_" + roomID.String()
	rdata, err := s.Get([]byte(rid))
	if err != nil {
		s.Logger.Errorf("Error in LoadRoom(%s) while getting data: %v", roomID.String(), err)
		return nil
	}
	room, _ := DecodeRoom([]byte(rdata))
	if err != nil {
		s.Logger.Errorf("Error in LoadRoom(%s) while decoding room data: %v", roomID.String(), err)
		return nil
	}
	return room
}

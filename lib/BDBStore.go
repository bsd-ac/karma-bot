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
	"go.uber.org/zap"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/id"

	badger "github.com/dgraph-io/badger/v3"
)

type zapLogger struct {
	L *zap.SugaredLogger
}

func (z *zapLogger) Errorf(f string, v ...interface{}) {
	z.L.Errorf(f, v...)
}
func (z *zapLogger) Warningf(f string, v ...interface{}) {
	z.L.Warnf(f, v...)
}
func (z *zapLogger) Infof(f string, v ...interface{}) {
	z.L.Infof(f, v...)
}
func (z *zapLogger) Debugf(f string, v ...interface{}) {
	z.L.Debugf(f, v...)
}

/*
   Main functions for usage

   Set(key, val []byte) error       - set value of 'key' to 'val' (as []byte)
   Get(key []byte) ([]byte, error)  - get value of 'key'          (as []byte)
   SSet(key, val string) error      - set value of 'key' to 'val' (as string)
   SGet(key string) ([]byte, error) - get value of 'key'          (as string)
*/
type BDBStore struct {
	DB *badger.DB
}

func NewBDBStore(dbPath string) (*BDBStore, error) {
	opts := badger.DefaultOptions(dbPath)
	zlog := new(zapLogger)
	zlog.L = zap.S()
	opts.Logger = zlog
	bdb, err := badger.Open(opts)
	if err != nil {
		return nil, err
	}
	bdbStore := new(BDBStore)
	bdbStore.DB = bdb
	return bdbStore, nil
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
				return err
			} else {
				return nil
			}
		})
		if err == nil {
			return nil
		}
		if err != badger.ErrTxnTooBig && err != badger.ErrConflict {
			return err
		}
	}
	return err
}

func (s *BDBStore) SSet(key, val string) error {
	return s.Set([]byte(key), []byte(val))
}

func (s *BDBStore) SaveFilterID(userID id.UserID, filterID string) {
	uid := "userid_filter_" + userID.String()
	_ = s.Set([]byte(uid), []byte(filterID))
}

func (s *BDBStore) LoadFilterID(userID id.UserID) string {
	uid := "userid_filter_" + userID.String()
	filter, _ := s.Get([]byte(uid))
	return string(filter[:])
}

func (s *BDBStore) SaveNextBatch(userID id.UserID, nextBatchToken string) {
	uid := "userid_batch_" + userID.String()
	s.Set([]byte(uid), []byte(nextBatchToken))
}

func (s *BDBStore) LoadNextBatch(userID id.UserID) string {
	uid := "userid_batch_" + userID.String()
	batch, _ := s.Get([]byte(uid))
	return string(batch[:])
}

func (s *BDBStore) SaveRoom(room *mautrix.Room) {
	rdata, _ := EncodeRoom(room)
	rid := "roomid_" + room.ID.String()
	s.Set([]byte(rid), rdata)
}

func (s *BDBStore) LoadRoom(roomID id.RoomID) *mautrix.Room {
	rid := "roomid_" + roomID.String()
	rdata, _ := s.Get([]byte(rid))
	room, _ := DecodeRoom([]byte(rdata))
	return room
}

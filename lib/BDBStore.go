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
	"bytes"
	"encoding/gob"

	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/id"

	badger "github.com/dgraph-io/badger/v3"
)

type BDBStore struct {
	BDB *badger.DB
}

func NewBDBStore(dbPath string) *BDBStore {
	bdb, err := badger.Open(badger.DefaultOptions(dbPath))
	if err != nil {
		return nil
	}
	bdbStore := new(BDBStore)
	bdbStore.BDB = bdb
	return bdbStore
}

func (s *BDBStore) encodeRoom(room *mautrix.Room) ([]byte, error) {
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	err := enc.Encode(room)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (s *BDBStore) decodeRoom(room []byte) (*mautrix.Room, error) {
	var r *mautrix.Room
	buf := bytes.NewBuffer(room)
	dec := gob.NewDecoder(buf)
	err := dec.Decode(&r)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (s *BDBStore) Get(key []byte) ([]byte, error) {
	var val []byte
	txn := s.BDB.NewTransaction(false)
	item, err := txn.Get(key)
	if err != nil {
		return nil, err
	}
	err = item.Value(func(iVal []byte) error {
		val = append([]byte{}, iVal...)
		return nil
	})
	txn.Discard()
	return val, err
}

func (s *BDBStore) Set(key []byte, val []byte) error {
	txn := s.BDB.NewTransaction(true)
	err := txn.Set(key, val)
	if err != nil {
		return err
	}
	err = txn.Commit()
	return err
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
	rdata, _ := s.encodeRoom(room)
	rid := "roomid_" + room.ID.String()
	s.Set([]byte(rid), rdata)
}

func (s *BDBStore) LoadRoom(roomID id.RoomID) *mautrix.Room {
	rid := "roomid_" + roomID.String()
	rdata, _ := s.Get([]byte(rid))
	room, _ := s.decodeRoom([]byte(rdata))
	return room
}

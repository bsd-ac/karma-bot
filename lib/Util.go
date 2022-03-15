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
	"bytes"
	"database/sql"
	"database/sql/driver"
	"encoding/gob"
	"fmt"
	"math/rand"
	"reflect"
	"time"

	"maunium.net/go/mautrix"
)

var RNG = rand.New(rand.NewSource(time.Now().UnixNano()))

func RandomStringFromChars(length int, chars string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = chars[RNG.Intn(len(chars))]
	}
	return string(b)

}

func RandomString(length int) string {
	return RandomStringFromChars(length, "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
}

var sqlDriverMap map[reflect.Type]string

func SQLDriverName(driver driver.Driver) (string, error) {
	if sqlDriverMap == nil {
		sqlDriverMap = map[reflect.Type]string{}
		for _, driverName := range sql.Drivers() {
			db, _ := sql.Open(driverName, "")
			if db != nil {
				driverType := reflect.TypeOf(db.Driver())
				sqlDriverMap[driverType] = driverName
			}
		}
	}
	driverType := reflect.TypeOf(driver)
	driverName, found := sqlDriverMap[driverType]
	if found {
		return driverName, nil
	}
	return "", fmt.Errorf("Could not find driver type and name")
}

func EncodeRoom(room *mautrix.Room) ([]byte, error) {
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	err := enc.Encode(room)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func DecodeRoom(room []byte) (*mautrix.Room, error) {
	var r *mautrix.Room
	buf := bytes.NewBuffer(room)
	dec := gob.NewDecoder(buf)
	err := dec.Decode(&r)
	if err != nil {
		return nil, err
	}
	return r, nil
}

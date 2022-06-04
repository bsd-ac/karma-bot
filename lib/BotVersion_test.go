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
	"sort"
	"testing"
)

func TestBotVersion(t *testing.T) {
	////// t1
	a1 := BotVersion{1, 0, 0, nil}
	a2 := BotVersion{1, 0, 1, nil}

	if BVLess(a2, a1) {
		t.Errorf("Incorrect comparison result - t1.1")
	}

	if !BVLess(a1, a2) {
		t.Errorf("Incorrect comparison result - t1.2")
	}

	////// t2
	b1 := BotVersion{1, 0, 0, nil}
	b2 := BotVersion{1, 0, 0, nil}

	if BVLess(b1, b2) || BVLess(b2, b1) {
		t.Errorf("Incorrect comparison result - t2")
	}

	///// t3
	varray := BotVersionArr{{1, 2, 3, nil}, {3, 1, 2, nil}, {1, 1, 1, nil}, {1, 0, 1, nil}, {1, 1, 1, nil}, {1, 0, 1, nil}}
	sort.Sort(varray)
	vlen := len(varray)
	for i := 0; i < vlen-1; i++ {
		if BVLess(varray[i+1], varray[i]) {
			t.Errorf("Incorrect comparison result - t3.%d", i)
		}
	}
}

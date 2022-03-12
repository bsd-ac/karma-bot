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
)

type KarmaVersion struct {
	Major    int
	Minor    int
	Patch    int
	SQLPatch func(db *sql.DB) bool
}

// strict inequality checker
func KVLess(v1, v2 KarmaVersion) bool {
	if v1.Major != v2.Major {
		return v1.Major < v2.Major
	} else if v1.Minor != v2.Minor {
		return v1.Minor < v2.Minor
	} else if v1.Patch != v2.Patch {
		return v1.Patch < v2.Patch
	}
	return false
}

type KarmaVersionArr []KarmaVersion

func (karr KarmaVersionArr) Len() int {
	return len(karr)
}
func (karr KarmaVersionArr) Swap(i, j int) {
	karr[i], karr[j] = karr[j], karr[i]
}
func (karr KarmaVersionArr) Less(i, j int) bool {
	return KVLess(karr[i], karr[j])
}

var KVPatches = KarmaVersionArr{SQLPatchv_1_0_0}

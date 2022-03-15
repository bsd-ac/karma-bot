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
	"fmt"
	"regexp"
	"strings"

	"go.uber.org/zap"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/id"
)

type GetKarmaBot struct {
}

const gk_rstr = `(?i)\!karma(.*)$`

var gk_rexp = regexp.MustCompile(gk_rstr)

func (u *GetKarmaBot) MatchMessage(body string) bool {
	return gk_rexp.MatchString(body)
}

func (u *GetKarmaBot) ProcessMessage(body string, cli *mautrix.Client, source mautrix.EventSource, evt *event.Event, bdb *BDBStore, sqlDB *SQLStore) error {
	htmlBody := strings.TrimSpace(evt.Content.AsMessage().FormattedBody)
	zap.S().Debugf("Processing html '%s'", htmlBody)
	href := gk_rexp.ReplaceAllString(htmlBody, "$1")
	zap.S().Debugf("Processing href '%s'", href)
	targetID := HTMLToUserID(href)
	zap.S().Debugf("Got userID '%s'", targetID)
	targetID = strings.TrimSpace(targetID)
	if targetID == "" {
		zap.S().Warnf("Could not parse user from html, defaulting to sender")
		targetID = evt.Sender.String() // query self
		dname, err := cli.GetDisplayName(id.UserID(targetID))
		if err != nil {
			zap.S().Warnf("Could not get display name for '%s'", targetID)
			href = "<name unknown>"
		} else {
			href = dname.DisplayName
		}
	}
	targetID = strings.TrimSpace(targetID)

	msg := ""
	optOut, _ := KarmaIsOptOut(targetID, bdb)
	zap.S().Debugf("Current opt out status for '%s': %t", targetID, optOut)
	if optOut {
		msg = "Unknown user"
	} else {
		karma, err := GetKarma(targetID, bdb)
		if err != nil {
			zap.S().Warnf("Could not get karma, defaulting to 0, for '%s': %v", targetID, err)
		}
		msg = fmt.Sprintf("Current karma for %s: %d", href, karma)
	}

	cli.SendText(evt.RoomID, msg)
	return nil
}

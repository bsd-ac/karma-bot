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

	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/format"
	"maunium.net/go/mautrix/id"
)

type OptStatusBot struct {
}

func (u *OptStatusBot) ID() string {
	return "OptStatusBot"
}

func (u *OptStatusBot) NeedsTimer() bool {
	return true
}

func (u *OptStatusBot) Re() *regexp.Regexp {
	return regexp.MustCompile(`(?i)^\!optstatus(.*)$`)
}

func (u *OptStatusBot) ProcessMessage(body string, source mautrix.EventSource, evt *event.Event, kBot *KarmaBot) bool {
	rexp := u.Re()
	if !rexp.MatchString(body) {
		return false
	}
	cli := kBot.mClient
	htmlBody := strings.TrimSpace(evt.Content.AsMessage().FormattedBody)
	href := rexp.ReplaceAllString(htmlBody, "$1")
	targetID := HTMLToUserID(href)
	kBot.logger.Debugf("Got userID '%s'", targetID)
	targetID = strings.TrimSpace(targetID)
	if targetID == "" {
		kBot.logger.Warnf("Could not parse user from html, defaulting to sender")
		targetID = evt.Sender.String() // query self
		dname, err := cli.GetDisplayName(id.UserID(targetID))
		if err != nil {
			kBot.logger.Warnf("Could not get display name for '%s'", targetID)
			href = "<name unknown>"
		} else {
			href = dname.DisplayName
		}
	}

	msg := ""
	optOut := kBot.IsOptOut(targetID)
	optStatus := "in"
	if optOut {
		optStatus = "out of"
	}
	msg = fmt.Sprintf("%s is currently opted %s the system", href, optStatus)

	htmlMsg := format.RenderMarkdown(msg, true, true)
	cli.SendMessageEvent(evt.RoomID, event.EventMessage, &htmlMsg)
	return true
}

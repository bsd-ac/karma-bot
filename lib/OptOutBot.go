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
	"regexp"

	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/event"
)

type OptOutBot struct {
}

func (u *OptOutBot) MatchMessage(body string) bool {
	rexp := regexp.MustCompile(`(?i)^\!optout\s*$`)
	return rexp.MatchString(body)
}

func (u *OptOutBot) ProcessMessage(body string, source mautrix.EventSource, evt *event.Event, kBot *KarmaBot) {
	senderID := evt.Sender.String()
	kBot.OptOut(senderID)
}

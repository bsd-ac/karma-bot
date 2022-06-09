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

	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/format"
)

type Command_OptStatus struct {
}

func (u *Command_OptStatus) NeedsTimer() bool {
	return true
}

func (u *Command_OptStatus) Process(evt *event.Event, kBot *KarmaBot, targetID, targetHREF string) bool {
	optOut := kBot.IsOptOut(targetID)
	msg := ""
	if optOut {
		msg = fmt.Sprintf("%s is not allowed to be tracked in the karma system", targetHREF)
	} else {
		msg = fmt.Sprintf("%s can be tracked in the karma system", targetHREF)
	}
	msgHTML := format.RenderMarkdown(msg, true, true)
	kBot.mClient.SendMessageEvent(evt.RoomID, event.EventMessage, &msgHTML)
	return true
}

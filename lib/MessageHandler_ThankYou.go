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

	"maunium.net/go/mautrix/event"
)

type MessageHandler_ThankYou struct {
}

func (u *MessageHandler_ThankYou) NeedsTimer() bool {
	return false
}

func (u *MessageHandler_ThankYou) FastMatch(body, bodyHTML string) bool {
	return bodyHTML != "" && regexp.MustCompile(`(?i)\bthank(s)?\b`).MatchString(body)
}

func (u *MessageHandler_ThankYou) Re() []*regexp.Regexp {
	return []*regexp.Regexp{
		regexp.MustCompile(`(?i)(thanks(\s+a\s+(lot|bunch))?|thank\s+you(\s+very\s+much)?)\s+(<a\s+href=".+/(.+)">.+</a>)`),
		regexp.MustCompile(`(?i)(<a\s+href=".+/(.+)">.+</a>)(\s*:\s*)?\s*(thanks(\s+a\s+(lot|bunch))?|thank\s+you(\s+very\s+much)?|\++)`),
	}
}

func (u *MessageHandler_ThankYou) ReIndex() []int {
	return []int{6, 2}
}

func (u *MessageHandler_ThankYou) MatchMessage(body string) bool {
	rexp_arr := u.Re()
	for _, rexp := range rexp_arr {
		if rexp.MatchString(body) {
			return true
		}
	}
	return false
}

func (u *MessageHandler_ThankYou) ProcessMessage(evt *event.Event, kBot *KarmaBot, body, bodyHTML string) bool {
	if bodyHTML == "" {
		return false
	}
	kBot.logger.Infof("Called MessageHandler_ThankYou")
	rexp_arr := u.Re()
	rind_arr := u.ReIndex()
	senderID := evt.Sender.String()
	for i, rexp := range rexp_arr {
		groups := rexp.FindAllStringSubmatch(bodyHTML, -1)
		rind := rind_arr[i]
		found := false
		for j := 0; j < len(groups); j++ {
			targetID := groups[j][rind]
			if targetID == "" {
				continue
			} else {
				kBot.KarmaAdd(senderID, targetID, evt.ID.String(), evt.RoomID.String(), 1)
				found = true
				break
			}
		}
		if found {
			break
		}
	}
	return false
}

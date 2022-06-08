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
	"strings"

	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/event"
)

type ThankYouBot struct {
}

func (u *ThankYouBot) ID() string {
	return "ThankYouBot"
}

func (u *ThankYouBot) NeedsTimer() bool {
	return false
}

func (u *ThankYouBot) Re() []*regexp.Regexp {
	return []*regexp.Regexp{
		regexp.MustCompile(`(?i)(thanks(\s+a\s+(lot|bunch))?|thank\s+you(\s+very\s+much)?)\s+(<a\s+href=".+/(.+)">.+</a>)`),
		regexp.MustCompile(`(?i)(<a\s+href=".+/(.+)">.+</a>)(\s*:\s*)?\s*(thanks(\s+a\s+(lot|bunch))?|thank\s+you(\s+very\s+much)?)`),
	}
}

func (u *ThankYouBot) ReIndex() []int {
	return []int{6, 2}
}

func (u *ThankYouBot) MatchMessage(body string) bool {
	rexp_arr := u.Re()
	for _, rexp := range rexp_arr {
		if rexp.MatchString(body) {
			return true
		}
	}
	return false
}

func (u *ThankYouBot) ProcessMessage(body string, source mautrix.EventSource, evt *event.Event, kBot *KarmaBot) bool {
	htmlBody := strings.TrimSpace(evt.Content.AsMessage().FormattedBody)
	if htmlBody == "" {
		htmlBody = body
	}
	if !u.MatchMessage(htmlBody){
		return false
	}
	kBot.logger.Infof("Called ThankYouBot")
	rexp_arr := u.Re()
	rind_arr := u.ReIndex()
	senderID := evt.Sender.String()
	for i, rexp := range rexp_arr {
		groups := rexp.FindAllStringSubmatch(htmlBody, -1)
		rind := rind_arr[i]
		found := false
		for j := 0; j < len(groups); j++ {
			targetID := groups[j][rind]
			if targetID == "" {
				continue
			} else {
				kBot.KarmaAdd(senderID, targetID, evt.ID.String(), evt.RoomID.String(), 1)
				found = true
			}
		}
		if found {
			break
		}
	}
	return false
}

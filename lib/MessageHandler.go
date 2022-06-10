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
	"time"

	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/event"
)

var RoomTimers = make(map[string]int64)

type KarmaCommand interface {
	NeedsTimer() bool
	Process(evt *event.Event, kBot *KarmaBot, targetID, targetHREF string) bool
}

var KarmaCommands = map[string]KarmaCommand{
	"karma":     &Command_Karma{},
	"tkarma":    &Command_KarmaTotal{},
	"optin":     &Command_OptIn{},
	"optout":    &Command_OptOut{},
	"optstatus": &Command_OptStatus{},
	"uptime":    &Command_Uptime{},
}

type KarmaMessageHandler interface {
	NeedsTimer() bool
	FastMatch(body, bodyHTML string) bool
	ProcessMessage(evt *event.Event, kBot *KarmaBot, body, bodyHTML string) bool
}

var KarmaMessageHandlers = []KarmaMessageHandler{
	&MessageHandler_ThankYou{},
}

func MessageHandler(source mautrix.EventSource, evt *event.Event, kBot *KarmaBot) {
	if evt.Sender == kBot.WhoAmI() {
		return
	}
	body := strings.TrimSpace(evt.Content.AsMessage().Body)
	bodyHTML := strings.TrimSpace(evt.Content.AsMessage().FormattedBody)
	if bodyHTML == "" {
		bodyHTML = body
	}
	tnow := time.Now().UnixMicro()
	roomID := evt.RoomID.String()
	if _, ok := RoomTimers[roomID]; !ok {
		RoomTimers[roomID] = 0
	}
	if body[0] == '!' {
		rexp := regexp.MustCompile(`(?i)^\!([a-z]+)(\s+.*)?$`)
		groups := rexp.FindAllStringSubmatch(bodyHTML, -1)
		if groups == nil || len(groups) != 1 || len(groups[0]) == 0 {
			kBot.logger.Warnf("Could not parse command: %s", bodyHTML)
			return
		}
		commandName := groups[0][1]
		if command, ok := KarmaCommands[commandName]; ok {
			href := strings.TrimSpace(groups[0][2])
			targetID := HTMLToUserID(href)
			if targetID == "" {
				targetID = evt.Sender.String()
				href = targetID
			}
			go func() {
				roomID = evt.RoomID.String()
				if !command.NeedsTimer() || tnow-RoomTimers[roomID] > kBot.kConf.ResponseFreq {
					if command.Process(evt, kBot, targetID, href) {
						RoomTimers[roomID] = tnow
					}
				}
			}()
		}
	} else {
		for _, handler := range KarmaMessageHandlers {
			if !handler.NeedsTimer() || (tnow-RoomTimers[roomID] > kBot.kConf.ResponseFreq && handler.FastMatch(body, bodyHTML)) {
				if handler.ProcessMessage(evt, kBot, body, bodyHTML) {
					RoomTimers[roomID] = tnow
				}
			}
		}
	}
}

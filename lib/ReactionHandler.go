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
	"strings"

	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/event"
)

func ReactionHandler(source mautrix.EventSource, evt *event.Event, kBot *KarmaBot) {
	relatesTo := evt.Content.AsReaction().GetRelatesTo()
	emoji := relatesTo.GetAnnotationKey()
	senderID := evt.Sender.String()
	targetEvent, err := kBot.mClient.GetEvent(evt.RoomID, relatesTo.EventID)
	if err != nil {
		kBot.logger.Warnf("Error while retrieving target event: %v", err)
		return
	}
	targetID := targetEvent.Sender.String()
	for _, pemoji := range strings.Split(kBot.kConf.PositiveEmojis, ",") {
		if emoji == pemoji {
			kBot.KarmaAdd(senderID, targetID, evt.ID.String(), evt.RoomID.String(), 1)
			return
		}
	}
	for _, nemoji := range strings.Split(kBot.kConf.NegativeEmojis, ",") {
		if emoji == nemoji {
			kBot.KarmaAdd(senderID, targetID, evt.ID.String(), evt.RoomID.String(), -1)
			return
		}
	}
}

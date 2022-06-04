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
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/event"
)

func MessageHandler(source mautrix.EventSource, evt *event.Event, kBot *KarmaBot) {

	body := evt.Content.AsMessage().Body
	for _, plugin := range Plugins {
		if plugin.MatchMessage(body) {
			plugin.ProcessMessage(body, source, evt, kBot)
		}
	}

	/*
		for _, pat := range regex_arr {
			re := regexp.MustCompile(pat)
			if re.MatchString(body) {
				person := re.ReplaceAllString(body, "$1")
				cli.Logger.Debugfln("person = %v\n", person)
			} else {
				cli.Logger.Debugfln("no person found\n")
			}
		}
	*/
}

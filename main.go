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

package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/id"
	"suah.dev/protect"

	"bsd.ac/karma-bot/lib"
)

var config = flag.String("f", "/etc/karma-bot.ini", "alternative configuration file")

func main() {
	flag.Parse()
	if *config == "" {
		_, _ = fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}
	cfg, _ := lib.ReadConfig(*config)

	fmt.Printf("Logging into %s as %s\n", cfg.Homeserver, cfg.Username)
	client, err := mautrix.NewClient(cfg.Homeserver, id.UserID(cfg.Username), cfg.AccessToken)
	if err != nil {
		panic(err)
	}

	protect.Pledge("stdio unveil rpath wpath cpath flock dns inet tty")
	protect.Unveil("/etc/resolv.conf", "r")
	protect.Unveil("/etc/ssl/cert.pem", "r")
	protect.Unveil(filepath.Dir(cfg.Database), "rwx")
	protect.UnveilBlock()

	syncer := client.Syncer.(*mautrix.DefaultSyncer)
	syncer.OnEventType(event.EventMessage, func(source mautrix.EventSource, evt *event.Event) { lib.MessageHandler(client, source, evt) })
	syncer.OnEventType(event.EventReaction, func(source mautrix.EventSource, evt *event.Event) { lib.ReactionHandler(client, source, evt) })
	syncer.OnEventType(event.EventRedaction, func(source mautrix.EventSource, evt *event.Event) { lib.RedactionHandler(client, source, evt) })

	err = client.Sync()
	if err != nil {
		panic(err)
	}
}

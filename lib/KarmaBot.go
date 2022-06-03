/*
 * Copyright (c) 2022 Aisha Tammy <aisha@bsd.ac>
 * Copyright (c) 2021 Aaron Bieber <aaron@bolddaemon.com>
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
	"path/filepath"

	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/id"
)

type KarmaBot struct {
	kConf   *KarmaConfig
	logger  *BotLogger
	mClient *mautrix.Client
	bDB     *BDBStore
	sqlDB   *SQLStore
}

func NewKarmaBot(kConf *KarmaConfig) *KarmaBot {
	kBot := new(KarmaBot)
	kBot.kConf = kConf
	kBot.logger = NewBotLogger()
	return kBot
}

func (kBot *KarmaBot) Start() error {
	var err error
	kBot.bDB, err = NewBDBStore(filepath.Join(kBot.kConf.DBDirectory, "badger"), kBot.logger)
	if err != nil {
		return err
	}

	kBot.sqlDB, err = NewSQLStore(kBot.kConf.DBtype, kBot.kConf.DBdsn)
	if err != nil {
		kBot.bDB.Close()
		return err
	}

	kBot.mClient, err = mautrix.NewClient(kBot.kConf.Homeserver, id.UserID(kBot.kConf.Username), kBot.kConf.AccessToken)
	if err != nil {
		kBot.bDB.Close()
		kBot.sqlDB.Close()
		return err
	}
	kBot.mClient.Store = kBot.bDB
	kBot.mClient.Logger = kBot.logger

	syncer := kBot.mClient.Syncer.(*mautrix.DefaultSyncer)
	syncer.OnEventType(event.EventMessage, func(source mautrix.EventSource, evt *event.Event) {
		MessageHandler(kBot.mClient, source, evt, kBot.bDB, kBot.sqlDB)
	})
	/*
	syncer.OnEventType(event.EventReaction, func(source mautrix.EventSource, evt *event.Event) {
		ReactionHandler(client, source, evt, kBot.bDB, kBot.sqlDB)
	})
	syncer.OnEventType(event.EventRedaction, func(source mautrix.EventSource, evt *event.Event) {
		RedactionHandler(client, source, evt, kBot.bDB, kBot.sqlDB)
	})
	*/
	err = kBot.mClient.Sync()

	if (err != nil) {
		kBot.bDB.Close()
		kBot.sqlDB.Close()
	}
	return err
}

func (kBot *KarmaBot) Stop() {
	kBot.mClient.StopSync()
	kBot.bDB.Close()
	kBot.sqlDB.Close()
}
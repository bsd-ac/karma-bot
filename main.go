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
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/id"
	"suah.dev/protect"

	"bsd.ac/karma-bot/lib"
)

func main() {
	flag := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	debugLevel := flag.String("d", "warn", "debug level of output (debug, info, warn, error, dpanic, panic, fatal)")
	config := flag.String("f", "/etc/karma-bot.ini", "alternative configuration file")
	outputFormat := flag.String("o", "console", "debug output format (console, json)")
	flag.Parse(os.Args[1:])

	zconf := zap.NewProductionConfig()
	zconf.Encoding = *outputFormat
	zlevel, err := zapcore.ParseLevel(*debugLevel)
	if err != nil {
		log.Fatalf("ERROR: could not set debug level: %v", err)
	}
	zconf.Level = zap.NewAtomicLevelAt(zlevel)
	zlog, err := zconf.Build()
	if err != nil {
		log.Fatalf("ERROR: could not initialize logger: %v", err)
	}
	zap.ReplaceGlobals(zlog)
	klog := zlog.Sugar()
	defer klog.Sync()
	// only use klog as logger from here on

	klog.Debugf("Reading config file '%s'", *config)
	cfg, err := lib.ReadConfig(*config)
	if cfg == nil {
		klog.Fatalf("Could not read the config file '%s': %v", *config, err)
	}
	klog.Debugf("Finished reading config file")

	klog.Debugf("Securing with pledge and unveil")
	protect.Pledge("stdio unveil rpath wpath cpath flock dns inet tty")
	protect.Unveil("/etc/resolv.conf", "r")
	protect.Unveil("/etc/ssl/cert.pem", "r")
	protect.Unveil(cfg.DBDirectory, "rwxc")
	protect.UnveilBlock()
	klog.Debugf("Finished securing")

	bdbStore, err := lib.NewBDBStore(filepath.Join(cfg.DBDirectory, "badger"))
	if err != nil {
		klog.Fatalf("Could not open the BDBStore: %v", err)
	}
	defer bdbStore.BDB.Close()

	sqlStore, err := lib.NewSQLStore(cfg.DBtype, cfg.DBdsn)
	if err != nil {
		bdbStore.BDB.Close()
		klog.Fatalf("Could not open the SQLStore: %v", err)
	}
	defer sqlStore.DB.Close()

	klog.Infof("Creating client for %s with username %s", cfg.Homeserver, cfg.Username)
	client, err := mautrix.NewClient(cfg.Homeserver, id.UserID(cfg.Username), cfg.AccessToken)
	if err != nil {
		panic(err)
	}
	client.Store = bdbStore
	mautrixLogger := new(lib.MautrixLogger)
	mautrixLogger.Logger = klog
	client.Logger = mautrixLogger

	klog.Debugf("Creating bot client")
	syncer := client.Syncer.(*mautrix.DefaultSyncer)
	klog.Debugf("Adding even handlers")
	syncer.OnEventType(event.EventMessage, func(source mautrix.EventSource, evt *event.Event) { lib.MessageHandler(client, source, evt) })
	syncer.OnEventType(event.EventReaction, func(source mautrix.EventSource, evt *event.Event) { lib.ReactionHandler(client, source, evt) })
	syncer.OnEventType(event.EventRedaction, func(source mautrix.EventSource, evt *event.Event) { lib.RedactionHandler(client, source, evt) })

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		err = client.Sync()
		if err != nil {
			panic(err)
		}
	}()

	sig := <-done
	klog.Infof("Caught signal '%v'", sig)
	klog.Infof("Stopping the client...")
	client.StopSync()
	klog.Infof("Stopped")
	klog.Infof("Shutting down")
}

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
	"syscall"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"suah.dev/protect"

	"bsd.ac/karma-bot/lib"
)

func main() {
	flag := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	debugLevel := flag.String("d", "error", "debug level of output (debug, info, warn, error, dpanic, panic, fatal)")
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
	kConf, err := lib.ReadConfig(*config)
	if err != nil {
		klog.Fatalf("Error while reading the config file: %s", err.Error())
	}
	klog.Debugf("Finished reading config file")

	klog.Debugf("Securing with pledge and unveil")
	protect.Unveil("/etc/resolv.conf", "r")
	protect.Unveil("/etc/ssl/cert.pem", "r")
	protect.Unveil(kConf.DBDirectory, "rwxc")
	protect.UnveilBlock()
	protect.Pledge("stdio rpath wpath cpath flock dns inet tty")
	klog.Debugf("Finished securing")

	kbot := lib.NewKarmaBot(kConf)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		err := kbot.Start()
		if err != nil {
			klog.Fatalf("Could not start the bot: %v", err)
		}
	}()

	sig := <-done
	klog.Infof("Caught signal '%v'", sig)
	klog.Infof("Shutting down...")
	kbot.Stop()
}

[![build](https://github.com/bsd-ac/karma-bot/actions/workflows/build.yml/badge.svg)](https://github.com/bsd-ac/karma-bot/actions/workflows/build.yml)
[![CodeQL](https://github.com/bsd-ac/karma-bot/actions/workflows/codeql-analysis.yml/badge.svg)](https://github.com/bsd-ac/karma-bot/actions/workflows/codeql-analysis.yml)
[![Maintainability](https://img.shields.io/codeclimate/maintainability/bsd-ac/karma-bot.svg)](https://codeclimate.com/github/bsd-ac/karma-bot)
[![GitHub license](https://img.shields.io/github/license/bsd-ac/karma-bot.svg)](https://github.com/bsd-ac/karma-bot/blob/master/LICENSE)
[![GitHub issues](https://img.shields.io/github/issues-raw/bsd-ac/karma-bot)](https://github.com/bsd-ac/karma-bot/issues)
[![contributions welcome](https://img.shields.io/badge/contributions-welcome-brightgreen.svg?style=flat)](https://github.com/bsd-ac/karma-bot/issues)

# Karma Bot

Karma bot for [Matrix](https://matrix.org/), invite [Banana-Bot:matrix.org](https://matrix.to/#/@banana-bot:matrix.org) to the channel.

## Features

On top of standard features which parse messages of the form `thanks <abcxyz>`, this bot also supports a few more enhancements

- Support message reactions:
  - positive emojis give karma: ❤️,👍️,💯,🍌,🎉,💞,💗,💓,💖,💘,💝,💕,😻,😍,❤️‍🔥
  - negative emojis reduce karma: 👎️,💔,😠,👿,🙁,☹️,🤬,☠️,💀
  - removing the reactions removes the karma contribution
- Per room and global karma stats.
- Ability to opt out/in of tracking: `!optout`, `!optin`

## Commands

| command             | notes                                                                                                                   |
|---------------------|-------------------------------------------------------------------------------------------------------------------------|
| `!karma [user]`     | get karma of a user in a room, defaults to sender if user is not specified                                              |
| `!tkarma [user]`    | get karma of a user across all rooms, defaults to sender if user is not specified                                       |
| `!optin`            | remove the sender from the karma tracking system (all votes give to and by the sender are deleted and permanently lost) |
| `!optout`           | allow the sender to be tracked in the karma tracking system (past events are not tracked)                               |
| `!optstatus [user]` | check if a user has opted in/out of the karma tracking system                                                           |
| `!uptime`           | check how long the bot has been up                                                                                      |

## Usage

```
$ karma-bot -h
Usage of karma-bot:
  -d string
        debug level of output (debug, info, warn, error, dpanic, panic, fatal) (default "warn")
  -f string
        alternative configuration file (default "/etc/karma-bot.ini")
  -o string
        debug output format (console, json) (default "console")
```

The [sample config file](karma-bot.ini.sample) contains detailed explanations of options to configure.

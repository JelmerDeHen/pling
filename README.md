# Pling

Play sound and lock screen when user is afk.

## Config

Looks for yaml config file named `.pling` in `$HOME` and `pwd`. The following options exist:

|Name|Type|Default|Description|
| --- | --- | --- | --- |
| `afk_timeout` | `time.Duration` | `10m` | Duration until user is afk |
| `i3lock` | `bool` | `false` | Exec i3lock when `afk_timeout` is reached |
| `i3lock_color` | `string` | `000000` | Change i3lock screen color to rrggbb |
| `mp3_file` | `string` | | mp3 file played when user is afk |
| `mp3_hour_start` | `int` |  | Hour of day to start playing mp3 |
| `mp3_hour_stop` | `int` | | Hour of day to stop playing mp3 |
| `mp3_interval` | `time.Duration` | `5s` | Duration to wait until playing mp3 file again |
| `dsn` | `string` | | sqlite3 db path |

Environment variables can be used to overwrite the configuration. The environment variable name for a config value is the name uppercased and prefixed with `PLING_`.

Flags can also be used to overwrite environment and config values, use `--help` to see the available flags.

## Install

Configuration file and systemd unit files can be deployed using `make install`, `make uninstall` reverts this process.

## Sqlite3

The database created by configuring `dsn` contains a table `activity` where presence is logged. The database can be inspected by using --list:

```
pling --dsn=$HOME/pling.db --list
```

`activity` table:

```
sqlite> pragma table_info('activity');
cid  name   type     notnull  dflt_value  pk
---  -----  -------  -------  ----------  --
0    id     INTEGER  1                    1
1    state  TEXT     0                    0
2    start  DATE     0                    0
3    stop   DATE     0                    0
```


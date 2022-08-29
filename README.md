# Pling

Play sound and lock screen when user is afk.

## Config

Looks for config file named `.pling` in `$HOME` and current working directory. Configuration file format is yaml.

Configuration file named `.pling` are searched in `$HOME` or current working directory. The following options exist:


|Name|Type|Default|Description|
| --- | --- | --- | --- |
| `afk_timeout` | `time.Duration` | `10m` | Time until user is considered afk and mp3 is played and i3lock is executed |
| `i3lock` | `bool` | `false` | When true executes i3lock when `afk_timeout` is reached |
| `i3lock_color` | `string` | `000000` | Configures i3lock screen in rrggbb color instead of black |
| `mp3_file` | `string` | in-memory file | Configure the mp3 file played when user is afk |
| `mp3_hour_start` | `int` | ``  | Hour of day to start playing mp3 |
| `mp3_hour_stop` | `int` | `` | Hour of day to stop playing mp3 |
| `mp3_interval` | `time.Duration` | `5s` | Duration to wait until playing mp3 file again |

See `.pling` file for default config file.

Example config file:

```yaml
afk_timeout: 10m
mp3_file: /some/path.mp3
mp3_interval: 1m
i3lock: false
i3lock_color: 00ff00
```

Environment variables can be used to overwrite config values. The environment variable name for a config value is the name uppercased and prefixed with `PLING_`.

Example of overwriting `afk_timeout` config value using environment variable `PLING_AFK_TIMEOUT`:

```sh
make && PLING_AFK_TIMEOUT=3s ./bin/pling
```

## Install

Configuration file and systemd unit files can be deployed using `make install`, `make uninstall` reverts this process.



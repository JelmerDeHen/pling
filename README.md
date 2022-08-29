# Pling

Play sound and lock screen when user is afk.

## Config

Looks for config file named `.pling` in `$HOME` and current working directory. Configuration file format is yaml.

Configuration file named `.pling` are searched in `$HOME` or current working directory. The following options exist:

- `afk_timeout` determines afk duration until afk handler is executed
- `i3lock` setting this to true will execute i3lock when user is afk
- `i3lock_color` turn i3lock screen in rrggbb color instead of black
- `mp3` not implemented yet, will allow changing mp3 played
- `mp3_interval` interval between playing mp3 when user is afk
- `mp3_hour_start` hour of day to start playing mp3
- `mp3_hour_stop` hour of day to stop playing mp3

Example config file:

```yaml
afk_timeout: 10m
mp3: /some/path.mp3
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



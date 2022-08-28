# Pling

Play sound and lock screen when user is afk.

## Config

Looks for config file named `.pling` in `$HOME` and current working directory. Configuration file format is yaml.

```yaml
idle_over_timeout: 10m
poll_interval: 1m
i3lock: false
i3lock_color: 00ff00
mp3: /some/path.mp3
```

Environment variables can be used to overwrite configuration file like this:

```
make && PLING_IDLE_OVER_TIMEOUT=3s PLING_POLL_INTERVAL=1s ./bin/pling
```

- `idle_over_timeout` determines afk duration until afk handler is executed
- `poll_interval` will affect the interval by which the mp3 sound is played
- `i3lock` setting this to true will execute i3lock when user is afk
- `i3lock_color` what rrggbb color i3lock screen will be
- `mp3` not implemented yet, will allow changing mp3 played
- `mp3_hour_start` configures what hour to start playing mp3 sound
- `mp3_hour_stop` configures what hour to stop playing mp3 sound

## Install

Configuration file and systemd unit files can be deployed using `make install`, `make uninstall` reverts this process.



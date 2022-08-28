# Pling

Play sound when user is afk

## Config

Looks for config file named `.pling` in `$HOME` and `.`. Configuration file format is yaml.

```yaml
idle_over_timeout: 10m
poll_interval: 1m
```

Environment variables can be used to overwrite configuration file like this:

```
make && PLING_IDLE_OVER_TIMEOUT=3s PLING_POLL_INTERVAL=1s ./bin/pling
```


## Install

Configuration file and systemd unit files can be deployed using `make install`, `make uninstall` reverts this process.

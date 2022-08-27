
ROOT_DIR:=$(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))
BIN:=${ROOT_DIR}/bin/pling

build:
				mkdir -pv bin
				go build -o bin/pling

clean:
				rm -rf bin

install:
				echo "${BIN}"
				printf '[Service]\nExecStart=${BIN}\n' > "${HOME}/.config/systemd/user/pling.service"
				printf '[Timer]\nOnBootSec=0min\nOnUnitActiveSec=1min\n\n[Install]\nWantedBy=timers.target' > "${HOME}/.config/systemd/user/pling.timer"
				systemctl --user daemon-reload
				systemctl --user enable pling.timer

uninstall:
				systemctl --user stop pling.service pling.timer
				systemctl --user disable pling.timer
				rm -v "${HOME}/.config/systemd/pling.service" "${HOME}/.config/systemd/pling.timer"
				systemctl --user daemon-reload

start:
				systemctl --user start pling.timer pling.service

stop:
				systemctl --user stop pling.timer pling.service

status:
				systemctl --user status pling.timer pling.service

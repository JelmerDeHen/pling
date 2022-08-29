
ROOT_DIR:=$(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))
BIN:=${ROOT_DIR}/bin/pling
SYSTEMD_USER_DIR:=${HOME}/.config/systemd/user

build:
				mkdir -pv bin
				go build -o bin/pling

clean:
				rm -rf bin

install:
				cp -v .pling "${HOME}/.pling"
				printf '[Service]\nExecStart=${BIN}\n' > "${SYSTEMD_USER_DIR}/pling.service"
				printf '[Timer]\nOnBootSec=0min\nOnUnitActiveSec=1min\n\n[Install]\nWantedBy=timers.target' > "${SYSTEMD_USER_DIR}/pling.timer"
				systemctl --user daemon-reload
				systemctl --user enable pling.timer

uninstall:
				systemctl --user stop pling.service pling.timer
				systemctl --user disable pling.timer
				rm -v "${SYSTEMD_USER_DIR}/pling.service" "${SYSTEMD_USER_DIR}/pling.timer" "${HOME}/.pling"
				systemctl --user daemon-reload

start:
				systemctl --user start pling.timer pling.service

stop:
				systemctl --user stop pling.timer pling.service

status:
				systemctl --user status pling.timer pling.service

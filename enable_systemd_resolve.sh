# !/bin/env sh

sudo systemctl start systemd-resolved
sudo systemctl enabled systemd-resolved
sudo systemctl unmask systemd-resolved
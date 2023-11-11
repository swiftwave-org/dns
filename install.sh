# !/bin/env sh

# Move binary to /usr/bin
sudo mv ./swiftwave-dns /usr/bin/

# Move service file to /etc/systemd/system/
sudo mv ./swiftwave-dns.service /etc/systemd/system/

# Daemon reload and enable
sudo systemctl daemon-reload
sudo systemctl enable swiftwave-dns.service
sudo systemctl start swiftwave-dns.service


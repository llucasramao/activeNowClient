#!/bin/bash

echo "[+] Iniciando instalação do agente"

mkdir -p /opt/activeNow
rm -R /opt/activeNow/*

curl -o /opt/activeNow/main http://192.168.1.31/main
chmod +x /opt/activeNow/main

echo "[Unit]
Description=ActiveNow BBTS Agent
After=network.target

[Service]
Type=simple
User=root
ExecStart=/opt/activeNow/main
Restart=on-failure

[Install]
WantedBy=multi-user.target" > /etc/systemd/system/activeNow.service

systemctl daemon-reload
systemctl restart activeNow
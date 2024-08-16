
### SystemD Services

`/etc/systemd/system/smartlabels.service`
```
[Unit]
Description=Smartlabels2
After=network-online.target
Wants=network-online.target systemd-networkd-wait-online.service

StartLimitIntervalSec=500

[Service]
Restart=on-failure
RestartSec=5s

WorkingDirectory=/home/install/smartlabels2
ExecStart=/home/install/smartlabels2/go_build_for_Linux_linux

[Install]
WantedBy=multi-user.target
```


`/etc/systemd/system/smartlabels2.service`
```
[Unit]
Description=Smartlabels2
After=network-online.target
Wants=network-online.target systemd-networkd-wait-online.service

StartLimitIntervalSec=500

[Service]
Restart=on-failure
RestartSec=5s

WorkingDirectory=/home/install/smartlabels2
ExecStart=/home/install/smartlabels2/go_build_for_Linux_linux

[Install]
WantedBy=multi-user.target
```

### Logrotate Config
`/etc/logrotate.d/smartlabels`
```
/var/log/smartlabels/*log {
        daily
        missingok
        rotate 3
        notifempty
        dateext
        dateformat -%Y-%m-%d
        create 640 root root
        compress
}
```
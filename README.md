# btfsautohostprice
BTFS auto host price automatically set BTFS host BTT price from coinmarketcap and a fixed USD amount everyday

# build

Simply `go build -ldflags "-s -w"` the project.

# Install

* Copy the `btfsautohostprice` binary from build to `/usr/local/bin/btfsautohostprice`
* Copy `systemd/btfsautohostprice.service` to `/etc/systemd/system/btfsautohostprice.service`
    * Make sure your btfs daemon is started by a unit called `btfs.service` or edit `After=` and `Requires=` in `btfsautohostprice.service`
* Copy `systemd/btfsautohostprice.timer` to `/etc/systemd/system/btfsautohostprice.timer`
* Copy `systemd/btfsautohostprice.env` to `/etc/default/btfsautohostprice`
* Edit `/etc/default/btfsautohostprice` with valid values
* Set perms with:

```bash
useradd --home-dir /var/lib/btfsautohostprice --create-home --system --shell /usr/sbin/nologin btfsautohostprice
chown root:btfsautohostprice /usr/local/bin/btfsautohostprice /etc/default/btfsautohostprice
chmod 750 /usr/local/bin/btfsautohostprice
chmod 640 /etc/default/btfsautohostprice
```

* Activate systemd with:

```bash
systemctl daemon-reload
systemctl enable --now btfsautohostprice.timer
```

* Check status with:

```bash
systemctl list-timers
systemctl status btfsautohostprice.timer
journalctl -u btfsautohostprice.service
```

* Force run with:

```bash
systemctl start btfsautohostprice.service
```
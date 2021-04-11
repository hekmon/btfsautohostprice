# btfsautohostprice
BTFS auto host price automatically set BTFS host BTT price from coinmarketcap and a fixed USD amount everyday

## build

Simply `go build -ldflags "-s -w"` the project.

## install

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

## run example

`journalctl -u btfsautohostprice.service`
```
Apr 11 08:33:57 serverhostname systemd[1]: Starting BTFS automatic host price update...
Apr 11 08:33:58 serverhostname btfsautohostprice[123483]: 10.00 USD is worth 1233.937434 BTT at 2021-04-11 08:33:07 +0000 UTC: with the 3x network redundancy, a host price must be 411.312478 BTT for a user to store 1TB/month on the network at this price
Apr 11 08:33:58 serverhostname btfsautohostprice[123483]: host pricing updated
Apr 11 08:33:58 serverhostname systemd[1]: btfsautohostprice.service: Succeeded.
Apr 11 08:33:58 serverhostname systemd[1]: Finished BTFS automatic host price update.
```
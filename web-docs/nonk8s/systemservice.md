# Run flagD as a systemd Service

To install as a systemd service clone the repo and run `sudo make install`.
This will place the binary by default in `/usr/local/bin`.

There will also be a default provider and sync enabled ( http / filepath ) both of which can be modified in the flagd.service.

Validation can be run with `systemctl status flagd`
And result similar to below will be seen

```console
  flagd.service - "A generic feature flag daemon"
     Loaded: loaded (/etc/systemd/system/flagd.service; disabled; vendor preset: enabled)
     Active: active (running) since Mon 2022-05-30 12:19:55 BST; 5min ago
   Main PID: 64610 (flagd)
      Tasks: 7 (limit: 4572)
     Memory: 1.4M
     CGroup: /system.slice/flagd.service
             └─64610 /usr/local/bin/flagd start -f=/etc/flagd/flags.json

May 30 12:19:55 foo systemd[1]: Started "A generic feature flag daemon".
```

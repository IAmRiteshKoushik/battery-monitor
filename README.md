# Battery Monitor and Notifier in Go

# Motivations
This stems from a personal requirement to trigger the notification
manager inside of linux to give me a heads-up at certain battery percentages.

Ideally I want the following"
1. At 20% a notification to start charging my device
2. At 90% a notification to stop charging my device (if plugged-in)
3. At 95% a warning that battery might be getting damaged (if plugged-in)

This is being written with consideration that a laptop battery undergoes
degradation when overcharged. This is especially true with batteries which
are OEM (have been replaced).

# About the Project

This is a monitoring tool which runs every 60 seconds and check for battery
details and releases based on the above mentioned constraints.

For error logs check
```bash
journalctl -u battery-monitor.service
```

# How to Run As a Background Process
1. Clone the repository
```bash
git clone github.com/IAmRiteshKoushik/battery-monitor
cd battery-monitor
go build .
```
2. Set the `ExecStart` path properly inside `battery-monitor.service`.
3. Set the `Environment` variables by using the echo command on `$DISPLAY` and 
`$XAUTHORITY`
4. Run using `systemd` (System and Service manager for Linux operating systems)
```bash
sudo cp battery-monitor.service /etc/systemd/user/battery-monitor.service
sudo systemctl --user enable battery-monitor.service
sudo systemctl --user start battery-monitor.service
```
Here, we are making the process user-specific. However, if you wish to configure
it globally for all users then you would have to look into setting up 
environment variables like `DBUS_SESSION_BUS_ADDRESS`. This is needed to 
trigger `notify-send` from a root-process to a user-session. You can read more
about it [here](https://bbs.archlinux.org/viewtopic.php?id=205867).

# Some Considerations

The battery module differs from workstation to workstation. It is best to find
out what your battery module is and where the data is located and appropriately
setup that data inside the following variables in the source code:
```go
	current_enery_file := "/sys/class/power_supply/BAT0/energy_now"
	energy_full_file := "/sys/class/power_supply/BAT0/energy_full"
	charging_status_file := "/sys/class/power_supply/BAT0/status"
```

It is usually found in `sys/class/power_supply`

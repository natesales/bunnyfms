# BunnyFMS

BunnyFMS is a lightweight field management system (FMS) for FRC robots.

## Hardware

In addition to a computer, BunnyFMS requires a wireless access point and an (unmanaged) network switch. An easy solution is to use a CPE-style router that incorporates an AP and switch into one device.

The FMS computer needs a connection to the field network and optionally an audio output device for game sounds.

### Event Configuration

1. Field access point
    1. IP address: `10.0.100.1`
    2. Subnet Mask: `255.255.255.0`
    3. 5 Ghz wireless:
        1. Mode: `N-only`
        2. SSID: `<event name>`
        3. Channel spacing: `20MHz`
        4. SSID Broadcast: `Disabled`
        5. Security: `WPA2-Personal + AES`
        6. Password: `<event password>`
        7. Key renewal: `3600`

2. Network switch
   1. All drive stations need a connection to the field network, so you may need a network switch for the 6 DS connections
   2. Configure all switch ports in the same VLAN (BunnyFMS does not use the same VLAN setup as an official FRC FMS) 

3. FMS computer
    1. Static IP address: `10.0.100.5/8`
    2. Default gateway and DNS if required: `10.0.100.1`
    4. Run (`bunnyfms -admin :8080 -viewer :8081 -auto-duration 10s -teleop-duration 2m20s -endgame-duration 30s`)

4. Robot radio kiosk
    1. Install the [FRC Radio Configuration Utility](https://docs.wpilib.org/en/stable/docs/zero-to-robot/step-3/radio-programming.html)
    2. Choose `Tools > FMS-Lite/Offseason FMS Mode`
    3. Enter your event name and password

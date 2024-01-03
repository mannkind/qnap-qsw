# qnap-qsw

A quick and dirty script to take action on a QNAP QSW-M2116P-2T2S.

## Commands

### Login
```
token=$(./qnap-qsw login --host switch.lan --password BLAH)
```

### Disable POE Ports
#### With Password
```
./qnap-qsw poeMode --host switch.lan --password BLAH --ports 1,2,3,4 --mode disable
```
#### With Token
```
./qnap-qsw poeMode --host switch.lan --token BLAH --ports 1,2,3,4 --mode disable
```

### Enable POE Ports
#### With Password
```
./qnap-qsw poeMode --host switch.lan --password BLAH --ports 1,2,3,4 --mode poePlusDot3at
```
#### With Token
```
./qnap-qsw poeMode --host switch.lan --token BLAH --ports 1,2,3,4 --mode poePlusDot3at
./qnap-qsw poeMode --host switch.lan --token BLAH --ports 19,20 --mode poePlusDot3bt
```

### Home Assistant

#### shell_commands.yaml
```
qnap_qsw_poe: "/config/qnap-qsw poeMode --host {{ host }} --password {{ password }} --ports {{ ports }} --mode {{ mode }}"
```

#### automations.yaml
```
alias: "Turn off non-essebtial POE ports during power outage"
description: ""
trigger:
  - type: turned_off
    platform: device
    device_id: <your device id>
    entity_id: binary_sensor.ups_online_status
    domain: binary_sensor
    for:
      hours: 0
      minutes: 1
      seconds: 0
condition: []
action:
  - service: shell_command.qnap_qsw_poe
    data:
      host: switch.lan
      password: !secret qnap_qsw_admin_pw
      ports: 1,2,3,4,5,6,7,8,9,10,11,12,13,14,16,19,20
      mode: disable
    mode: single
```

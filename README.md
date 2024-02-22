# qnap-qsw

A quick and dirty script to take action on a QNAP QSW-M2116P-2T2S.

## Commands

### Login
```
./qnap-qsw login --host switch.lan --password BLAH
```

### Disable/Enable POE Ports
```
./qnap-qsw poeMode --host switch.lan --password BLAH --disable-ports 1,2,3,4
./qnap-qsw poeMode --host switch.lan --password BLAH --disable-ports 1,2 --poeplus-ports 3,4 --poeplusplus-ports 19,20
```

### Home Assistant

#### shell_commands.yaml
```
qnap_qsw_poe: "/config/qnap-qsw poeMode --host {{ host }} --password {{ password }} --disable-ports '{{ disable_ports }}' --poe-ports '{{ poe_ports }}' --poeplus-ports '{{ poeplus_ports }}' --poeplusplus-ports '{{ poeplusplus_ports }}'"
```

#### automations.yaml
```
- alias: "Turn off non-essential POE ports during power outage"
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
        password: <password here>
        disable_ports: 1,2,3,4,5,6,7,8,9,10,11,12,13,14,16,19,20
  mode: single
- alias: "Turn on non-essential POE ports after power outage"
  description: ""
  trigger:
    - type: turned_on
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
        password: <password here>
        poeplus_ports: 1,2,3,4,5,6,7,8,9,10,11,12,13,14,16
        poeplusplus_ports: 19,20
  mode: single
```

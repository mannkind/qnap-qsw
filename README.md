# qnap-qsw

A quick and dirty script to take action on a QNAP QSW-M2116P-2T2S.

## Commands

### Login
```
./qnap-qsw login --host switch.lan --password BLAH
```

### Disable POE Ports
#### With Password
```
./qnap-qsw poeMode --host switch.lan --password BLAH --ports 1,2,3,4 --mode disabled
```
#### With Token
```
token=$(./qnap-qsw login --host switch.lan --password BLAH)
./qnap-qsw poeMode --host switch.lan --token $token --ports 1,2,3,4 --mode disabled
```

### Enable POE Ports
#### With Password
```
./qnap-qsw poeMode --host switch.lan --password BLAH --ports 1,2,3,4 --mode poe+
./qnap-qsw poeMode --host switch.lan --password BLAH --ports 19,20 --mode poe++
```
#### With Token
```
token=$(./qnap-qsw login --host switch.lan --password BLAH)
./qnap-qsw poeMode --host switch.lan --token $token --ports 1,2,3,4 --mode poe+
./qnap-qsw poeMode --host switch.lan --token $token --ports 19,20 --mode poe++
```

### Home Assistant

#### shell_commands.yaml
```
qnap_qsw: "/qnap-qsw login --host {{ host }} --password {{ password }}"
qnap_qsw_poe: "/config/qnap-qsw poeMode --host {{ host }} --token {{ token }} --ports {{ ports }} --mode {{ mode }}"
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
    - service: shell_command.qnap_qsw
      data:
        host: switch.lan
        password: <password here>
      response_variable: login
    - service: shell_command.qnap_qsw_poe
      data:
        host: switch.lan
        token: "{{ login['stdout'] }}"
        ports: 1,2,3,4,5,6,7,8,9,10,11,12,13,14,16,19,20
        mode: disabled
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
    - service: shell_command.qnap_qsw
      data:
        host: switch.lan
        password: <password here>
      response_variable: login
    - service: shell_command.qnap_qsw_poe
      data:
        host: switch.lan
        token: "{{ login['stdout'] }}"
        ports: 1,2,3,4,5,6,7,8,9,10,11,12,13,14,16
        mode: poe+
    - service: shell_command.qnap_qsw_poe
      data:
        host: switch.lan
        token: "{{ login['stdout'] }}"
        ports: 19,20
        mode: poe++
  mode: single
```

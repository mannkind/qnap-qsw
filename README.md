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


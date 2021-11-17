# Rabbit PubSub (Topics)

## Prerequisites

Terminal 1 - Publisher (to exchange/topic):
```bash
./pub -t messages -r system.prod.debug
```
NOTE: The routing key here is specifically system.debug

Terminal 2 - Subscriber:
```bash
./sub -t messages -r "system.*.*"
```
NOTE: This will receive a copy of all messages


Terminal 3 - Subscriber:
```bash
./sub -t messages -r "system.prod.warning"
```
NOTE: This will receive NO input, because it is looking for a routing key that doesn't exist
NOTE 2: Putting in a working routing key here (see below) will cause this to "load balance"

### Will Work
The following subscriber routing keys WILL work:
- `./sub -t messages -r "system.*.*"`
- `./sub -t messages -r "system.prod.*"`
- `./sub -t messages -r "system.*.debug"`

### Will NOT Work
The following subscriber routing keys WILL work:
- `./sub -t messages -r "system.*`
- `./sub -t messages -r "system.test.*"`
- `./sub -t messages -r "system.*.warning"`

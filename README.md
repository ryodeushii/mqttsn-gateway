# mqttsn-gateway

Default setup of MQTT-sn gateway forwarding messages to MQTT broker (RMQ+mqtt plugin for ex.)

Example configuration is placed in `data` directory. Project might be run either using `go run main.go` or using `go build` and then running produced binary `./mqttsngws` from project directory.


## Benchmarking results

####  parallel with 240 jobs
with message queue limit set to 100_000 messages
RAM: up to 5.5G with quick scale down after devicecs stop pushing messages
CPU: up to 3.5 cores in max load scenario

### from 1 to unknown amount of gateways 
with message queue limit set to 1_000 messages
RAM: up to 15MB
CPU: up to 0.1 cores


*Testing machine*
Lenovo Thinkpad T16 I7-1260P (12 cores - 4HP + 8 LP / 16 threads (vCores)), 32GB RAM


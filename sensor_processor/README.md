## Sensor Processor

A sensor processor is responsible for consuming various sensor messages/data and sending them to the notifier.

Currently, we have the edison processor, which processes sensor data from a touch sensor and sends events to the 
notifier via an MQTT broker.

The sensor processor uses the [gobot library](http://gobot.io/documentation/platforms/edison/) as a driver library 
to communicate with the Intel Edison.

### Building for Intel Edison

Cross-compile `go` for the linux architecture. This is assuming you are running on a windows or mac.
```
cd <src directory of golang>
GOOS=linux GOARCH=386 ./make.bash

cd ~/workspace/src/github.com/wfernandes/iot/sensor_processor/edison_processor
GOOS=linux GOARCH=386 go build main.go
```

Or you can use the build script in `bin/build_sensor_processor`


After the binary has been generated, `scp` it onto the Intel Edison

### Running the binary

```
./main --config=sensor.json
```
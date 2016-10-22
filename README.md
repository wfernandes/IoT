[![Build Status](https://travis-ci.org/wfernandes/IoT.svg?branch=master)](https://travis-ci.org/wfernandes/IoT)

## IoT Sensor Notification System

This is a first pass at creating a working end-to-end system that consumes sensor data/events and notifies a user.


The intention of this project was to create a platform where a user could be notified upon certain conditions of
sensor input data. 

The sensor processor which would run on something like a Raspberry Pi or Intel Edison would consume sensor data/events
 and send it to the notification processor over MQTT via a MQTT Broker. The notification processor could either run on a local machine or 
 on the cloud, say on a platform like [Cloud Foundry](http://docs.cloudfoundry.org/concepts/).
 
### Architecture

![architecture](docs/architecture.png)

### Tests
```
./bin/test
```

### Building the Components
Run `./bin/build_notification_processor` to build the notification processor.

Run `./bin/build_sensor_processor` to build the sensor processor. Currently, the sensor processor builds on a linux/amd64 system. This allows the
sensor processor binary to work on the Intel Edison and the Raspberry Pi.

### Wiki

The [wiki](https://github.com/wfernandes/IoT/wiki) contains more information regarding process, troubleshooting and setup I did so that I could reference my silliness in the future.

### Future Work

I've been tracking ideas on the [issues page](https://github.com/wfernandes/IoT/issues) and marking what is important to work on.

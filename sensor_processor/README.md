## Sensor Processor

A sensor processor is responsible for consuming various sensor messages/data and sending them to the notifier.

Currently, we have the edison processor, which processes sensor data from a touch sensor and sends events to the 
notifier via UDP.

It uses the [gobot library](http://gobot.io/documentation/platforms/edison/) to communicate with the Intel Edison.
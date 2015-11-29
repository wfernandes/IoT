## IoT Sensor Notification System

This is a first pass at creating a working end-to-end system that consumes sensor data/events and notifies a user.


The intention of this project was to create a platform where a user could be notified upon certain conditions of
sensor input data. 

The sensor processor which would run on something like a Raspberry Pi or Intel Edison would consume sensor data/events
 and send it to the notification processor over UDP. The notification processor could either run on a local machine or 
 on the cloud, say on a platform like [Cloud Foundry](http://docs.cloudfoundry.org/concepts/).
 
### Tests
```
ginkgo -r
```
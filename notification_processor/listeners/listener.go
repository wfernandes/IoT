package listeners

import (
	"net"

	"github.com/cloudfoundry/gosteno"
)

type UDPListener struct {
	address  string
	dataChan chan string
	logger   *gosteno.Logger
}

// TODO: Returns listener interface
func New(address string, outputChan chan string, logger *gosteno.Logger) *UDPListener {

	return &UDPListener{
		address:  address,
		dataChan: outputChan,
		logger:   logger,
	}
}

func (l *UDPListener) Start() {

	defer close(l.dataChan)

	conn, err := net.ListenPacket("udp", l.address)
	if err != nil {
		l.logger.Errorf("Error establishing connection: %s", err.Error())
		return
	}
	defer conn.Close()
	readBuffer := make([]byte, 65535) //buffer with size = max theoretical UDP size
	for {
		readCount, _, err := conn.ReadFrom(readBuffer)
		if err != nil {
			l.logger.Debugf("Error while reading. %s", err)
			return
		}

		readData := make([]byte, readCount) //pass on buffer in size only of read data
		copy(readData, readBuffer[:readCount])
		l.dataChan <- string(readData)
	}
}

func (l *UDPListener) Stop() {
	close(l.dataChan)
}

package writers

import "net"

type Writer struct {
	address string
	dataCh  chan string
}

// TODO: Add logging
func New(address string, dataCh chan string) *Writer {

	return &Writer{
		address: address,
		dataCh:  dataCh,
	}
}

func (w *Writer) Start() {

	conn, err := net.Dial("udp", w.address)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// TODO: log number of bytes written successfully
	for data := range w.dataCh {
		_, err := conn.Write([]byte(data))
		if err != nil {
			panic(err)
		}
	}
}

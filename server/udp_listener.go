package server

import "net"

import "github.com/Sirupsen/logrus"

import logger "../logger"

var log = logger.GetLogger()

func isError(err error) {
	if err != nil{
		log.Fatal("Error: ", err)
	}
}

type ListenerPool struct {
	UDPListeners []*net.UDPConn
	UDPBasePort int
	UDPListenerCount int
}

func (self *ListenerPool) Init() {
	log.Info("Starting Server...")
	for i:= 0; i < self.UDPListenerCount; i++ {
		addr := net.UDPAddr{
			Port: i + self.UDPBasePort,
			IP: net.ParseIP("127.0.0.1"),
		}
		conn, err := net.ListenUDP("udp", &addr)
		isError(err)
		self.UDPListeners = append(self.UDPListeners, conn)
		log.WithFields(logrus.Fields{
			"addr": addr,
		}).Info("Opened UDP")
	}
	log.WithFields(logrus.Fields{
		"listeners": self.UDPListeners,
	}).Info("Listeners Created")
}

func (self *ListenerPool) Producer (messageOut chan<- string, listener *net.UDPConn) {
	log.WithFields(logrus.Fields{
		"listener": listener,
	}).Info("Producer Start")

	defer listener.Close()

	buff := make([]byte, 1024)

	for {
		log.WithFields(logrus.Fields{
			"listener": listener,
		}).Debug("Looking At")

		n,_,err := listener.ReadFromUDP(buff)
		isError(err)

		log.WithFields(logrus.Fields{
			"num": n,
			"buff": buff,
		}).Debug("From UDP")

		if n > 0 {
			messageOut <- string(buff[0:n])
		}
	}
}

func (self *ListenerPool) StartProducers (messageOut chan<- string) {
	log.Info("Starting Producers")

	for _, listener := range self.UDPListeners {
		go self.Producer(messageOut, listener)
	}
	
	log.Info("Producers Started")
}

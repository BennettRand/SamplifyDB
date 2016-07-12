package main

import "github.com/Sirupsen/logrus"

import server "./server"
import logger "./logger"

var log = logger.GetLogger()

func worker(id int, jobs <-chan string, results chan<- int) {
	for j := range jobs {
		log.WithFields(logrus.Fields{
			"id": id,
			"job": j,
		}).Info("Working on new Job.")
	}
}

func main() {

	listeners := server.ListenerPool{UDPListenerCount: 1, UDPBasePort: 12345}
	jobs := make(chan string, 10)
	results := make(chan int)

	listeners.Init()
	listeners.StartProducers(jobs)

	for i := 0; i < 10; i++{
		go worker(i, jobs, results)
	}

	for{}
}

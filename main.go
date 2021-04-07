package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/juli3nk/openfaas-connector-stan/config"
	"github.com/juli3nk/openfaas-connector-stan/stan"
	nstan "github.com/nats-io/stan.go"
	"github.com/nats-io/stan.go/pb"
	"github.com/openfaas/connector-sdk/types"
)

const (
	clusterID = "faas-connector"
	clientID = "faas-connector-worker"
	queueGroup = "faas-connector-worker-group"
)

func main() {
	// Controller
	creds := types.GetCredentials()
	config := config.Get()

	controllerConfig := &types.ControllerConfig{
		GatewayURL:      config.GatewayURL,
		PrintResponse:   config.PrintResponse,
		RebuildInterval: config.RebuildInterval,
	}

	controller := types.NewController(creds, controllerConfig)
	controller.BeginMapBuilder()

	// Broker
	messageHandler := func(msg *nstan.Msg) {
		log.Printf("Received topic: %s, message: %s", msg.Subject, string(msg.Data))
		controller.Invoke(msg.Subject, &msg.Data)
	}

	stanQueue := stan.STANQueue{
		StanURL:   "nats://" + config.BrokerHost + ":4222",
		ClusterID: clusterID,
		ClientID:  clientID,

		MaxReconnect:   config.BrokerMaxReconnect,
		ReconnectDelay: config.BrokerReconnectDelay,

		Subjects:       config.Topics,
		QGroup:         queueGroup,
		MessageHandler: messageHandler,
		Durable:        config.BrokerDurable,
		StartOption:    nstan.StartAt(pb.StartPosition_NewOnly),
		AckWait:        config.BrokerAckWait,
		MaxInFlight:    nstan.MaxInflight(config.BrokerMaxInFlight),
	}

	if err := stanQueue.Connect(); err != nil {
		log.Panic(err)
	}

	if err := stanQueue.Subscribe(); err != nil {
		log.Panic(err)
	}

	// Wait for a SIGINT (perhaps triggered by user with CTRL-C)
	// Run cleanup when signal is received
	signalChan := make(chan os.Signal, 0)
	cleanupDone := make(chan bool)
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		for range signalChan {
			fmt.Printf("\nReceived an interrupt, unsubscribing and closing connection...\n\n")
			// Do not unsubscribe a durable on exit, except if asked to.
			if config.BrokerDurable == "" || config.BrokerUnsubscribe {
				if err := stanQueue.Unsubscribe(); err != nil {
					log.Panic(err)
				}
			}
			if err := stanQueue.CloseConnection(); err != nil {
				log.Panicf("Cannot close connection to %s because of an error: %v\n", stanQueue.StanURL, err)
			}
			cleanupDone <- true
		}
	}()
	<-cleanupDone
}

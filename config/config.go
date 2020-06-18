package config

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	BrokerHost           string
	BrokerMaxReconnect   int
	BrokerReconnectDelay time.Duration
	BrokerDurable        string
	BrokerAckWait        time.Duration
	BrokerMaxInFlight    int
	BrokerUnsubscribe    bool
	Topics               []string
	GatewayURL           string
	RebuildInterval      time.Duration
	PrintResponse        bool
}

func Get() Config {
	brokerHost := "bus"
	if val, exists := os.LookupEnv("broker_host"); exists {
		brokerHost = val
	}

	brokerMaxReconnect := 120
	if val, exists := os.LookupEnv("broker_max_reconnect"); exists {
		parsedVal, err := strconv.Atoi(val)
		if err == nil {
			brokerMaxReconnect = parsedVal
		}
	}

	brokerReconnectDelay := time.Second * 2
	if val, exists := os.LookupEnv("broker_reconnect_delay"); exists {
		parsedVal, err := time.ParseDuration(val)
		if err == nil {
			brokerReconnectDelay = parsedVal
		}
	}

	brokerDurable := ""
	if val, exists := os.LookupEnv("broker_durable"); exists {
		brokerDurable = val
	}

	brokerAckWait := time.Second * 30
	if val, exists := os.LookupEnv("broker_ack_wait"); exists {
		parsedVal, err := time.ParseDuration(val)
		if err == nil {
			brokerAckWait = parsedVal
		}
	}

	brokerMaxInFlight := 1
	if val, exists := os.LookupEnv("broker_max_in_flight"); exists {
		parsedVal, err := strconv.Atoi(val)
		if err == nil {
			brokerMaxInFlight = parsedVal
		}
	}

	brokerUnsubscribe := false
	if val, exists := os.LookupEnv("broker_unsubscribe"); exists {
		parsedVal, err := strconv.ParseBool(val)
		if err == nil {
			brokerUnsubscribe = parsedVal
		}
	}

	topics := []string{}
	if val, exists := os.LookupEnv("topics"); exists {
		for _, topic := range strings.Split(val, ",") {
			if len(topic) > 0 {
				topics = append(topics, topic)
			}
		}
	}
	if len(topics) == 0 {
		log.Fatal(`Provide a list of topics i.e. topics="payment_published,slack_joined"`)
	}

	gatewayURL := "http://gateway:8080"
	if val, exists := os.LookupEnv("gateway_url"); exists {
		gatewayURL = val
	}

	rebuildInterval := time.Second * 3
	if val, exists := os.LookupEnv("rebuild_interval"); exists {
		parsedVal, err := time.ParseDuration(val)
		if err == nil {
			rebuildInterval = parsedVal
		}
	}

	printResponse := false
	if val, exists := os.LookupEnv("print_response"); exists {
		printResponse = (val == "1" || val == "true")
	}

	return Config{
		BrokerHost:           brokerHost,
		BrokerMaxReconnect:   brokerMaxReconnect,
		BrokerReconnectDelay: brokerReconnectDelay,
		BrokerDurable:        brokerDurable,
		BrokerAckWait:        brokerAckWait,
		BrokerMaxInFlight:    brokerMaxInFlight,
		BrokerUnsubscribe:    brokerUnsubscribe,
		Topics:               topics,
		GatewayURL:           gatewayURL,
		RebuildInterval:      rebuildInterval,
		PrintResponse:        printResponse,
	}
}

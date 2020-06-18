package stan

import (
	"fmt"
	"log"
	"sync"
	"time"

	nstan "github.com/nats-io/stan.go"
)

// connect creates a subscription to NATS Streaming
func (q *STANQueue) Connect() error {
	log.Printf("Connect: %s\n", q.StanURL)

	q.connMutex = &sync.RWMutex{}
	q.quitCh = make(chan struct{})

	sc, err := nstan.Connect(
		q.ClusterID,
		q.ClientID,
		nstan.NatsURL(q.StanURL),
		nstan.SetConnectionLostHandler(func(conn nstan.Conn, err error) {
			log.Printf("Disconnected from %s\n", q.StanURL)

			q.reconnect()
		}),
	)
	if err != nil {
		return fmt.Errorf("can't connect to %s: %v", q.StanURL, err)
	}

	q.connMutex.Lock()
	defer q.connMutex.Unlock()

	q.conn = sc

	q.subscription = make(map[string]nstan.Subscription)

	return nil
}

func (q *STANQueue) reconnect() {
	log.Printf("Reconnect\n")

	for i := 0; i < q.MaxReconnect; i++ {
		select {
		case <-time.After(time.Duration(i) * q.ReconnectDelay):
			if err := q.Connect(); err == nil {
				log.Printf("Reconnecting (%d/%d) to %s succeeded\n", i+1, q.MaxReconnect, q.StanURL)

				if err := q.Subscribe(); err != nil {
					log.Panic(err)
				}

				return
			}

			nextTryIn := (time.Duration(i+1) * q.ReconnectDelay).String()

			log.Printf("Reconnecting (%d/%d) to %s failed\n", i+1, q.MaxReconnect, q.StanURL)
			log.Printf("Waiting %s before next try", nextTryIn)
		case <-q.quitCh:
			log.Println("Received signal to stop reconnecting...")

			return
		}
	}

	log.Printf("Reconnecting limit (%d) reached\n", q.MaxReconnect)
}

func (q *STANQueue) Subscribe() error {
	for _, subject := range q.Subjects {
		log.Printf("Subscribing to: %s at %s\n", subject, q.StanURL)
		log.Println("Wait for ", q.AckWait)

		subscription, err := q.conn.QueueSubscribe(
			subject,
			q.QGroup,
			q.MessageHandler,
			q.StartOption,
			nstan.DurableName(q.Durable),
			nstan.AckWait(q.AckWait),
			q.MaxInFlight,
		)
		if err != nil {
			return fmt.Errorf("couldn't subscribe to %s at %s. Error: %v", subject, q.StanURL, err)
		}

		log.Printf(
			"Listening on [%s], clientID=[%s], qgroup=[%s] durable=[%s]\n",
			subject,
			q.ClientID,
			q.QGroup,
			q.Durable,
		)

		q.subscription[subject] = subscription
	}

	return nil
}

func (q *STANQueue) Unsubscribe() error {
	q.connMutex.Lock()
	defer q.connMutex.Unlock()

	for _, subject := range q.Subjects {
		if q.subscription[subject] != nil {
			return fmt.Errorf("q.subscription[%s] is nil", subject)
		}

		if err := q.subscription[subject].Unsubscribe(); err != nil {
			return fmt.Errorf(
				"Cannot unsubscribe subject: %s from %s because of an error: %v",
				subject,
				q.StanURL,
				err,
			)
		}
	}

	return nil
}

func (q *STANQueue) CloseConnection() error {
	q.connMutex.Lock()
	defer q.connMutex.Unlock()

	if q.conn == nil {
		return fmt.Errorf("q.conn is nil")
	}

	close(q.quitCh)

	return q.conn.Close()
}

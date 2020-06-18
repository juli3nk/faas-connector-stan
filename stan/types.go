package stan

import (
	"sync"
	"time"

	nstan "github.com/nats-io/stan.go"
)

// STANQueue represents a subscription to NATS Streaming
type STANQueue struct {
	StanURL   string
	ClusterID string
	ClientID  string

	conn      nstan.Conn
	connMutex *sync.RWMutex
	quitCh    chan struct{}

	MaxReconnect   int
	ReconnectDelay time.Duration

	Subjects       []string
	QGroup         string
	MessageHandler func(*nstan.Msg)
	StartOption    nstan.SubscriptionOption
	Durable        string
	AckWait        time.Duration
	MaxInFlight    nstan.SubscriptionOption

	subscription map[string]nstan.Subscription
}

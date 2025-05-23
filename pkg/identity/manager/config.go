package manager

import (
	"time"

	"github.com/pomerium/pomerium/internal/events"
	"github.com/pomerium/pomerium/pkg/grpc/databroker"
)

var (
	defaultSessionRefreshGracePeriod     = 1 * time.Minute
	defaultSessionRefreshCoolOffDuration = 10 * time.Second
	defaultUpdateUserInfoInterval        = 10 * time.Minute
	defaultLeaseTTL                      = 30 * time.Second
)

type config struct {
	authenticator                 Authenticator
	dataBrokerClient              databroker.DataBrokerServiceClient
	sessionRefreshGracePeriod     time.Duration
	sessionRefreshCoolOffDuration time.Duration
	updateUserInfoInterval        time.Duration
	leaseTTL                      time.Duration
	now                           func() time.Time
	eventMgr                      *events.Manager
}

func newConfig(options ...Option) *config {
	cfg := new(config)
	WithSessionRefreshGracePeriod(defaultSessionRefreshGracePeriod)(cfg)
	WithSessionRefreshCoolOffDuration(defaultSessionRefreshCoolOffDuration)(cfg)
	WithNow(time.Now)(cfg)
	WithUpdateUserInfoInterval(defaultUpdateUserInfoInterval)(cfg)
	WithLeaseTTL(defaultLeaseTTL)(cfg)
	for _, option := range options {
		option(cfg)
	}
	return cfg
}

// An Option customizes the configuration used for the identity manager.
type Option func(*config)

// WithAuthenticator sets the authenticator in the config.
func WithAuthenticator(authenticator Authenticator) Option {
	return func(cfg *config) {
		cfg.authenticator = authenticator
	}
}

// WithDataBrokerClient sets the databroker client in the config.
func WithDataBrokerClient(dataBrokerClient databroker.DataBrokerServiceClient) Option {
	return func(cfg *config) {
		cfg.dataBrokerClient = dataBrokerClient
	}
}

// WithSessionRefreshGracePeriod sets the session refresh grace period used by the manager.
func WithSessionRefreshGracePeriod(dur time.Duration) Option {
	return func(cfg *config) {
		cfg.sessionRefreshGracePeriod = dur
	}
}

// WithSessionRefreshCoolOffDuration sets the session refresh cool-off duration used by the manager.
func WithSessionRefreshCoolOffDuration(dur time.Duration) Option {
	return func(cfg *config) {
		cfg.sessionRefreshCoolOffDuration = dur
	}
}

// WithNow customizes the time.Now function used by the manager.
func WithNow(now func() time.Time) Option {
	return func(cfg *config) {
		cfg.now = now
	}
}

// WithEventManager passes an event manager to record events
func WithEventManager(mgr *events.Manager) Option {
	return func(cfg *config) {
		cfg.eventMgr = mgr
	}
}

// WithUpdateUserInfoInterval sets the update user info interval in the config.
func WithUpdateUserInfoInterval(dur time.Duration) Option {
	return func(cfg *config) {
		cfg.updateUserInfoInterval = dur
	}
}

// WithLeaseTTL sets the TTL used by the leaser.
func WithLeaseTTL(ttl time.Duration) Option {
	return func(o *config) {
		o.leaseTTL = ttl
	}
}

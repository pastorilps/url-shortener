package cassandra

import (
	"fmt"

	"github.com/gocql/gocql"
)

type Config struct {
	Hosts    []string
	Keyspace string
	Username string
	Password string
}

func NewConnection(cfg Config) (*gocql.Session, error) {
	cluster := gocql.NewCluster(cfg.Hosts...)

	cluster.Keyspace = cfg.Keyspace
	cluster.Consistency = gocql.Quorum

	if cfg.Username != "" {
		cluster.Authenticator = gocql.PasswordAuthenticator{
			Username: cfg.Username,
			Password: cfg.Password,
		}
	}

	session, err := cluster.CreateSession()
	if err != nil {
		return nil, fmt.Errorf("Cassandra connection: %w", err)
	}

	return session, nil
}

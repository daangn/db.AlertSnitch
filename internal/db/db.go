package db

import (
	"fmt"

	"github.com/daangn/db.AlertSnitch/internal"
)

// SupportedModel stores the model that is supported by this application
// daangn use 1.0.0 Ver
const SupportedModel = "1.0.0"

// ConnectionArgs required to create a MySQL connection
type ConnectionArgs struct {
	DSN                    string
	MaxIdleConns           int
	MaxOpenConns           int
	MaxConnLifetimeSeconds int
}

// Connect connects to a backend database
func Connect(backend string, args ConnectionArgs) (internal.Storer, error) {
	switch backend {
	case "mysql":
		return connectMySQL(args)

	case "postgres":
		return connectPG(args)

	case "null":
		return NullDB{}, nil

	default:
		return nil, fmt.Errorf("Invalid backend %q", backend)
	}
}

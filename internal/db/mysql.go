package db

import (
	"context"
	"fmt"
	"time"

	"database/sql"

	"github.com/daangn/db.AlertSnitch/internal"
	"github.com/daangn/db.AlertSnitch/internal/metrics"
	"github.com/sirupsen/logrus"
)

// MySQLDB A database that does nothing
type MySQLDB struct {
	db *sql.DB
}

// ConnectMySQL connect to a MySQL database using the provided data source name
func connectMySQL(args ConnectionArgs) (*MySQLDB, error) {
	if args.DSN == "" {
		return nil, fmt.Errorf("Empty DSN provided, can't connect to MySQL database")
	}

	logrus.Debugf("Connecting to MySQL database with DSN: %s", args.DSN)

	connection, err := sql.Open("mysql", args.DSN)
	if err != nil {
		return nil, fmt.Errorf("failed to open MySQL connection: %s", err)
	}

	connection.SetMaxIdleConns(args.MaxIdleConns)
	connection.SetMaxOpenConns(args.MaxOpenConns)
	connection.SetConnMaxLifetime(time.Duration(args.MaxConnLifetimeSeconds) * time.Second)

	database := &MySQLDB{
		db: connection,
	}

	err = database.Ping()
	if err != nil {
		return nil, err
	}
	logrus.Debug("Connected to MySQL database")

	return database, database.CheckModel()
}

// Save implements Storer interface
func (d MySQLDB) Save(data *internal.AlertGroup) error {
	return d.unitOfWork(func(tx *sql.Tx) error {
		r, err := tx.Exec(`
			INSERT INTO alert_group (time, receiver, status, external_url, group_key)
			VALUES (now(), ?, ?, ?, ?)`, data.Receiver, data.Status, data.ExternalURL, data.GroupKey)
		if err != nil {
			return fmt.Errorf("failed to insert into alert_groups: %s", err)
		}

		alertGroupID, err := r.LastInsertId() // alertGroupID
		if err != nil {
			return fmt.Errorf("failed to get alert_groups inserted id: %s", err)
		}

		for k, v := range data.GroupLabels {
			_, err := tx.Exec(`
				INSERT INTO group_label (alert_group_id, group_label, value)
				VALUES (?, ?, ?)`, alertGroupID, k, v)
			if err != nil {
				return fmt.Errorf("failed to insert into group_label: %s", err)
			}
		}
		for k, v := range data.CommonLabels {
			_, err := tx.Exec(`
				INSERT INTO common_label (alert_group_id, label, value)
				VALUES (?, ?, ?)`, alertGroupID, k, v)
			if err != nil {
				return fmt.Errorf("failed to insert into common_label: %s", err)
			}
		}
		for k, v := range data.CommonAnnotations {
			_, err := tx.Exec(`
				INSERT INTO common_annotation (alert_group_id, annotation, value)
				VALUES (?, ?, ?)`, alertGroupID, k, v)
			if err != nil {
				return fmt.Errorf("failed to insert into common_annotation: %s", err)
			}
		}

		for _, alert := range data.Alerts {
			var result sql.Result
			if alert.EndsAt.Before(alert.StartsAt) {
				result, err = tx.Exec(`
				INSERT INTO alert (alert_group_id, status, starts_at, generator_url, fingerprint)
				VALUES (?, ?, ?, ?, ?)`,
					alertGroupID, alert.Status, alert.StartsAt, alert.GeneratorURL, alert.Fingerprint)
			} else {
				result, err = tx.Exec(`
				INSERT INTO alert (alert_group_id, status, starts_at, ends_at, generator_url, fingerprint)
				VALUES (?, ?, ?, ?, ?, ?)`,
					alertGroupID, alert.Status, alert.StartsAt, alert.EndsAt, alert.GeneratorURL, alert.Fingerprint)
			}
			if err != nil {
				return fmt.Errorf("failed to insert into alert: %s", err)
			}

			alertID, err := result.LastInsertId()
			if err != nil {
				return fmt.Errorf("failed to get alert inserted id: %s", err)
			}

			for k, v := range alert.Labels {
				_, err := tx.Exec(`
					INSERT INTO alert_label (alert_id, label, value)
					VALUES (?, ?, ?)`, alertID, k, v)
				if err != nil {
					return fmt.Errorf("failed to insert into alert_label: %s", err)
				}
			}
			for k, v := range alert.Annotations {
				_, err := tx.Exec(`
					INSERT INTO alert_annotation (alert_id, annotation, value)
					VALUES (?, ?, ?)`, alertID, k, v)
				if err != nil {
					return fmt.Errorf("failed to insert into alert_annotation: %s", err)
				}
			}
		}

		return nil
	})
}

func (d MySQLDB) unitOfWork(f func(*sql.Tx) error) error {
	tx, err := d.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %s", err)
	}

	err = f(tx)

	if err != nil {
		e := tx.Rollback()
		if e != nil {
			return fmt.Errorf("failed to rollback transaction (%s) after failing execution: %s", e, err)
		}
		return fmt.Errorf("failed execution: %s", err)
	}
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %s", err)
	}
	return nil
}

// Ping implements Storer interface
func (d MySQLDB) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	if err := d.db.PingContext(ctx); err != nil {
		metrics.DatabaseUp.Set(0)
		logrus.Debugf("Failed to ping database: %s", err)
		return err
	}
	metrics.DatabaseUp.Set(1)

	logrus.Debugf("Pinged database...")
	return nil
}

// CheckModel implements Storer interface
func (d MySQLDB) CheckModel() error {
	rows, err := d.db.Query("SELECT version FROM Model")
	if err != nil {
		return fmt.Errorf("failed to fetch model version from the database: %s", err)
	}
	defer rows.Close()

	if !rows.Next() {
		return fmt.Errorf("failed to read model version from the database: empty resultset")
	}

	var model string
	if err := rows.Scan(&model); err != nil {
		return fmt.Errorf("failed to read model version from the database: %s", err)
	}

	if model != SupportedModel {
		return fmt.Errorf("database model '%s' is not supported by this application (%s)", model, SupportedModel)
	}

	return nil
}

func (MySQLDB) String() string {
	return "mysql database driver"
}

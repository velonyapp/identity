package mysql

import (
	"database/sql"

	"github.com/velony-app/identity/internal/conf"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"

	"github.com/XSAM/otelsql"
	"github.com/go-sql-driver/mysql"
)

func NewConnection(c *conf.Data) (*sql.DB, error) {
	cfg, err := mysql.ParseDSN(c.GetMysql().GetDsn())
	if err != nil {
		return nil, err
	}

	cfg.ParseTime = true
	cfg.InterpolateParams = true

	dsn := cfg.FormatDSN()

	attrs := append(otelsql.AttributesFromDSN(dsn), semconv.DBSystemNameMySQL)

	db, err := otelsql.Open("mysql", dsn,
		otelsql.WithAttributes(attrs...),
		otelsql.WithSpanOptions(otelsql.SpanOptions{
			OmitConnResetSession: true,
			OmitRows:             true,
		}),
	)
	if err != nil {
		return nil, err
	}

	if c.GetMysql().GetMaxOpenConnections() != 0 {
		db.SetMaxOpenConns(int(c.GetMysql().GetMaxOpenConnections()))
	}
	if c.GetMysql().GetMaxIdleConnections() != 0 {
		db.SetMaxIdleConns(int(c.GetMysql().GetMaxIdleConnections()))
	}
	if c.GetMysql().GetMaxConnectionLifetime() != nil {
		db.SetConnMaxLifetime(c.GetMysql().GetMaxConnectionLifetime().AsDuration())
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

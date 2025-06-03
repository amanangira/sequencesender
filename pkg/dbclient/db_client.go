package dbclient

import (
	"fmt"
	"log"
	"os"
	"strings"

	"sequencesender"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func GetDBConnectionString() string {
	dsn := fmt.Sprintf(
		`host=%s port=%s user=%s password=%s database=%s`,
		os.Getenv(sequencesender.EnvDBHostKey),
		os.Getenv(sequencesender.EnvDBPortKey),
		os.Getenv(sequencesender.EnvDBUsernameKey),
		os.Getenv(sequencesender.EnvDBPasswordKey),
		os.Getenv(sequencesender.EnvDBNameKey),
	)

	if strings.ToLower(os.Getenv(sequencesender.EnvDBSSL)) == sequencesender.EnvTrueValue {
		dsn += fmt.Sprintf(" sslmode=require")
	} else {
		dsn += fmt.Sprintf(" sslmode=disable")
	}

	log.Println("connection", dsn)

	return dsn
}

func NewSQLXConnection(dsn string) (*sqlx.DB, error) {
	sqlDB, err := sqlx.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	return sqlDB, nil
}

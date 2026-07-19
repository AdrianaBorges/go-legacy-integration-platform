package database

import (
	"database/sql"
	"fmt"
	"os"
)

type OracleConfig struct {
	User          string
	Password      string
	ConnectString string
}

func LoadOracleConfig() OracleConfig {
	return OracleConfig{
		User:          os.Getenv("ORACLE_USER"),
		Password:      os.Getenv("ORACLE_PASSWORD"),
		ConnectString: os.Getenv("ORACLE_CONNECT_STRING"),
	}
}

func OpenOracle(config OracleConfig) (*sql.DB, error) {
	if config.User == "" ||
		config.Password == "" ||
		config.ConnectString == "" {
		return nil, fmt.Errorf("configuração do Oracle incompleta")
	}

	connectionString := fmt.Sprintf(
		`oracle://%s:%s@%s`,
		config.User,
		config.Password,
		config.ConnectString,
	)

	db, err := sql.Open("oracle", connectionString)
	if err != nil {
		return nil, fmt.Errorf("erro ao abrir conexão Oracle: %w", err)
	}

	return db, nil
}

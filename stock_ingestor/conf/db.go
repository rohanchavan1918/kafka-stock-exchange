package conf

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type DB struct {
	DBType string `viper:"string" mapstructure:"db_type"`
	DBHost string `viper:"string" mapstructure:"db_host"`
	DBPort string `viper:"string" mapstructure:"db_port"`
	DBUser string `viper:"string" mapstructure:"db_user"`
	DBPass string `viper:"string" mapstructure:"db_pass"`
	DBName string `viper:"string" mapstructure:"db_name"`
}

func (db *DB) GetConnectionString(dbType string) string {

	var connectionString string
	switch dbType {

	case "mysql":
		connectionString = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
			db.DBUser, db.DBPass, db.DBHost, db.DBPort, db.DBName)

	case "postgres":
		connectionString = fmt.Sprintf("host=%s port=%s user=%s "+
			"password=%s dbname=%s sslmode=disable",
			db.DBHost, db.DBPort, db.DBUser, db.DBPass, db.DBName)
	}

	return connectionString
}

func GetDBConnection(dbConfig *DB) *sql.DB {
	connectionString := dbConfig.GetConnectionString(dbConfig.DBType)

	// sql.Register(dbConfig.DBType, &mysql.MySQLDriver{})
	db, err := sql.Open(dbConfig.DBType, connectionString)

	if err != nil {
		panic(err)
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	return db
}

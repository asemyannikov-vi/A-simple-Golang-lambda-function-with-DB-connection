package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	_ "github.com/lib/pq"
)

var (
	user     = os.Getenv("DB_USER")
	password = os.Getenv("DB_PASSWORD")
	port     = os.Getenv("DB_PORT")
	host     = os.Getenv("DB_HOST")
	schema   = os.Getenv("DB_SCHEMA")
)

const (
	internalErrorCode = 500
	successfullCode   = 200
	resourceNotFound  = 404
	dbHealthSql       = "./db.health.sql"
	migrationUpSql    = "./migration.up.sql"
	migrationDownSql  = "./migration.down.sql"
)

func stringToMessage(str string) string {
	type Message struct {
		Message string `json:"message"`
	}
	bytes, err := json.Marshal(Message{Message: str})
	if err != nil {
		panic(err)
	}
	return string(bytes)
}

func connectToDataBase() (*sql.DB, string, int, error) {
	db, err := sql.Open("postgres", fmt.Sprintf("port=%s host=%s user=%s password=%s sslmode=disable", port, host, user, password))
	if err != nil {
		return db, "Failed to connect to database.", internalErrorCode, err
	}
	err = db.Ping()
	if err != nil {
		fmt.Println()
		return db, "Failed to ping to database.", internalErrorCode, err
	}
	return db, "Database connection successfully established.", successfullCode, nil
}

func queryRun(db *sql.DB, queryFilePath string) (*sql.Rows, error) {
	query, err := ioutil.ReadFile(queryFilePath)
	if err != nil {
		return nil, err
	}
	rows, err := db.Query(string(query))
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func queryHandler(request events.APIGatewayProxyRequest, db *sql.DB) (string, int, error) {
	switch request.HTTPMethod {
	case "POST":
		_, err := queryRun(db, migrationUpSql)
		if err != nil {
			return "Failed to create Schema to DataBase.", internalErrorCode, err
		}
		return "Schema '" + schema + "' has been created.", successfullCode, nil
	case "GET":
		rows, err := queryRun(db, dbHealthSql)
		if err != nil {
			return "Failed to check Health of DataBase query.", internalErrorCode, err
		}
		i := 0
		for rows.Next() {
			err := rows.Scan(&i)
			if err != nil {
				return "Failed to parse result of query.", internalErrorCode, err
			}
		}
		if i == 0 {
			return "Schema '" + schema + "' has not been created.", resourceNotFound, nil
		}
		return "Schema '" + schema + "' has been created and alive.", successfullCode, nil
	case "DELETE":
		_, err := queryRun(db, migrationDownSql)
		if err != nil {
			return "Failed to drop Schema from DataBase.", internalErrorCode, err
		}
		return "Schema '" + schema + "' has been removed.", successfullCode, nil
	default:
		return "Internal Server Error", internalErrorCode, nil
	}
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	db, msg, code, err := connectToDataBase()
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       stringToMessage(msg),
			StatusCode: code,
		}, nil
	} else {
		defer db.Close()
		msg, code, err := queryHandler(request, db)
		if err != nil {
			return events.APIGatewayProxyResponse{
				Body:       stringToMessage(msg),
				StatusCode: code,
			}, nil
		}
		return events.APIGatewayProxyResponse{
			Body:       stringToMessage(msg),
			StatusCode: code,
		}, nil
	}
}

func main() {
	lambda.Start(handler)
}

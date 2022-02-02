package config

import (
	"fmt"
	"os"
)

var (
	DBURL                              = ""
	SECRETKEY                          = ""
	ENV                                = "DEV"
	GOOGLE_PROJECT_ID                  = ""
	GOOGLE_TOPIC_ID                    = ""
	GOOGLE_TYPE                        = ""
	GOOGLE_PRIVATE_KEY_ID              = ""
	GOOGLE_PRIVATE_KEY                 = ""
	GOOGLE_CLIENT_EMAIL                = ""
	GOOGLE_CLIENT_ID                   = ""
	GOOGLE_AUTH_URI                    = ""
	GOOGLE_TOKEN_URI                   = ""
	GOOGLE_AUTH_PROVIDER_x509_CERT_URL = ""
	GOOGLE_CLIENT_x509_CERT_URL        = ""
)

func LoadEnv() {
	ENV = os.Getenv("ENV")
	var (
		username string
		password string
		host     string
		database string
		port     string
	)

	if ENV != "TEST" {
		username = os.Getenv("DB_USER")
		password = os.Getenv("DB_PASSWORD")
		host = os.Getenv("DB_HOST")
		database = os.Getenv("DB_NAME")
		port = os.Getenv("DB_PORT")

		SECRETKEY = os.Getenv("SECRET_KEY")
		GOOGLE_PROJECT_ID = os.Getenv("GOOGLE_PROJECT_ID")
		GOOGLE_TOPIC_ID = os.Getenv("GOOGLE_TOPIC_ID")
		GOOGLE_TYPE = os.Getenv("GOOGLE_TYPE")
		GOOGLE_PRIVATE_KEY_ID = os.Getenv("GOOGLE_PRIVATE_KEY_ID")
		GOOGLE_PRIVATE_KEY = os.Getenv("GOOGLE_PRIVATE_KEY")
		GOOGLE_CLIENT_EMAIL = os.Getenv("GOOGLE_CLIENT_EMAIL")
		GOOGLE_CLIENT_ID = os.Getenv("GOOGLE_CLIENT_ID")
		GOOGLE_AUTH_URI = os.Getenv("GOOGLE_AUTH_URI")
		GOOGLE_TOKEN_URI = os.Getenv("GOOGLE_TOKEN_URI")
		GOOGLE_AUTH_PROVIDER_x509_CERT_URL = os.Getenv("GOOGLE_AUTH_PROVIDER_x509_CERT_URL")
		GOOGLE_CLIENT_x509_CERT_URL = os.Getenv("GOOGLE_CLIENT_x509_CERT_URL")
	} else {
		username = os.Getenv("TEST_DB_USER")
		password = os.Getenv("TEST_DB_PASSWORD")
		host = os.Getenv("TEST_DB_HOST")
		database = os.Getenv("TEST_DB_NAME")
		port = os.Getenv("TEST_DB_PORT")

		SECRETKEY = os.Getenv("TEST_SECRET_KEY")
		GOOGLE_PROJECT_ID = os.Getenv("TEST_GOOGLE_PROJECT_ID")
		GOOGLE_TOPIC_ID = os.Getenv("TEST_GOOGLE_TOPIC_ID")
		GOOGLE_TYPE = os.Getenv("TEST_GOOGLE_TYPE")
		GOOGLE_PRIVATE_KEY_ID = os.Getenv("TEST_GOOGLE_PRIVATE_KEY_ID")
		GOOGLE_PRIVATE_KEY = os.Getenv("TEST_GOOGLE_PRIVATE_KEY")
		GOOGLE_CLIENT_EMAIL = os.Getenv("TEST_GOOGLE_CLIENT_EMAIL")
		GOOGLE_CLIENT_ID = os.Getenv("TEST_GOOGLE_CLIENT_ID")
		GOOGLE_AUTH_URI = os.Getenv("TEST_GOOGLE_AUTH_URI")
		GOOGLE_TOKEN_URI = os.Getenv("TEST_GOOGLE_TOKEN_URI")
		GOOGLE_AUTH_PROVIDER_x509_CERT_URL = os.Getenv("TEST_GOOGLE_AUTH_PROVIDER_x509_CERT_URL")
		GOOGLE_CLIENT_x509_CERT_URL = os.Getenv("TEST_GOOGLE_CLIENT_x509_CERT_URL")
	}

	DBURL = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		username,
		password,
		host,
		port,
		database,
	)
}

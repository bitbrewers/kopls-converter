package main

import (
	"log"
	"os"

	"github.com/bitbrewers/kopls-converter/server"
)

func main() {
	httpPort, dbUrl, admin, adminPasswd, domain := parseArgs()

	dbClient, err := server.NewClient(dbUrl)
	if err != nil {
		log.Fatal(err)
	}

	if err = server.NewServer(":"+httpPort, dbClient, admin, adminPasswd, domain).Start(); err != nil {
		log.Fatal(err)
	}
}

func parseArgs() (port, dbUrl, admin, adminPasswd, domain string) {
	args := []struct {
		name   string
		target *string
	}{
		{"PORT", &port},
		{"DATABASE_URL", &dbUrl},
		{"ADMIN_USER", &admin},
		{"ADMIN_PASSWD", &adminPasswd},
		{"DOMAIN", &domain},
	}

	for _, arg := range args {
		if *arg.target = os.Getenv(arg.name); *arg.target == "" {
			log.Fatalf("$%s must be set", arg.name)
		}
	}

	return
}

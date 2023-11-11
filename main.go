package main

import (
	"bufio"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gocql/gocql"
)

func main() {
	// Defining command line flags
	host := flag.String("host", "127.0.0.1", "Database host")
	port := flag.String("port", "9042", "Database port")
	user := flag.String("user", "", "Username for authentication")
	password := flag.String("password", "", "Password for authentication")
	caPath := flag.String("ca", "", "Path to CA certificate")
	certPath := flag.String("cert", "", "Path to client certificate")
	keyPath := flag.String("key", "", "Path to client key")
	verify := flag.Bool("verify", false, "Enable SSL host verification")

	// Parsing the flags
	flag.Parse()

	// Setting up the cluster configuration
	cluster := gocql.NewCluster(*host + ":" + *port)

	// SSL options
	var caCertPool *x509.CertPool
	if *caPath != "" {
		caCert, err := os.ReadFile(*caPath)
		if err != nil {
			log.Fatalf("could not read CA certificate: %v", err)
		}
		caCertPool = x509.NewCertPool()
		if !caCertPool.AppendCertsFromPEM(caCert) {
			log.Fatalf("failed to append CA certificate")
		}
	}

	var clientCert tls.Certificate
	var err error
	if *certPath != "" && *keyPath != "" {
		clientCert, err = tls.LoadX509KeyPair(*certPath, *keyPath)
		if err != nil {
			log.Fatalf("could not load client key pair: %v", err)
		}
	}

	cluster.SslOpts = &gocql.SslOptions{
		CaPath:   *caPath,
		CertPath: *certPath,
		KeyPath:  *keyPath,
		Config: &tls.Config{
			Certificates:       []tls.Certificate{clientCert},
			RootCAs:            caCertPool,
			InsecureSkipVerify: !*verify,
		},
		EnableHostVerification: *verify,
	}

	// Authentication
	if *user != "" && *password != "" {
		cluster.Authenticator = gocql.PasswordAuthenticator{
			Username: *user,
			Password: *password,
		}
	}

	cluster.Consistency = gocql.Quorum

	// Creating a session
	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatalf("Could not connect to ScyllaDB: %v", err)
	}
	defer session.Close()

	// Interactive command line interface
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Connected to ScyllaDB. Enter your CQL commands, or type 'exit' to quit:")

	for {
		fmt.Print("CQL> ")
		query, _ := reader.ReadString('\n')

		query = strings.TrimSpace(query)

		if query == "exit" {
			break
		}

		// Executing the query
		if err := session.Query(query).Exec(); err != nil {
			log.Printf("Error executing query: %v\n", err)
			continue
		}

		// Querying and displaying results
		iter := session.Query(query).Iter()
		columns := iter.Columns()
		if columns == nil {
			fmt.Println("Query executed successfully (no output to display).")
			continue
		}

		fmt.Println(strings.Repeat("-", 40))
		for _, col := range columns {
			fmt.Printf("%-20s", col.Name)
		}
		fmt.Println()
		fmt.Println(strings.Repeat("-", 40))

		rowData := make(map[string]interface{})
		for iter.MapScan(rowData) {
			for _, col := range columns {
				fmt.Printf("%-20v", rowData[col.Name])
			}
			fmt.Println()
			rowData = make(map[string]interface{})
		}
		fmt.Println(strings.Repeat("-", 40))
		if err := iter.Close(); err != nil {
			log.Printf("Error closing iterator: %v\n", err)
		}
	}

	fmt.Println("Exiting ScyllaDB CLI.")
}

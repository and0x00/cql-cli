# 📟 CQL-cli

## Overview
This project provides a command-line interface (CLI) tool to interact with ScyllaDB databases. Written in Go, it's designed to offer a simple and direct way to execute CQL commands against a ScyllaDB instance. It supports various configurations, including SSL options and authentication.

## Prerequisites
- Access to a ScyllaDB database
- SSL certificates (if connecting to a database using SSL)

## Installation
Clone the repository and build the project:

```bash
git clone https://github.com/and0x00/cql-cli.git
cd cql-cli
go build -o cql-cli
```

## Usage
Run the CLI tool using the following syntax:

```bash
./cql-cli [options]
```

### Options
- `-host <host>`: Database host (default "127.0.0.1").
- `-port <port>`: Database port (default "9042").
- `-user <username>`: Username for authentication.
- `-password <password>`: Password for authentication.
- `-ca <ca_certificate_path>`: Path to CA certificate.
- `-cert <client_certificate_path>`: Path to client certificate.
- `-key <client_key_path>`: Path to client key.
- `-verify`: Enable SSL host verification (default false).

### Example
Connect to a ScyllaDB instance:

```bash
./cql-cli -host 54.214.170.10 -port 9042 -user cassandra -password cassandra -ca /path/to/ca.crt -cert /path/to/client.crt -key /path/to/client.key -verify
```

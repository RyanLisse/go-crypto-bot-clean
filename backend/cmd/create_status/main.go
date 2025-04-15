package main

import (
	"fmt"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/scripts"
)

func main() {
	fmt.Println("Creating status table...")
	scripts.CreateStatusTable()
}

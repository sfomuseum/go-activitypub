package main

//  ./bin/create-dynamodb-tables -dynamodb-client-uri 'awsdynamodb://?region=us-east-1&credentials=default' -table activities -table-prefix collection_ap_

import (
	"context"
	"log"

	"github.com/sfomuseum/go-activitypub/app/dynamodb/tables/create"
)

func main() {

	ctx := context.Background()
	err := create.Run(ctx)

	if err != nil {
		log.Fatalf("Failed to create tables, %v", err)
	}
}

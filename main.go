package main

import (
	"log"
	"time"

	"github.com/luke92/LocalStackCloudTestGolang/internal/bucket"
	"github.com/luke92/LocalStackCloudTestGolang/internal/pkg/cloud/aws"
)

func main() {
	// Create a session instance.
	ses, err := aws.New(aws.Config{
		Address: "http://localhost:4566",
		Region:  "sa-east-1",
		Profile: "localstack",
		ID:      "c7l59f",
		Secret:  "03mhyk",
	})
	if err != nil {
		log.Fatalln(err)
	}

	// Test bucket
	bucket.Bucket(aws.NewS3(ses, time.Second*5))
}

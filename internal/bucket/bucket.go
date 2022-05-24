package bucket

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/luke92/LocalStackCloudTestGolang/internal/pkg/cloud"
)

func Bucket(client cloud.BucketClient) {
	ctx := context.Background()

	create(ctx, client)
	uploadObject(ctx, client)
	downloadObject(ctx, client)
	deleteObject(ctx, client)
	listObjects(ctx, client)
}

func create(ctx context.Context, client cloud.BucketClient) {
	if err := client.Create(ctx, "aws-test"); err != nil {
		log.Fatalln(err)
	}
	log.Println("create: ok")
}

func uploadObject(ctx context.Context, client cloud.BucketClient) {
	file, err := os.Open("./assets/id.txt")
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	url, err := client.UploadObject(ctx, "aws-test", "id.txt", file)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("upload object:", url)
}

func downloadObject(ctx context.Context, client cloud.BucketClient) {
	file, err := os.Create("./tmp/id.txt")
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	if err := client.DownloadObject(ctx, "aws-test", "id.txt", file); err != nil {
		log.Fatalln(err)
	}
	log.Println("download object: ok")
}

func deleteObject(ctx context.Context, client cloud.BucketClient) {
	if err := client.DeleteObject(ctx, "aws-test", "id.txt"); err != nil {
		log.Fatalln(err)
	}
	log.Println("delete object: ok")
}

func listObjects(ctx context.Context, client cloud.BucketClient) {
	objects, err := client.ListObjects(ctx, "aws-test")
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("list objects:")
	for _, object := range objects {
		fmt.Printf("%+v\n", object)
	}
}

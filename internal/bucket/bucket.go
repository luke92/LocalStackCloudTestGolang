package bucket

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/luke92/LocalStackCloudTestGolang/internal/pkg/cloud"
)

var bucketName string
var filename string
var folder string
var foldertmp string

func Bucket(client cloud.BucketClient) {
	bucketName = "aws-test"
	filename = "id.txt"
	folder = "./assets/"
	foldertmp = "./tmp/"
	ctx := context.Background()

	create(ctx, client)
	uploadObject(ctx, client)
	downloadObject(ctx, client)
	deleteObject(ctx, client)
	listObjects(ctx, client)
}

func create(ctx context.Context, client cloud.BucketClient) {

	existsBucket, err := client.Exists(ctx, bucketName)

	if err != nil {
		log.Println(err.Error())
	}

	if existsBucket {
		log.Println("BucketName " + bucketName + " exists")
	} else {
		if err := client.Create(ctx, bucketName); err != nil {
			log.Fatalln(err)
		}
		log.Println("create: ok")
	}
}

func uploadObject(ctx context.Context, client cloud.BucketClient) {
	file, err := os.Open(folder + filename) // "./assets/id.txt"
	if err != nil {
		log.Fatalln(err)
	}
	//defer is for clean resources
	defer file.Close()

	url, err := client.UploadObject(ctx, bucketName, filename, file)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("upload object:", url)
}

func downloadObject(ctx context.Context, client cloud.BucketClient) {
	createFolderIfNotExists("tmp")
	file, err := os.Create(foldertmp + filename) // "./tmp/id.txt"
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	if err := client.DownloadObject(ctx, bucketName, filename, file); err != nil {
		log.Fatalln(err)
	}
	log.Println("download object: ok")
}

func createFolderIfNotExists(path string) {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(path, os.ModePerm)
		if err != nil {
			log.Println(err)
		}
	}
}

func deleteObject(ctx context.Context, client cloud.BucketClient) {
	if err := client.DeleteObject(ctx, bucketName, filename); err != nil {
		log.Fatalln(err)
	}
	log.Println("delete object: ok")
}

func listObjects(ctx context.Context, client cloud.BucketClient) {
	objects, err := client.ListObjects(ctx, bucketName)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("list objects:")
	for _, object := range objects {
		fmt.Printf("%+v\n", object)
	}
}

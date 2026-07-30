package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"bandplan/src/storage"
)

func main() {
	log.Println("- Testing R2")

	ctx := context.Background()

	r2, err := storage.NewR2Storage(ctx)
	if err != nil {
		panic(err)
	}

	url, err := r2.Upload(
		ctx,
		"test/hello.txt",
		strings.NewReader("Hello form BandPlan!"),
		"text/plain",
	)
	if err != nil {
		panic(err)
	}

	fmt.Println("Upload Successful!")
	fmt.Println(url)
}

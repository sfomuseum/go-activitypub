package main

import (
	"context"
	"flag"
	"log"

	"github.com/sfomuseum/go-activitypub/database"
	"github.com/sfomuseum/go-activitypub/posts"
	"github.com/sfomuseum/go-activitypub/uris"
)

func main() {

	var posts_database_uri string

	flag.StringVar(&posts_database_uri, "posts-database-uri", "", "...")

	flag.Parse()

	ctx := context.Background()

	posts_db, err := database.NewPostsDatabase(ctx, posts_database_uri)

	if err != nil {
		log.Fatalf("Failed to create posts database, %v", err)
	}

	defer posts_db.Close(ctx)

	uris_table := uris.DefaultURIs()

	for _, uri := range flag.Args() {

		post, err := posts.GetPostFromObjectURI(ctx, uris_table, posts_db, uri)

		if err != nil {
			log.Fatalf("Failed to get for post for '%s', %v", uri, err)
		}

		log.Println(uri, post.Id)
	}
}

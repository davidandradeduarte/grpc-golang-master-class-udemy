package main

import (
	"blog/blogpb"
	"context"
	"fmt"
	"io"
	"log"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("running gRPC blog client")

	opts := grpc.WithInsecure()

	cc, err := grpc.Dial("localhost:50051", opts)
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer cc.Close()

	c := blogpb.NewBlogServiceClient(cc)

	// create blog
	fmt.Println("creating the blog...")
	blog := &blogpb.Blog{
		AuthorId: "David",
		Title:    "My First Blog",
		Content:  "Content of the first blog",
	}

	resp, err := c.CreateBlog(context.Background(), &blogpb.CreateBlogRequest{Blog: blog})
	if err != nil {
		log.Fatalf("error creating blog: %v", err)
	}
	blogID := resp.GetBlog().GetId()

	fmt.Printf("blog has been created: %v", resp)

	// read blog
	fmt.Println("reading the blog...")

	_, noBlogErr := c.ReadBlog(context.Background(), &blogpb.ReadBlogRequest{
		BlogId: "something",
	})
	if noBlogErr != nil {
		log.Fatalf("error while calling read blog: %v", noBlogErr)
	}

	readBlogReq := &blogpb.ReadBlogRequest{
		BlogId: blogID,
	}
	readBlogRes, readBlogErr := c.ReadBlog(context.Background(), readBlogReq)
	if readBlogErr != nil {
		log.Fatalf("error while calling read blog: %v", readBlogErr)
	}

	fmt.Printf("blog was read: %v", readBlogRes)

	// update Blog
	newBlog := &blogpb.Blog{
		Id:       blogID,
		AuthorId: "Changed Author",
		Title:    "My First Blog (edited)",
		Content:  "Content of the first blog, with some awesome additions!",
	}
	updateRes, updateErr := c.UpdateBlog(context.Background(), &blogpb.UpdateBlogRequest{Blog: newBlog})
	if updateErr != nil {
		fmt.Printf("Error happened while updating: %v \n", updateErr)
	}
	fmt.Printf("Blog was updated: %v\n", updateRes)

	// delete Blog
	deleteRes, deleteErr := c.DeleteBlog(context.Background(), &blogpb.DeleteBlogRequest{BlogId: blogID})

	if deleteErr != nil {
		fmt.Printf("Error happened while deleting: %v \n", deleteErr)
	}
	fmt.Printf("Blog was deleted: %v \n", deleteRes)

	// list Blogs
	stream, err := c.ListBlog(context.Background(), &blogpb.ListBlogRequest{})
	if err != nil {
		log.Fatalf("error while calling ListBlog RPC: %v", err)
	}
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Something happened: %v", err)
		}
		fmt.Println(res.GetBlog())
	}
}

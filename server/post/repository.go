package post

import (
	"context"
	"fmt"
	"time"

	"github.com/abaltra/blog/server/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Repository struct {
	Config *config.Config
}

var postsCollection *mongo.Collection
var DB *mongo.Client

func (m *Repository) Init() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var err error
	DB, err = mongo.Connect(ctx, options.Client().ApplyURI(m.Config.DBConnectionString))

	if err != nil {
		panic(err)
	}

	postsCollection = DB.Database("blog").Collection("posts")

	if err := m.Ping(); err != nil {
		panic(err)
	}

	fmt.Printf("Connected to Mongo DB at %s\n", m.Config.DBConnectionString)
}

func (m *Repository) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return DB.Ping(ctx, readpref.Primary())
}

func (m *Repository) Create(post Post) (Post, error) {
	fmt.Println("Creating a post")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := postsCollection.InsertOne(ctx, post)

	return post, err
}

func (m *Repository) Save(p Post) error {
	fmt.Println("Updating a post")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := map[string]string{
		"slug": p.Slug,
	}

	_, err := postsCollection.UpdateOne(ctx, filter, p)

	return err
}

func (m *Repository) DeleteBySlug(slug string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	filter := map[string]string{
		"slug": slug,
	}
	_, err := postsCollection.DeleteMany(ctx, filter)

	return err
}

func (m *Repository) DeleteByID(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	filter := map[string]string{
		"id": id,
	}
	_, err := postsCollection.DeleteMany(ctx, filter)

	return err
}

func (m *Repository) List(from int, size int, filters map[string]interface{}) ([]*Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	fmt.Printf("Listing posts. From %d, page size %d\n", from, size)

	query := make(map[string]interface{})

	if filters != nil {
		for key, value := range filters {
			query[key] = value
		}
	}

	_f := int64(from)
	_s := int64(size)
	options := &options.FindOptions{
		Skip:  &_f,
		Limit: &_s,
	}

	results := []*Post{}

	curr, err := postsCollection.Find(ctx, query, options)

	if err != nil {
		return nil, err
	}

	defer curr.Close(ctx)

	for curr.Next(ctx) {
		var result Post
		err := curr.Decode(&result)
		if err != nil {
			return nil, err
		}

		results = append(results, &result)
	}

	return results, nil
}

func (m *Repository) GetBySlug(slug string) (*Post, error) {
	fmt.Printf("Getting post by slug %s\n", slug)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	filter := map[string]string{
		"slug": slug,
	}

	var result Post
	err := postsCollection.FindOne(ctx, filter).Decode(&result)

	return &result, err
}

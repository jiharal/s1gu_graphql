package starwars

import (
	"log"
	"time"

	"github.com/neelance/graphql-go"
)

// Type interface
type (
	character interface {
		ID() graphql.ID
		Name() string
		Friends() *[]*characterResolver
		Height(struct{ Unit string }) float64
	}
)

type Resolver struct{}

// Root Type
type (
	Human struct {
		ID      graphql.ID
		Name    string
		Friends []graphql.ID
		Height  float64
		Post    []graphql.ID
	}
	Post struct {
		ID   graphql.ID
		Date string
		Text string
	}
)

// Type resolver
type (
	humanResolver struct {
		a *Human
	}
	characterResolver struct {
		character
	}
	postResolver struct {
		b *Post
	}
)

var (
	humanData = make(map[graphql.ID]*Human)
	postData  = make(map[graphql.ID]*Post)
)

var humans = []*Human{
	{
		ID:      "1000",
		Name:    "Luke Skywalker",
		Friends: []graphql.ID{"1002", "1003", "2000", "2001"},
		Height:  7.34,
		Post:    []graphql.ID{"1", "3"},
	},
	{
		ID:      "1001",
		Name:    "Muhammad Faris",
		Friends: []graphql.ID{"1001"},
		Height:  1.43,
	},
	{
		ID:      "1002",
		Name:    "Jihar Al Gifari",
		Friends: []graphql.ID{"1001"},
		Height:  1.75,
		Post:    []graphql.ID{"1", "2"},
	},
	{
		ID:      "1003",
		Name:    "Fendi Jatmiko",
		Friends: []graphql.ID{"1002"},
		Height:  1.60,
	},
}

var posts = []*Post{
	{
		ID:   "1",
		Date: "12 Januari 2018",
		Text: "Hello",
	},
	{
		ID:   "2",
		Date: "12 Januari 2018",
		Text: "Hi juga",
	},
	{
		ID:   "3",
		Date: "12 Januari 2018",
		Text: "Apa kabar ?",
	},
}

// root
func (r *Resolver) Human(args struct{ ID graphql.ID }) *humanResolver {
	startTime := time.Now()
	for _, h := range humans {
		humanData[h.ID] = h
	}
	if h := humanData[args.ID]; h != nil {
		endTime := time.Now()
		log.Println("Duration:", endTime.Sub(startTime).Seconds())
		return &humanResolver{h}
	}
	return nil
}

// human resolver
func (s *humanResolver) ID() graphql.ID {
	return s.a.ID
}

func (s *humanResolver) Name() string {
	return s.a.Name
}

func (s *humanResolver) Friends() *[]*characterResolver {
	return resolveCharacters(s.a.Friends)
}
func (s *humanResolver) Height(args struct{ Unit string }) float64 {
	return convertLength(s.a.Height, args.Unit)
}

func (s *humanResolver) Post() *[]*postResolver {
	return resolvePosts(s.a.Post)
}

// resolver posts
func resolvePosts(ids []graphql.ID) *[]*postResolver {
	var posts []*postResolver
	for _, id := range ids {
		if c := resolvePost(id); c != nil {
			posts = append(posts, c)
		}
	}
	return &posts
}
func resolvePost(id graphql.ID) *postResolver {
	for _, h := range posts {
		postData[h.ID] = h
	}
	if h, ok := postData[id]; ok {
		return &postResolver{h}
	}
	return nil
}

// Chracter resolver
func resolveCharacters(ids []graphql.ID) *[]*characterResolver {
	var characters []*characterResolver
	for _, id := range ids {
		if c := resolveCharacter(id); c != nil {
			characters = append(characters, c)
		}
	}
	return &characters
}

func resolveCharacter(id graphql.ID) *characterResolver {
	if h, ok := humanData[id]; ok {
		return &characterResolver{&humanResolver{h}}
	}
	return nil
}

// Post resolver
func (p *postResolver) ID() graphql.ID {
	return p.b.ID
}

func (p *postResolver) Date() string {
	return p.b.Date
}

func (p *postResolver) Text() string {
	return p.b.Text
}

// Other function
func (r *characterResolver) ToHuman() (*humanResolver, bool) {
	c, ok := r.character.(*humanResolver)
	return c, ok
}

func convertLength(meters float64, unit string) float64 {
	switch unit {
	case "METER":
		return meters
	case "FOOT":
		return meters * 3.28084
	default:
		panic("invalid unit")
	}
}

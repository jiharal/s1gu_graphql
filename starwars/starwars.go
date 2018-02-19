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
	}
)

type Resolver struct{}

// Root Type
type (
	Human struct {
		ID      graphql.ID
		Name    string
		Friends []graphql.ID
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
)

var (
	humanData = make(map[graphql.ID]*Human)
)

var humans = []*Human{
	{
		ID:      "1000",
		Name:    "Luke Skywalker",
		Friends: []graphql.ID{"1002", "1003", "2000", "2001"},
	},
	{
		ID:      "1001",
		Name:    "Muhammad Faris",
		Friends: []graphql.ID{"1001"},
	},
	{
		ID:      "1002",
		Name:    "Jihar Al Gifari",
		Friends: []graphql.ID{"1001"},
	},
	{
		ID:      "1003",
		Name:    "Fendi Jatmiko",
		Friends: []graphql.ID{"1002"},
	},
}

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

func (r *characterResolver) ToHuman() (*humanResolver, bool) {
	c, ok := r.character.(*humanResolver)
	return c, ok
}

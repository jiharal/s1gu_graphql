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
		FriendsConnection(friendsConnectionArgs) (*friendsConnectionResolver, error)
		AppearsIn() []string
	}
)

type Resolver struct{}

// Root Type
type (
	Human struct {
		ID        graphql.ID
		Name      string
		Friends   []graphql.ID
		AppearsIn []string
		Height    float64
		Mass      int
		Starships []graphql.ID
	}
	friendsConnectionArgs struct {
		First *int32
		After *graphql.ID
	}
	droid struct {
		ID              graphql.ID
		Name            string
		Friends         []graphql.ID
		AppearsIn       []string
		PrimaryFunction string
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
	friendsConnectionResolver struct {
		ids  []graphql.ID
		from int
		to   int
	}
	droidResolver struct {
		d *droid
	}
)

var (
	humanData = make(map[graphql.ID]*Human)
	droidData = make(map[graphql.ID]*droid)
)

var humans = []*Human{
	{
		ID:        "1000",
		Name:      "Luke Skywalker",
		Friends:   []graphql.ID{"1002", "1003", "2000", "2001"},
		AppearsIn: []string{"NEWHOPE", "EMPIRE", "JEDI"},
		Height:    1.72,
		Mass:      77,
		Starships: []graphql.ID{"3001", "3003"},
	},
	{
		ID:        "1004",
		Name:      "Wilhuff Tarkin",
		Friends:   []graphql.ID{"1001"},
		AppearsIn: []string{"NEWHOPE"},
		Height:    1.8,
		Mass:      0,
	},
}

// root human
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

// foot human
func (s *humanResolver) ID() graphql.ID {
	return s.a.ID
}

func (s *humanResolver) Name() string {
	return s.a.Name
}

func (s *humanResolver) Friends() *[]*characterResolver {
	return resolveCharacters(s.a.Friends)
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
	if _, ok := humanData[id]; ok {
		return &characterResolver{}
	}
	if _, ok := droidData[id]; ok {
		return &characterResolver{}
	}
	return nil
}

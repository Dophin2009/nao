package gqlschema

import (
	"github.com/graphql-go/graphql"
	"gitlab.com/Dophin2009/nao/internal/data"
)

// MediaType is the GraphQL object type for Media.
var MediaType = BuildQueryType(mediaBuilderConfig)

// MediaInputType is the GraphQL input object type for Media.
var MediaInputType = BuildMutationType(mediaBuilderConfig)

var mediaBuilderConfig = TypeBuilderConfig{
	Name: "Media",
	Fields: []FieldBuilderConfig{
		FieldBuilderConfig{
			Name:         "id",
			OutputType:   graphql.NewNonNull(graphql.Int),
			InputType:    graphql.Int,
			Description:  "An integer identifier.",
			DefaultValue: 1,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if md, ok := p.Source.(data.Media); ok {
					return md.ID, nil
				}
				return nil, nil
			},
		},
		FieldBuilderConfig{
			Name:         "titles",
			OutputType:   InfoListType,
			InputType:    InfoInputListType,
			Description:  "A list of titles used to name the Media.",
			DefaultValue: []data.Info{},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if md, ok := p.Source.(data.Media); ok {
					return md.Titles, nil
				}
				return nil, nil
			},
		},
		FieldBuilderConfig{
			Name:         "synopses",
			OutputType:   InfoListType,
			InputType:    InfoInputListType,
			Description:  "A list of synopses, typically in different languages.",
			DefaultValue: []data.Info{},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if md, ok := p.Source.(data.Media); ok {
					return md.Synopses, nil
				}
				return nil, nil
			},
		},
		FieldBuilderConfig{
			Name:         "background",
			OutputType:   InfoListType,
			InputType:    InfoInputListType,
			Description:  "A list of background information segments, typically in different languages.",
			DefaultValue: []data.Info{},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if md, ok := p.Source.(data.Media); ok {
					return md.Background, nil
				}
				return nil, nil
			},
		},
		FieldBuilderConfig{
			Name:         "startDate",
			OutputType:   graphql.DateTime,
			InputType:    graphql.DateTime,
			Description:  "The date the Media began release.",
			DefaultValue: nil,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if md, ok := p.Source.(data.Media); ok {
					return *md.StartDate, nil
				}
				return nil, nil
			},
		},
		FieldBuilderConfig{
			Name:         "endDate",
			OutputType:   graphql.DateTime,
			InputType:    graphql.DateTime,
			Description:  "The date the Media ended release.",
			DefaultValue: nil,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if md, ok := p.Source.(data.Media); ok {
					return *md.EndDate, nil
				}
				return nil, nil
			},
		},
	},
}

var queryMediaByIDField = &graphql.Field{
	Type:        MediaType,
	Description: "Query Media by ID.",
	Args: graphql.FieldConfigArgument{
		"id": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.Int),
		},
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		id, ok := p.Args["id"].(int)
		if !ok {
			return nil, nil
		}

		ds, err := getDataServices(p)
		if err != nil {
			return nil, nil
		}
		mSer := ds.MediaService

		md, err := mSer.GetByID(id)
		if err != nil {
			return nil, err
		}

		return *md, nil
	},
}

var createMediaMutationField = &graphql.Field{
	Type:        MediaType,
	Description: "Create a new Media",
	Args: graphql.FieldConfigArgument{
		"media": &graphql.ArgumentConfig{
			Type: MediaInputType,
		},
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		md, ser, err := serializeMedia(p, "media")
		if err != nil {
			return nil, err
		}

		err = ser.Create(md)
		if err != nil {
			return nil, err
		}

		return *md, err
	},
}

var updateMediaMutationField = &graphql.Field{
	Type:        MediaType,
	Description: "Update a Media",
	Args: graphql.FieldConfigArgument{
		"media": &graphql.ArgumentConfig{
			Type: MediaInputType,
		},
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		md, ser, err := serializeMedia(p, "media")
		if err != nil {
			return nil, err
		}

		err = ser.Update(md)
		if err != nil {
			return nil, err
		}

		return *md, err
	},
}

func serializeMedia(p graphql.ResolveParams, argName string) (*data.Media, *data.MediaService, error) {
	ds, err := getDataServices(p)
	if err != nil {
		return nil, nil, err
	}
	ser := ds.MediaService

	m, err := serializeArg(p.Args[argName], ser)
	if err != nil {
		return nil, nil, err
	}

	md, err := ser.AssertType(m)
	if err != nil {
		return nil, nil, err
	}

	return md, ser, err
}

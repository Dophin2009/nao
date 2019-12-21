package gqlschema

import (
	"github.com/graphql-go/graphql"
	"gitlab.com/Dophin2009/nao/internal/data"
)

// MediaType is the GraphQL object type for Media.
var MediaType = BuildQueryType(mediaBuilderConfig)

// MediaInputType is the GraphQL input object type for Media.
var MediaInputType = BuildMutationType(mediaBuilderConfig)

// SeasonType is the GraphQL object type for Season.
var SeasonType = BuildQueryType(seasonBuilderConfig)

// SeasonInputType is the GraphQL input object type for Season.
var SeasonInputType = BuildMutationType(seasonBuilderConfig)

// QuarterEnumType is the GraphQL enum type for Quarter of
// a Season.
var QuarterEnumType = graphql.NewEnum(graphql.EnumConfig{
	Name: "Quarter",
	Values: graphql.EnumValueConfigMap{
		"Winter": &graphql.EnumValueConfig{
			Value: data.Winter,
		},
		"Spring": &graphql.EnumValueConfig{
			Value: data.Spring,
		},
		"Summer": &graphql.EnumValueConfig{
			Value: data.Summer,
		},
		"Fall": &graphql.EnumValueConfig{
			Value: data.Fall,
		},
	},
})

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
				if md, ok := p.Source.(data.Media); ok && md.StartDate != nil {
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
				if md, ok := p.Source.(data.Media); ok && md.EndDate != nil {
					return *md.EndDate, nil
				}
				return nil, nil
			},
		},
		FieldBuilderConfig{
			Name:        "seasonPremiered",
			OutputType:  SeasonType,
			InputType:   SeasonInputType,
			Description: "The year and season the Media premiered.",
			DefaultValue: data.Season{
				Quarter: nil,
				Year:    nil,
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if md, ok := p.Source.(data.Media); ok {
					return md.SeasonPremiered, nil
				}
				return nil, nil
			},
		},
		FieldBuilderConfig{
			Name:         "type",
			OutputType:   graphql.String,
			InputType:    graphql.String,
			Description:  "The type of the Media.",
			DefaultValue: nil,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if md, ok := p.Source.(data.Media); ok {
					return md.Type, nil
				}
				return nil, nil
			},
		},
		FieldBuilderConfig{
			Name:         "source",
			OutputType:   graphql.String,
			InputType:    graphql.String,
			Description:  "The source material the Media is derived from.",
			DefaultValue: nil,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if md, ok := p.Source.(data.Media); ok {
					return md.Source, nil
				}
				return nil, nil
			},
		},
	},
}

var seasonBuilderConfig = TypeBuilderConfig{
	Name: "Season",
	Fields: []FieldBuilderConfig{
		FieldBuilderConfig{
			Name:        "quarter",
			OutputType:  QuarterEnumType,
			InputType:   QuarterEnumType,
			Description: "A quarter of the year, by season, of the year.",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if sn, ok := p.Source.(data.Season); ok {
					return sn.Quarter, nil
				}
				return nil, nil
			},
		},
		FieldBuilderConfig{
			Name:        "year",
			OutputType:  graphql.Int,
			InputType:   graphql.Int,
			Description: "A single year.",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if sn, ok := p.Source.(data.Season); ok {
					return sn.Year, nil
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

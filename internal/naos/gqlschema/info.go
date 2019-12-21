package gqlschema

import (
	"github.com/graphql-go/graphql"
	"gitlab.com/Dophin2009/nao/internal/data"
)

// InfoListType is the GraphQL object type for a list
// of Info.
var InfoListType = graphql.NewList(graphql.NewNonNull(InfoType))

// InfoType is the GraphQL object type for Info.
var InfoType = BuildQueryType(infoBuilderConfig)

// InfoInputListType is the GraphQL object type for a
// list of InfoInput.
var InfoInputListType = graphql.NewList(graphql.NewNonNull(InfoInputType))

// InfoInputType sit he GraphQL input object type
// for Info.
var InfoInputType = BuildMutationType(infoBuilderConfig)

var infoBuilderConfig = TypeBuilderConfig{
	Name: "Info",
	Fields: []FieldBuilderConfig{
		FieldBuilderConfig{
			Name:         "data",
			OutputType:   graphql.NewNonNull(graphql.String),
			InputType:    graphql.NewNonNull(graphql.String),
			Description:  "The data, typically a text.",
			DefaultValue: "",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if info, ok := p.Source.(data.Info); ok {
					return info.Data, nil
				}
				return nil, nil
			},
		},
		FieldBuilderConfig{
			Name:         "language",
			OutputType:   graphql.NewNonNull(graphql.String),
			InputType:    graphql.NewNonNull(graphql.String),
			Description:  "The langauge the data is in.",
			DefaultValue: "English",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if info, ok := p.Source.(data.Info); ok {
					return info.Language, nil
				}
				return nil, nil
			},
		},
	},
}

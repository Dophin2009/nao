package gqlschema

import (
	"github.com/graphql-go/graphql"
	"gitlab.com/Dophin2009/nao/internal/data"
)

// InfoListType is the GraphQL object type for a list
// of Info.
var InfoListType = graphql.NewNonNull(graphql.NewList(graphql.NewNonNull(InfoType)))

// InfoType is the GraphQL object type for Info.
var InfoType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Info",
	Fields: graphql.Fields{
		"data": &graphql.Field{
			Type:        graphql.NewNonNull(graphql.String),
			Description: "The data, typically text.",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if info, ok := p.Source.(data.Info); ok {
					return info.Data, nil
				}
				return nil, nil
			},
		},
		"language": &graphql.Field{
			Type:        graphql.NewNonNull(graphql.String),
			Description: "The language the data is in.",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if info, ok := p.Source.(data.Info); ok {
					return info.Language, nil
				}
				return nil, nil
			},
		},
	},
})

// InfoInputListType is the GraphQL object type for a
// list of InfoInput.
var InfoInputListType = graphql.NewNonNull(graphql.NewList(graphql.NewNonNull(InfoInputType)))

// InfoInputType sit he GraphQL input object type
// for Info.
var InfoInputType = graphql.NewInputObject(graphql.InputObjectConfig{
	Name: "InfoInput",
	Fields: graphql.InputObjectConfigFieldMap{
		"data": &graphql.InputObjectFieldConfig{
			Type:         graphql.NewNonNull(graphql.String),
			Description:  "The data, typically text.",
			DefaultValue: "",
		},
		"language": &graphql.InputObjectFieldConfig{
			Type:         graphql.NewNonNull(graphql.String),
			Description:  "The language the data is in.",
			DefaultValue: "English",
		},
	},
})

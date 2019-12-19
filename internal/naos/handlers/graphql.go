package handlers

import (
	"errors"

	"github.com/graphql-go/graphql"
	"gitlab.com/Dophin2009/nao/internal/data"
)

// DataServices is a struct passed around in the context
// that contains pointers to all data layer services.
type DataServices struct {
	MediaService *data.MediaService
}

const contextDataServices = "services"

// Schema returns a Schema object for the data.
func Schema() (graphql.Schema, error) {
	// Schema is the GraphQL schema definition for the data.
	return graphql.NewSchema(graphql.SchemaConfig{
		Query: QueryType,
	})
}

// QueryType defines the root query type.
var QueryType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"Media": &graphql.Field{
				Type:        MediaType,
				Description: "Query media by id",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					id, ok := p.Args["id"].(int)
					if !ok {
						return nil, nil
					}

					ds, err := getDataService(p)
					if err != nil {
						return nil, nil
					}
					mSer := ds.MediaService
					return mSer.GetByID(id)
				},
			},
		},
	},
)

// MediaType is the GraphQL object type for Media.
var MediaType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Media",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.Int,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if md, ok := p.Source.(data.Media); ok {
						return md.ID, nil
					}
					return nil, nil
				},
			},
			"titles": &graphql.Field{
				Type: graphql.NewList(InfoType),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if md, ok := p.Source.(data.Media); ok {
						return md.Titles, nil
					}
					return nil, nil
				},
			},
			"synopses": &graphql.Field{
				Type: graphql.NewList(InfoType),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if md, ok := p.Source.(data.Media); ok {
						return md.Synopses, nil
					}
					return nil, nil
				},
			},
		},
	},
)

// InfoType is the GraphQL object type for Info.
var InfoType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Info",
		Fields: graphql.Fields{
			"data": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if info, ok := p.Source.(data.Info); ok {
						return info.Data, nil
					}
					return nil, nil
				},
			},
			"language": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if info, ok := p.Source.(data.Info); ok {
						return info.Language, nil
					}
					return nil, nil
				},
			},
		},
	},
)

func getDataService(p graphql.ResolveParams) (*DataServices, error) {
	v := p.Context.Value(contextDataServices)
	if v == nil {
		return nil, errors.New("data services not found")
	}
	ds, ok := v.(*DataServices)
	if !ok {
		return nil, errors.New("found value is not a DataServices")
	}
	return ds, nil
}

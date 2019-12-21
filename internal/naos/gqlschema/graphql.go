package gqlschema

import (
	"errors"

	"github.com/graphql-go/graphql"
	json "github.com/json-iterator/go"
	"gitlab.com/Dophin2009/nao/internal/data"
)

// DataServices is a struct passed around in the context
// that contains pointers to all data layer services.
type DataServices struct {
	MediaService *data.MediaService
}

// ContextDataServices is the context key for DataServices.
const ContextDataServices = "services"

// Schema returns a Schema object for the data.
func Schema() (graphql.Schema, error) {
	return graphql.NewSchema(graphql.SchemaConfig{
		Query:    QueryType,
		Mutation: MutationType,
	})
}

// QueryType defines the root query type.
var QueryType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Query",
	Fields: graphql.Fields{
		"QueryMediaByID": queryMediaByIDField,
	},
})

// MutationType defines the root mutation type.
var MutationType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Mutation",
	Fields: graphql.Fields{
		"CreateMedia": createMediaMutationField,
		"UpdateMedia": updateMediaMutationField,
	},
})

// TypeBuilderConfig contains information to create
// output and input GraphQL object types.
type TypeBuilderConfig struct {
	Name   string
	Fields []FieldBuilderConfig
}

// FieldBuilderConfig contains information to create
// output and input GraphQL fields.
type FieldBuilderConfig struct {
	Name         string
	Description  string
	OutputType   graphql.Type
	InputType    graphql.Type
	Resolve      func(p graphql.ResolveParams) (interface{}, error)
	DefaultValue interface{}
}

// BuildQueryType returns an output object type and
// an input object type for the given type information.
func BuildQueryType(config TypeBuilderConfig) *graphql.Object {
	qfields := make(graphql.Fields)
	for _, f := range config.Fields {
		qfields[f.Name] = &graphql.Field{
			Type:        f.OutputType,
			Description: f.Description,
			Resolve:     f.Resolve,
		}
	}
	qtype := graphql.NewObject(graphql.ObjectConfig{
		Name:   config.Name,
		Fields: qfields,
	})

	return qtype
}

// BuildMutationType returns an output object type and
// an input object type for the given type information.
func BuildMutationType(config TypeBuilderConfig) *graphql.InputObject {
	mfields := make(graphql.InputObjectConfigFieldMap)
	for _, f := range config.Fields {
		mfields[f.Name] = &graphql.InputObjectFieldConfig{
			Type:         f.InputType,
			Description:  f.Description,
			DefaultValue: f.DefaultValue,
		}
	}
	mtype := graphql.NewInputObject(graphql.InputObjectConfig{
		Name:   config.Name + "Input",
		Fields: mfields,
	})

	return mtype
}

func serializeArg(arg interface{}, ser data.Service) (data.Model, error) {
	jsonString, err := json.Marshal(arg)
	if err != nil {
		return nil, err
	}

	m, err := ser.Unmarshal([]byte(jsonString))
	if err != nil {
		return nil, err
	}

	return m, err
}

func getDataServices(p graphql.ResolveParams) (*DataServices, error) {
	v := p.Context.Value(ContextDataServices)
	if v == nil {
		return nil, errors.New("data services not found")
	}
	ds, ok := v.(*DataServices)
	if !ok {
		return nil, errors.New("found value is not a DataServices")
	}
	return ds, nil
}

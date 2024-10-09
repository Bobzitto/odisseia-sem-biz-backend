package graph

import (
	"backend/internal/models"
	"errors"
	"strings"

	"github.com/graphql-go/graphql"
)

type Graph struct {
	Aulas       []*models.Aula
	QueryString string
	Config      graphql.SchemaConfig
	fields      graphql.Fields
	aulaType    *graphql.Object
}

func New(aulas []*models.Aula) *Graph {
	var aulaType = graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Aula",
			Fields: graphql.Fields{
				"id": &graphql.Field{
					Type: graphql.Int,
				},
				"name": &graphql.Field{
					Type: graphql.String,
				},
				"active": &graphql.Field{
					Type: graphql.Boolean,
				},
				"size": &graphql.Field{
					Type: graphql.Int,
				},
				"review": &graphql.Field{
					Type: graphql.Float,
				},
				"created_at": &graphql.Field{
					Type: graphql.DateTime,
				},
				"updated_at": &graphql.Field{
					Type: graphql.DateTime,
				},
			},
		},
	)
	var fields = graphql.Fields{

		"list": &graphql.Field{
			Type:        graphql.NewList(aulaType),
			Description: "get all aulas",
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				return aulas, nil
			},
		},

		"search": &graphql.Field{
			Type:        graphql.NewList(aulaType),
			Description: "search aulas by name",
			Args: graphql.FieldConfigArgument{
				"nameContains": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				var list []*models.Aula
				search, ok := params.Args["nameContains"].(string)
				if ok {
					for _, currentAula := range aulas {
						if strings.Contains(strings.ToLower(currentAula.Name), strings.ToLower(search)) {
							list = append(list, currentAula)
						}
					}
				}
				return list, nil
			},
		},

		"get": &graphql.Field{
			Type:        aulaType,
			Description: "get aula by id",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.Int,
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				id, ok := p.Args["id"].(int)
				if ok {
					for _, aula := range aulas {
						if aula.ID == id {
							return aula, nil
						}
					}
				}
				return nil, nil
			},
		},
	}

	return &Graph{
		Aulas:    aulas,
		fields:   fields,
		aulaType: aulaType,
	}
}

func (g *Graph) Query() (*graphql.Result, error){
	rootQuery := graphql.ObjectConfig{
		Name: "RootQuery",
		Fields: g.fields,
	}
	schemaConfig := graphql.SchemaConfig{Query: graphql.NewObject(rootQuery)}
	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		return nil, err
	}

	params := graphql.Params{Schema: schema, RequestString: g.QueryString}
	resp := graphql.Do(params)
	if len(resp.Errors) > 0 {
		return nil, errors.New("error executing query")
	}
	return resp, nil
}
package graphqlschema_test

import (
	"fmt"

	"github.com/OpenUdon/graphqlschema"
)

func ExampleParseSDL() {
	model, err := graphqlschema.ParseSDL([]byte(`schema { query: Query }
type Query {
  book(id: ID!): Book
}
type Book {
  id: ID!
  title: String!
}`))
	if err != nil {
		panic(err)
	}
	fmt.Println(model.QueryType)
	fmt.Println(len(model.Operations))
	// Output:
	// Query
	// 1
}

func ExampleParseIntrospection() {
	model, err := graphqlschema.ParseIntrospection([]byte(`{
  "data": {
    "__schema": {
      "queryType": {"name": "Query"},
      "types": [
        {
          "kind": "OBJECT",
          "name": "Query",
          "fields": [
            {
              "name": "book",
              "args": [
                {
                  "name": "id",
                  "type": {"kind": "NON_NULL", "ofType": {"kind": "SCALAR", "name": "ID"}}
                }
              ],
              "type": {"kind": "OBJECT", "name": "Book"}
            }
          ]
        },
        {"kind": "OBJECT", "name": "Book", "fields": []}
      ]
    }
  }
}`))
	if err != nil {
		panic(err)
	}
	op, _ := model.OperationByID("query.book")
	fmt.Println(op.ID)
	fmt.Println(op.Arguments[0].Type)
	// Output:
	// query.book
	// ID!
}

func ExampleModel_OperationByID() {
	model, err := graphqlschema.ParseSDL([]byte(`type Query { search(term: String): [Book!]! }
type Book { title: String! }`))
	if err != nil {
		panic(err)
	}
	op, ok := model.OperationByID("query.search")
	fmt.Println(ok)
	fmt.Println(op.Type)
	// Output:
	// true
	// [Book!]!
}

func ExampleModel_ResolveSelector() {
	model, err := graphqlschema.ParseSDL([]byte(`type Query { search(term: String): [Book!]! }
type Book { title: String! }`))
	if err != nil {
		panic(err)
	}
	target, ok := model.ResolveSelector("#/types/Query/fields/search")
	fmt.Println(ok)
	fmt.Println(target.Operation.ID)
	// Output:
	// true
	// query.search
}

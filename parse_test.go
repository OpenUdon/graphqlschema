package graphqlschema_test

import (
	"os"
	"strings"
	"testing"

	"github.com/OpenUdon/graphqlschema"
)

func TestParseSDLFixture(t *testing.T) {
	model := parseSDLFixture(t)
	if model.SourceKind != graphqlschema.SourceKindSDL {
		t.Fatalf("source kind = %q", model.SourceKind)
	}
	if got, want := model.Description, "Library API schema."; got != want {
		t.Fatalf("description = %q, want %q", got, want)
	}
	assertRoots(t, model)

	book, ok := model.TypeByName("Book")
	if !ok {
		t.Fatal("Book type missing")
	}
	if book.Kind != "OBJECT" || len(book.Interfaces) != 1 || book.Interfaces[0] != "Node" {
		t.Fatalf("Book = %#v", book)
	}
	if len(book.Directives) != 1 || book.Directives[0].Arguments["name"] != `"catalog"` {
		t.Fatalf("Book directives = %#v", book.Directives)
	}
	statusField := fieldByName(book.Fields, "status")
	if statusField == nil || statusField.Type.String() != "BookStatus!" {
		t.Fatalf("Book.status = %#v", statusField)
	}

	query, _ := model.TypeByName("Query")
	bookField := fieldByName(query.Fields, "book")
	if bookField == nil || bookField.Description != "Find a book by id." {
		t.Fatalf("Query.book = %#v", bookField)
	}
	if len(bookField.Arguments) != 1 || bookField.Arguments[0].Type.String() != "ID!" {
		t.Fatalf("Query.book args = %#v", bookField.Arguments)
	}
	if got := bookField.Arguments[0].Directives[0].Arguments["name"]; got != `"lookup"` {
		t.Fatalf("argument directive = %q", got)
	}
	searchField := fieldByName(query.Fields, "search")
	if searchField == nil || searchField.Type.String() != "[SearchResult!]!" {
		t.Fatalf("Query.search = %#v", searchField)
	}
	node, _ := model.TypeByName("Node")
	if len(node.PossibleTypes) != 1 || node.PossibleTypes[0] != "Book" {
		t.Fatalf("Node possible types = %#v", node.PossibleTypes)
	}
	searchResult, _ := model.TypeByName("SearchResult")
	if len(searchResult.PossibleTypes) != 1 || searchResult.PossibleTypes[0] != "Book" {
		t.Fatalf("SearchResult possible types = %#v", searchResult.PossibleTypes)
	}

	filter, _ := model.TypeByName("BookFilter")
	statusInput := inputByName(filter.InputFields, "status")
	if statusInput == nil || statusInput.DefaultValue != "AVAILABLE" {
		t.Fatalf("BookFilter.status = %#v", statusInput)
	}

	enum, _ := model.TypeByName("BookStatus")
	checkedOut := enumValueByName(enum.EnumValues, "CHECKED_OUT")
	if checkedOut == nil || checkedOut.Directives[0].Arguments["name"] != `"busy"` {
		t.Fatalf("enum value = %#v", checkedOut)
	}

	tag := model.Directives["tag"]
	if tag == nil || !tag.Repeatable || len(tag.Locations) == 0 || tag.Arguments[0].Type.String() != "String!" {
		t.Fatalf("tag directive = %#v", tag)
	}

	if _, ok := model.TypeByName("__Schema"); ok {
		t.Fatal("introspection type should not be exposed from SDL parser")
	}
	if scalar, ok := model.TypeByName("String"); !ok || !scalar.BuiltIn {
		t.Fatalf("built-in scalar String missing or not marked built-in: %#v", scalar)
	}
	assertOperationsAndSelectors(t, model)
}

func TestParseIntrospectionFixture(t *testing.T) {
	data := readFixture(t, "testdata/library-introspection.json")
	model, err := graphqlschema.ParseIntrospection(data)
	if err != nil {
		t.Fatalf("ParseIntrospection failed: %v", err)
	}
	if model.SourceKind != graphqlschema.SourceKindIntrospection {
		t.Fatalf("source kind = %q", model.SourceKind)
	}
	if model.Raw == nil {
		t.Fatal("Raw should preserve introspection document")
	}
	assertRoots(t, model)

	query, _ := model.TypeByName("Query")
	search := fieldByName(query.Fields, "search")
	if search == nil || search.Type.String() != "[SearchResult!]!" {
		t.Fatalf("Query.search = %#v", search)
	}
	if got := search.Arguments[0].Type.String(); got != "BookFilter" {
		t.Fatalf("search filter type = %q", got)
	}

	book, _ := model.TypeByName("Book")
	if len(book.Interfaces) != 1 || book.Interfaces[0] != "Node" {
		t.Fatalf("Book interfaces = %#v", book.Interfaces)
	}
	node, _ := model.TypeByName("Node")
	if len(node.PossibleTypes) != 1 || node.PossibleTypes[0] != "Book" {
		t.Fatalf("Node possible types = %#v", node.PossibleTypes)
	}
	filter, _ := model.TypeByName("BookFilter")
	if got := inputByName(filter.InputFields, "limit").DefaultValue; got != "10" {
		t.Fatalf("limit default = %q", got)
	}
	if tag := model.Directives["tag"]; tag == nil || !tag.Repeatable || tag.Arguments[0].DefaultValue != `"default"` {
		t.Fatalf("tag directive = %#v", tag)
	}
	assertOperationsAndSelectors(t, model)
}

func TestParseAutoDetectsSourceKind(t *testing.T) {
	if model, err := graphqlschema.Parse(readFixture(t, "testdata/library.graphqls")); err != nil || model.SourceKind != graphqlschema.SourceKindSDL {
		t.Fatalf("SDL Parse = %#v, %v", model, err)
	}
	if model, err := graphqlschema.Parse(readFixture(t, "testdata/library-introspection.json")); err != nil || model.SourceKind != graphqlschema.SourceKindIntrospection {
		t.Fatalf("introspection Parse = %#v, %v", model, err)
	}
}

func TestParseIntrospectionMapRootForms(t *testing.T) {
	model, err := graphqlschema.ParseIntrospectionMap(map[string]any{
		"__schema": map[string]any{
			"queryType": map[string]any{"name": "Query"},
			"types": []any{
				map[string]any{"kind": "OBJECT", "name": "Query", "fields": []any{}},
			},
		},
	})
	if err != nil {
		t.Fatalf("ParseIntrospectionMap failed: %v", err)
	}
	if model.QueryType != "Query" {
		t.Fatalf("query root = %q", model.QueryType)
	}
}

func TestTypeByNameMisses(t *testing.T) {
	model := parseSDLFixture(t)
	if _, ok := model.TypeByName(""); ok {
		t.Fatal("empty type unexpectedly resolved")
	}
	if _, ok := model.TypeByName("Missing"); ok {
		t.Fatal("missing type unexpectedly resolved")
	}
	var nilModel *graphqlschema.Model
	if _, ok := nilModel.TypeByName("Query"); ok {
		t.Fatal("nil model unexpectedly resolved")
	}
}

func TestOperationByIDMisses(t *testing.T) {
	model := parseSDLFixture(t)
	if _, ok := model.OperationByID(""); ok {
		t.Fatal("empty operation unexpectedly resolved")
	}
	if _, ok := model.OperationByID("query.missing"); ok {
		t.Fatal("missing operation unexpectedly resolved")
	}
	var nilModel *graphqlschema.Model
	if _, ok := nilModel.OperationByID("query.book"); ok {
		t.Fatal("nil model unexpectedly resolved operation")
	}
}

func TestMalformedInputs(t *testing.T) {
	cases := []struct {
		name string
		run  func() error
		want string
	}{
		{
			name: "empty parse",
			run: func() error {
				_, err := graphqlschema.Parse(nil)
				return err
			},
			want: "empty",
		},
		{
			name: "invalid SDL",
			run: func() error {
				_, err := graphqlschema.ParseSDL([]byte(`type Query { broken(: String): String }`))
				return err
			},
			want: "parse graphql SDL",
		},
		{
			name: "invalid JSON",
			run: func() error {
				_, err := graphqlschema.ParseIntrospection([]byte(`{`))
				return err
			},
			want: "parse graphql introspection JSON",
		},
		{
			name: "trailing JSON",
			run: func() error {
				_, err := graphqlschema.ParseIntrospection([]byte(`{"__schema":{"types":[]}} {"__schema":{"types":[]}}`))
				return err
			},
			want: "trailing data",
		},
		{
			name: "missing schema",
			run: func() error {
				_, err := graphqlschema.ParseIntrospection([]byte(`{"data":{}}`))
				return err
			},
			want: "missing __schema",
		},
		{
			name: "top level errors",
			run: func() error {
				_, err := graphqlschema.ParseIntrospection([]byte(`{"errors":[{"message":"nope"}],"data":{"__schema":{"types":[]}}}`))
				return err
			},
			want: "contains errors",
		},
		{
			name: "malformed types",
			run: func() error {
				_, err := graphqlschema.ParseIntrospection([]byte(`{"__schema":{"types":{}}}`))
				return err
			},
			want: "types must be a list",
		},
		{
			name: "malformed fields",
			run: func() error {
				_, err := graphqlschema.ParseIntrospection([]byte(`{"__schema":{"types":[{"kind":"OBJECT","name":"Query","fields":{}}]}}`))
				return err
			},
			want: "fields must be a list",
		},
		{
			name: "malformed type ref",
			run: func() error {
				_, err := graphqlschema.ParseIntrospection([]byte(`{"__schema":{"types":[{"kind":"OBJECT","name":"Query","fields":[{"name":"x","type":{"kind":"NON_NULL","ofType":null}}]}]}}`))
				return err
			},
			want: "ofType must be an object for NON_NULL",
		},
		{
			name: "unknown type ref kind",
			run: func() error {
				_, err := graphqlschema.ParseIntrospection([]byte(`{"__schema":{"types":[{"kind":"OBJECT","name":"Query","fields":[{"name":"x","type":{"kind":"BOGUS","name":"Thing"}}]}]}}`))
				return err
			},
			want: "not a supported GraphQL type-ref kind",
		},
		{
			name: "missing type ref kind",
			run: func() error {
				_, err := graphqlschema.ParseIntrospection([]byte(`{"__schema":{"types":[{"kind":"OBJECT","name":"Query","fields":[{"name":"x","type":{"name":"Thing"}}]}]}}`))
				return err
			},
			want: "kind must be a string",
		},
		{
			name: "malformed directive",
			run: func() error {
				_, err := graphqlschema.ParseIntrospection([]byte(`{"__schema":{"types":[],"directives":[{"name":"bad","locations":[1]}]}}`))
				return err
			},
			want: "locations[0] must be a string",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.run()
			if err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("error = %v, want %q", err, tc.want)
			}
		})
	}
}

func assertOperationsAndSelectors(t *testing.T, model *graphqlschema.Model) {
	t.Helper()
	got := make([]string, 0, len(model.Operations))
	for _, op := range model.Operations {
		got = append(got, op.ID)
	}
	want := []string{"query.book", "query.search", "mutation.checkout", "subscription.bookUpdated"}
	if len(got) != len(want) {
		t.Fatalf("operations = %#v, want %#v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("operations = %#v, want %#v", got, want)
		}
	}
	book, ok := model.OperationByID("query.book")
	if !ok {
		t.Fatal("query.book operation missing")
	}
	if book.Kind != "query" || book.RootType != "Query" || book.FieldName != "book" || book.Type.String() != "Book" {
		t.Fatalf("query.book = %#v", book)
	}
	if len(book.Arguments) != 1 || book.Arguments[0].Name != "id" || book.Arguments[0].Type.String() != "ID!" {
		t.Fatalf("query.book args = %#v", book.Arguments)
	}
	for _, selector := range []string{
		"query.book",
		"#/operations/query.book",
		"#/types/Query/fields/book",
		"#/types/Mutation/fields/checkout",
		"#/types/Subscription/fields/bookUpdated",
	} {
		target, ok := model.ResolveSelector(selector)
		if !ok {
			t.Fatalf("selector %q did not resolve", selector)
		}
		if target.Operation == nil || target.Field == nil || target.RootType == nil {
			t.Fatalf("selector %q target incomplete: %#v", selector, target)
		}
		if !model.SelectorAliases()[selector] {
			t.Fatalf("selector %q missing from aliases %#v", selector, model.SelectorAliases())
		}
	}
	if target, ok := model.ResolveSelector("#/operations/query%2Ebook"); !ok || target.Operation.ID != "query.book" {
		t.Fatalf("percent-encoded operation selector target = %#v, %v", target, ok)
	}
	for _, selector := range []string{
		"query.missing",
		"#/operations/query.missing",
		"#/types/Query/fields/missing",
		"#/types/Book/fields/title",
		"#/components/schemas/Book",
	} {
		if target, ok := model.ResolveSelector(selector); ok {
			t.Fatalf("selector %q unexpectedly resolved to %#v", selector, target)
		}
	}
	var nilModel *graphqlschema.Model
	if len(nilModel.SelectorAliases()) != 0 {
		t.Fatal("nil model returned selector aliases")
	}
	if target, ok := nilModel.ResolveSelector("query.book"); ok {
		t.Fatalf("nil model resolved selector to %#v", target)
	}
}

func parseSDLFixture(t *testing.T) *graphqlschema.Model {
	t.Helper()
	model, err := graphqlschema.ParseSDL(readFixture(t, "testdata/library.graphqls"))
	if err != nil {
		t.Fatalf("ParseSDL failed: %v", err)
	}
	return model
}

func readFixture(t *testing.T, path string) []byte {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func assertRoots(t *testing.T, model *graphqlschema.Model) {
	t.Helper()
	if model.QueryType != "Query" || model.MutationType != "Mutation" || model.SubscriptionType != "Subscription" {
		t.Fatalf("roots = query:%q mutation:%q subscription:%q", model.QueryType, model.MutationType, model.SubscriptionType)
	}
}

func fieldByName(fields []*graphqlschema.FieldDefinition, name string) *graphqlschema.FieldDefinition {
	for _, field := range fields {
		if field.Name == name {
			return field
		}
	}
	return nil
}

func inputByName(inputs []*graphqlschema.InputValueDefinition, name string) *graphqlschema.InputValueDefinition {
	for _, input := range inputs {
		if input.Name == name {
			return input
		}
	}
	return nil
}

func enumValueByName(values []*graphqlschema.EnumValueDefinition, name string) *graphqlschema.EnumValueDefinition {
	for _, value := range values {
		if value.Name == name {
			return value
		}
	}
	return nil
}

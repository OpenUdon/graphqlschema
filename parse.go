package graphqlschema

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/vektah/gqlparser/v2"
	"github.com/vektah/gqlparser/v2/ast"
)

// Parse parses GraphQL SDL or introspection JSON into native metadata.
func Parse(data []byte) (*Model, error) {
	trimmed := bytes.TrimSpace(data)
	if len(trimmed) == 0 {
		return nil, fmt.Errorf("graphql schema document is empty")
	}
	if trimmed[0] == '{' {
		return ParseIntrospection(trimmed)
	}
	return ParseSDL(trimmed)
}

// ParseSDL parses GraphQL schema definition language into native metadata.
func ParseSDL(data []byte) (*Model, error) {
	if len(bytes.TrimSpace(data)) == 0 {
		return nil, fmt.Errorf("graphql SDL document is empty")
	}
	schema, err := gqlparser.LoadSchema(&ast.Source{Name: "schema.graphql", Input: string(data)})
	if err != nil {
		return nil, fmt.Errorf("parse graphql SDL: %w", err)
	}
	return modelFromAST(schema), nil
}

// ParseIntrospection parses a GraphQL introspection JSON response into native
// metadata.
func ParseIntrospection(data []byte) (*Model, error) {
	var raw map[string]any
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.UseNumber()
	if err := dec.Decode(&raw); err != nil {
		return nil, fmt.Errorf("parse graphql introspection JSON: %w", err)
	}
	var trailing any
	if err := dec.Decode(&trailing); err == nil {
		return nil, fmt.Errorf("parse graphql introspection JSON: trailing data after root object")
	} else if err != io.EOF {
		return nil, fmt.Errorf("parse graphql introspection JSON: %w", err)
	}
	if raw == nil {
		return nil, fmt.Errorf("graphql introspection root must be an object")
	}
	return ParseIntrospectionMap(raw)
}

// ParseIntrospectionMap parses an already-decoded GraphQL introspection JSON
// response into native metadata.
func ParseIntrospectionMap(raw map[string]any) (*Model, error) {
	if raw == nil {
		return nil, fmt.Errorf("graphql introspection root must be an object")
	}
	if hasGraphQLErrors(raw["errors"]) {
		return nil, fmt.Errorf("graphql introspection response contains errors")
	}
	schemaMap, err := introspectionSchemaMap(raw)
	if err != nil {
		return nil, err
	}
	model := &Model{
		SourceKind: SourceKindIntrospection,
		Types:      map[string]*TypeDefinition{},
		Directives: map[string]*DirectiveDefinition{},
		Raw:        raw,
	}
	model.Description = stringValue(schemaMap["description"])
	model.QueryType, err = namedTypeField(schemaMap, "queryType", "graphql introspection.__schema")
	if err != nil {
		return nil, err
	}
	model.MutationType, err = namedTypeField(schemaMap, "mutationType", "graphql introspection.__schema")
	if err != nil {
		return nil, err
	}
	model.SubscriptionType, err = namedTypeField(schemaMap, "subscriptionType", "graphql introspection.__schema")
	if err != nil {
		return nil, err
	}
	if err := parseIntrospectionTypes(model, schemaMap); err != nil {
		return nil, err
	}
	if err := parseIntrospectionDirectives(model, schemaMap); err != nil {
		return nil, err
	}
	model.buildOperations()
	return model, nil
}

func modelFromAST(schema *ast.Schema) *Model {
	model := &Model{
		SourceKind:  SourceKindSDL,
		Description: schema.Description,
		Types:       map[string]*TypeDefinition{},
		Directives:  map[string]*DirectiveDefinition{},
	}
	if schema.Query != nil {
		model.QueryType = schema.Query.Name
	}
	if schema.Mutation != nil {
		model.MutationType = schema.Mutation.Name
	}
	if schema.Subscription != nil {
		model.SubscriptionType = schema.Subscription.Name
	}
	for _, name := range sortedASTTypeNames(schema.Types) {
		def := schema.Types[name]
		if def == nil || def.Name == "" {
			continue
		}
		if def.BuiltIn && (def.Kind != ast.Scalar || strings.HasPrefix(def.Name, "__")) {
			continue
		}
		model.Types[def.Name] = typeDefinitionFromAST(schema, def)
	}
	for _, name := range sortedASTDirectiveNames(schema.Directives) {
		def := schema.Directives[name]
		if def == nil || def.Name == "" {
			continue
		}
		model.Directives[def.Name] = directiveDefinitionFromAST(def)
	}
	model.buildOperations()
	return model
}

func typeDefinitionFromAST(schema *ast.Schema, def *ast.Definition) *TypeDefinition {
	out := &TypeDefinition{
		Kind:          string(def.Kind),
		Name:          def.Name,
		Description:   def.Description,
		Interfaces:    append([]string(nil), def.Interfaces...),
		PossibleTypes: possibleTypesFromAST(schema, def),
		Directives:    directiveUsesFromAST(def.Directives),
		BuiltIn:       def.BuiltIn,
	}
	for _, field := range def.Fields {
		if field == nil {
			continue
		}
		if def.Kind == ast.InputObject {
			out.InputFields = append(out.InputFields, inputValueFromASTField(field))
			continue
		}
		out.Fields = append(out.Fields, fieldDefinitionFromAST(field))
	}
	for _, value := range def.EnumValues {
		if value == nil {
			continue
		}
		out.EnumValues = append(out.EnumValues, &EnumValueDefinition{
			Name:        value.Name,
			Description: value.Description,
			Directives:  directiveUsesFromAST(value.Directives),
		})
	}
	return out
}

func possibleTypesFromAST(schema *ast.Schema, def *ast.Definition) []string {
	if schema == nil || def == nil {
		return nil
	}
	seen := map[string]bool{}
	for _, name := range def.Types {
		if strings.TrimSpace(name) != "" {
			seen[name] = true
		}
	}
	for _, possible := range schema.PossibleTypes[def.Name] {
		if possible != nil && strings.TrimSpace(possible.Name) != "" {
			seen[possible.Name] = true
		}
	}
	if len(seen) == 0 {
		return nil
	}
	out := make([]string, 0, len(seen))
	for name := range seen {
		out = append(out, name)
	}
	sort.Strings(out)
	return out
}

func fieldDefinitionFromAST(field *ast.FieldDefinition) *FieldDefinition {
	out := &FieldDefinition{
		Name:        field.Name,
		Description: field.Description,
		Type:        typeRefFromAST(field.Type),
		Directives:  directiveUsesFromAST(field.Directives),
	}
	for _, arg := range field.Arguments {
		if arg == nil {
			continue
		}
		out.Arguments = append(out.Arguments, inputValueFromASTArgument(arg))
	}
	return out
}

func inputValueFromASTArgument(arg *ast.ArgumentDefinition) *InputValueDefinition {
	out := &InputValueDefinition{
		Name:        arg.Name,
		Description: arg.Description,
		Type:        typeRefFromAST(arg.Type),
		Directives:  directiveUsesFromAST(arg.Directives),
	}
	if arg.DefaultValue != nil {
		out.DefaultValue = arg.DefaultValue.String()
	}
	return out
}

func inputValueFromASTField(field *ast.FieldDefinition) *InputValueDefinition {
	out := &InputValueDefinition{
		Name:        field.Name,
		Description: field.Description,
		Type:        typeRefFromAST(field.Type),
		Directives:  directiveUsesFromAST(field.Directives),
	}
	if field.DefaultValue != nil {
		out.DefaultValue = field.DefaultValue.String()
	}
	return out
}

func directiveDefinitionFromAST(def *ast.DirectiveDefinition) *DirectiveDefinition {
	out := &DirectiveDefinition{
		Name:        def.Name,
		Description: def.Description,
		Repeatable:  def.IsRepeatable,
	}
	for _, location := range def.Locations {
		out.Locations = append(out.Locations, string(location))
	}
	sort.Strings(out.Locations)
	for _, arg := range def.Arguments {
		if arg == nil {
			continue
		}
		out.Arguments = append(out.Arguments, inputValueFromASTArgument(arg))
	}
	return out
}

func directiveUsesFromAST(list ast.DirectiveList) []*DirectiveUse {
	var out []*DirectiveUse
	for _, directive := range list {
		if directive == nil || directive.Name == "" {
			continue
		}
		use := &DirectiveUse{Name: directive.Name}
		for _, arg := range directive.Arguments {
			if arg == nil || arg.Name == "" || arg.Value == nil {
				continue
			}
			if use.Arguments == nil {
				use.Arguments = map[string]string{}
			}
			use.Arguments[arg.Name] = arg.Value.String()
		}
		out = append(out, use)
	}
	return out
}

func typeRefFromAST(t *ast.Type) *TypeRef {
	if t == nil {
		return nil
	}
	return &TypeRef{
		NamedType: t.NamedType,
		Elem:      typeRefFromAST(t.Elem),
		NonNull:   t.NonNull,
	}
}

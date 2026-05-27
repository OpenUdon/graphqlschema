package graphqlschema

import "strings"

const (
	// SourceKindSDL marks models parsed from GraphQL schema definition language.
	SourceKindSDL = "sdl"
	// SourceKindIntrospection marks models parsed from GraphQL introspection JSON.
	SourceKindIntrospection = "introspection"
)

// Model is a parsed GraphQL schema artifact.
type Model struct {
	SourceKind       string
	Description      string
	QueryType        string
	MutationType     string
	SubscriptionType string
	Types            map[string]*TypeDefinition
	Directives       map[string]*DirectiveDefinition
	Raw              map[string]any
}

// TypeDefinition describes a GraphQL type.
type TypeDefinition struct {
	Kind          string
	Name          string
	Description   string
	Fields        []*FieldDefinition
	InputFields   []*InputValueDefinition
	Interfaces    []string
	PossibleTypes []string
	EnumValues    []*EnumValueDefinition
	Directives    []*DirectiveUse
	BuiltIn       bool
}

// FieldDefinition describes a field on an object or interface type.
type FieldDefinition struct {
	Name        string
	Description string
	Type        *TypeRef
	Arguments   []*InputValueDefinition
	Directives  []*DirectiveUse
}

// InputValueDefinition describes a field argument, input-object field, or
// directive argument.
type InputValueDefinition struct {
	Name         string
	Description  string
	Type         *TypeRef
	DefaultValue string
	Directives   []*DirectiveUse
}

// EnumValueDefinition describes a GraphQL enum value.
type EnumValueDefinition struct {
	Name        string
	Description string
	Directives  []*DirectiveUse
}

// DirectiveDefinition describes a GraphQL directive definition.
type DirectiveDefinition struct {
	Name        string
	Description string
	Locations   []string
	Arguments   []*InputValueDefinition
	Repeatable  bool
}

// DirectiveUse describes a directive applied to schema metadata.
type DirectiveUse struct {
	Name      string
	Arguments map[string]string
}

// TypeRef describes a GraphQL named, list, or non-null type reference.
type TypeRef struct {
	NamedType string
	Elem      *TypeRef
	NonNull   bool
}

// String returns the GraphQL type-reference spelling.
func (t *TypeRef) String() string {
	if t == nil {
		return ""
	}
	suffix := ""
	if t.NonNull {
		suffix = "!"
	}
	if t.NamedType != "" {
		return t.NamedType + suffix
	}
	return "[" + t.Elem.String() + "]" + suffix
}

// TypeByName returns a parsed type by GraphQL type name.
func (m *Model) TypeByName(name string) (*TypeDefinition, bool) {
	if m == nil {
		return nil, false
	}
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, false
	}
	typ, ok := m.Types[name]
	return typ, ok
}

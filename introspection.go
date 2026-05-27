package graphqlschema

import (
	"fmt"
	"sort"
	"strings"

	"github.com/vektah/gqlparser/v2/ast"
)

func hasGraphQLErrors(raw any) bool {
	if raw == nil {
		return false
	}
	if list, ok := raw.([]any); ok {
		return len(list) > 0
	}
	return true
}

func introspectionSchemaMap(raw map[string]any) (map[string]any, error) {
	if schema, ok, err := optionalMapValue(raw["__schema"], "graphql introspection.__schema"); err != nil || ok {
		return schema, err
	}
	data, ok, err := optionalMapValue(raw["data"], "graphql introspection.data")
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("graphql introspection response missing __schema")
	}
	schema, ok, err := optionalMapValue(data["__schema"], "graphql introspection.data.__schema")
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("graphql introspection response missing __schema")
	}
	return schema, nil
}

func parseIntrospectionTypes(model *Model, schemaMap map[string]any) error {
	types, ok, err := optionalListValue(schemaMap["types"], "graphql introspection.__schema.types")
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("graphql introspection.__schema.types must be a list")
	}
	for i, raw := range types {
		context := fmt.Sprintf("graphql introspection.__schema.types[%d]", i)
		typeMap, ok, err := optionalMapValue(raw, context)
		if err != nil {
			return err
		}
		if !ok {
			return fmt.Errorf("%s must be an object", context)
		}
		def, err := typeDefinitionFromIntrospection(typeMap, context)
		if err != nil {
			return err
		}
		if def.Name != "" {
			model.Types[def.Name] = def
		}
	}
	return nil
}

func typeDefinitionFromIntrospection(typeMap map[string]any, context string) (*TypeDefinition, error) {
	def := &TypeDefinition{
		Kind:        stringValue(typeMap["kind"]),
		Name:        stringValue(typeMap["name"]),
		Description: stringValue(typeMap["description"]),
	}
	if def.Kind == "" {
		return nil, fmt.Errorf("%s.kind must be a string", context)
	}
	if interfaces, err := namedTypeList(typeMap["interfaces"], context+".interfaces"); err != nil {
		return nil, err
	} else {
		def.Interfaces = interfaces
	}
	if possible, err := namedTypeList(typeMap["possibleTypes"], context+".possibleTypes"); err != nil {
		return nil, err
	} else {
		def.PossibleTypes = possible
	}
	if fields, err := fieldsFromIntrospection(typeMap["fields"], context+".fields"); err != nil {
		return nil, err
	} else {
		def.Fields = fields
	}
	if inputFields, err := inputValuesFromIntrospection(typeMap["inputFields"], context+".inputFields"); err != nil {
		return nil, err
	} else {
		def.InputFields = inputFields
	}
	if enumValues, err := enumValuesFromIntrospection(typeMap["enumValues"], context+".enumValues"); err != nil {
		return nil, err
	} else {
		def.EnumValues = enumValues
	}
	return def, nil
}

func parseIntrospectionDirectives(model *Model, schemaMap map[string]any) error {
	directives, ok, err := optionalListValue(schemaMap["directives"], "graphql introspection.__schema.directives")
	if err != nil {
		return err
	}
	if !ok {
		return nil
	}
	for i, raw := range directives {
		context := fmt.Sprintf("graphql introspection.__schema.directives[%d]", i)
		dirMap, ok, err := optionalMapValue(raw, context)
		if err != nil {
			return err
		}
		if !ok {
			return fmt.Errorf("%s must be an object", context)
		}
		def, err := directiveFromIntrospection(dirMap, context)
		if err != nil {
			return err
		}
		if def.Name != "" {
			model.Directives[def.Name] = def
		}
	}
	return nil
}

func directiveFromIntrospection(dirMap map[string]any, context string) (*DirectiveDefinition, error) {
	def := &DirectiveDefinition{
		Name:        stringValue(dirMap["name"]),
		Description: stringValue(dirMap["description"]),
		Repeatable:  boolValue(dirMap["isRepeatable"]),
	}
	if locations, ok, err := optionalListValue(dirMap["locations"], context+".locations"); err != nil {
		return nil, err
	} else if ok {
		for i, raw := range locations {
			location, ok := raw.(string)
			if !ok {
				return nil, fmt.Errorf("%s.locations[%d] must be a string", context, i)
			}
			if strings.TrimSpace(location) != "" {
				def.Locations = append(def.Locations, location)
			}
		}
		sort.Strings(def.Locations)
	}
	if args, err := inputValuesFromIntrospection(dirMap["args"], context+".args"); err != nil {
		return nil, err
	} else {
		def.Arguments = args
	}
	return def, nil
}

func fieldsFromIntrospection(raw any, context string) ([]*FieldDefinition, error) {
	list, ok, err := optionalListValue(raw, context)
	if err != nil || !ok {
		return nil, err
	}
	var out []*FieldDefinition
	for i, rawField := range list {
		fieldContext := fmt.Sprintf("%s[%d]", context, i)
		fieldMap, ok, err := optionalMapValue(rawField, fieldContext)
		if err != nil {
			return nil, err
		}
		if !ok {
			return nil, fmt.Errorf("%s must be an object", fieldContext)
		}
		field := &FieldDefinition{
			Name:        stringValue(fieldMap["name"]),
			Description: stringValue(fieldMap["description"]),
		}
		field.Type, err = requiredTypeRef(fieldMap["type"], fieldContext+".type")
		if err != nil {
			return nil, err
		}
		field.Arguments, err = inputValuesFromIntrospection(fieldMap["args"], fieldContext+".args")
		if err != nil {
			return nil, err
		}
		out = append(out, field)
	}
	return out, nil
}

func inputValuesFromIntrospection(raw any, context string) ([]*InputValueDefinition, error) {
	list, ok, err := optionalListValue(raw, context)
	if err != nil || !ok {
		return nil, err
	}
	var out []*InputValueDefinition
	for i, rawInput := range list {
		inputContext := fmt.Sprintf("%s[%d]", context, i)
		inputMap, ok, err := optionalMapValue(rawInput, inputContext)
		if err != nil {
			return nil, err
		}
		if !ok {
			return nil, fmt.Errorf("%s must be an object", inputContext)
		}
		input := &InputValueDefinition{
			Name:         stringValue(inputMap["name"]),
			Description:  stringValue(inputMap["description"]),
			DefaultValue: stringValue(inputMap["defaultValue"]),
		}
		input.Type, err = requiredTypeRef(inputMap["type"], inputContext+".type")
		if err != nil {
			return nil, err
		}
		out = append(out, input)
	}
	return out, nil
}

func enumValuesFromIntrospection(raw any, context string) ([]*EnumValueDefinition, error) {
	list, ok, err := optionalListValue(raw, context)
	if err != nil || !ok {
		return nil, err
	}
	var out []*EnumValueDefinition
	for i, rawEnum := range list {
		enumContext := fmt.Sprintf("%s[%d]", context, i)
		enumMap, ok, err := optionalMapValue(rawEnum, enumContext)
		if err != nil {
			return nil, err
		}
		if !ok {
			return nil, fmt.Errorf("%s must be an object", enumContext)
		}
		out = append(out, &EnumValueDefinition{
			Name:        stringValue(enumMap["name"]),
			Description: stringValue(enumMap["description"]),
		})
	}
	return out, nil
}

func requiredTypeRef(raw any, context string) (*TypeRef, error) {
	ref, err := typeRefFromIntrospection(raw, context)
	if err != nil {
		return nil, err
	}
	if ref == nil {
		return nil, fmt.Errorf("%s must be an object", context)
	}
	return ref, nil
}

func typeRefFromIntrospection(raw any, context string) (*TypeRef, error) {
	refMap, ok, err := optionalMapValue(raw, context)
	if err != nil || !ok {
		return nil, err
	}
	kind := stringValue(refMap["kind"])
	name := stringValue(refMap["name"])
	ofType, err := typeRefFromIntrospection(refMap["ofType"], context+".ofType")
	if err != nil {
		return nil, err
	}
	switch kind {
	case "NON_NULL":
		if ofType == nil {
			return nil, fmt.Errorf("%s.ofType must be an object for NON_NULL", context)
		}
		ofType.NonNull = true
		return ofType, nil
	case "LIST":
		if ofType == nil {
			return nil, fmt.Errorf("%s.ofType must be an object for LIST", context)
		}
		return &TypeRef{Elem: ofType}, nil
	default:
		if name == "" {
			return nil, fmt.Errorf("%s.name must be a string", context)
		}
		return &TypeRef{NamedType: name}, nil
	}
}

func namedTypeField(parent map[string]any, field string, context string) (string, error) {
	raw, exists := parent[field]
	if !exists || raw == nil {
		return "", nil
	}
	value, ok, err := optionalMapValue(raw, context+"."+field)
	if err != nil {
		return "", err
	}
	if !ok {
		return "", fmt.Errorf("%s.%s must be an object", context, field)
	}
	return stringValue(value["name"]), nil
}

func namedTypeList(raw any, context string) ([]string, error) {
	list, ok, err := optionalListValue(raw, context)
	if err != nil || !ok {
		return nil, err
	}
	var out []string
	for i, item := range list {
		itemMap, ok, err := optionalMapValue(item, fmt.Sprintf("%s[%d]", context, i))
		if err != nil {
			return nil, err
		}
		if !ok {
			return nil, fmt.Errorf("%s[%d] must be an object", context, i)
		}
		if name := stringValue(itemMap["name"]); name != "" {
			out = append(out, name)
		}
	}
	return out, nil
}

func optionalMapValue(raw any, context string) (map[string]any, bool, error) {
	if raw == nil {
		return nil, false, nil
	}
	value, ok := raw.(map[string]any)
	if !ok {
		return nil, false, fmt.Errorf("%s must be an object", context)
	}
	return value, true, nil
}

func optionalListValue(raw any, context string) ([]any, bool, error) {
	if raw == nil {
		return nil, false, nil
	}
	value, ok := raw.([]any)
	if !ok {
		return nil, false, fmt.Errorf("%s must be a list", context)
	}
	return value, true, nil
}

func stringValue(raw any) string {
	value, _ := raw.(string)
	return strings.TrimSpace(value)
}

func boolValue(raw any) bool {
	value, _ := raw.(bool)
	return value
}

func sortedASTTypeNames(values map[string]*ast.Definition) []string {
	names := make([]string, 0, len(values))
	for name := range values {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func sortedASTDirectiveNames(values map[string]*ast.DirectiveDefinition) []string {
	names := make([]string, 0, len(values))
	for name := range values {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

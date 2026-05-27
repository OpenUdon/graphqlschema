package graphqlschema

import (
	"net/url"
	"sort"
	"strings"
)

// SelectorAliases returns all operation ids and local JSON Pointer selectors
// that ResolveSelector can resolve.
func (m *Model) SelectorAliases() map[string]bool {
	aliases := map[string]bool{}
	if m == nil {
		return aliases
	}
	for _, op := range m.Operations {
		if op == nil || op.ID == "" || op.RootType == "" || op.FieldName == "" {
			continue
		}
		aliases[op.ID] = true
		aliases[jsonPointer("operations", op.ID)] = true
		aliases[jsonPointer("types", op.RootType, "fields", op.FieldName)] = true
	}
	return aliases
}

// ResolveSelector resolves a canonical operation id or supported local JSON
// Pointer selector. Supported pointers are #/operations/{id} and
// #/types/{rootType}/fields/{fieldName}.
func (m *Model) ResolveSelector(selector string) (*SelectorTarget, bool) {
	if m == nil {
		return nil, false
	}
	selector = strings.TrimSpace(selector)
	if selector == "" {
		return nil, false
	}
	if !strings.HasPrefix(selector, "#/") {
		op, ok := m.OperationByID(selector)
		if !ok {
			return nil, false
		}
		return m.selectorTarget(op, selector)
	}
	parts, ok := decodeJSONPointer(selector)
	if !ok {
		return nil, false
	}
	switch {
	case len(parts) == 2 && parts[0] == "operations":
		op, ok := m.OperationByID(parts[1])
		if !ok {
			return nil, false
		}
		return m.selectorTarget(op, selector)
	case len(parts) == 4 && parts[0] == "types" && parts[2] == "fields":
		op, ok := m.operationByRootField(parts[1], parts[3])
		if !ok {
			return nil, false
		}
		return m.selectorTarget(op, selector)
	default:
		return nil, false
	}
}

func (m *Model) selectorTarget(op *Operation, selector string) (*SelectorTarget, bool) {
	if op == nil {
		return nil, false
	}
	root, field, ok := m.rootField(op.RootType, op.FieldName)
	if !ok {
		return nil, false
	}
	return &SelectorTarget{
		Kind:      "operation",
		Selector:  selector,
		Operation: op,
		Field:     field,
		RootType:  root,
	}, true
}

func (m *Model) operationByRootField(rootType string, fieldName string) (*Operation, bool) {
	for _, op := range m.Operations {
		if op != nil && op.RootType == rootType && op.FieldName == fieldName {
			return op, true
		}
	}
	return nil, false
}

func (m *Model) rootField(rootType string, fieldName string) (*TypeDefinition, *FieldDefinition, bool) {
	root, ok := m.TypeByName(rootType)
	if !ok || root == nil {
		return nil, nil, false
	}
	for _, field := range root.Fields {
		if field != nil && field.Name == fieldName {
			return root, field, true
		}
	}
	return nil, nil, false
}

func (m *Model) buildOperations() {
	if m == nil {
		return
	}
	var operations []*Operation
	for _, root := range []struct {
		kind string
		name string
	}{
		{kind: "query", name: m.QueryType},
		{kind: "mutation", name: m.MutationType},
		{kind: "subscription", name: m.SubscriptionType},
	} {
		rootType, ok := m.TypeByName(root.name)
		if !ok || rootType == nil {
			continue
		}
		for _, field := range rootType.Fields {
			if field == nil || field.Name == "" || strings.HasPrefix(field.Name, "__") {
				continue
			}
			operations = append(operations, &Operation{
				ID:          root.kind + "." + field.Name,
				Kind:        root.kind,
				RootType:    root.name,
				FieldName:   field.Name,
				Description: field.Description,
				Type:        cloneTypeRef(field.Type),
				Arguments:   cloneInputValues(field.Arguments),
				Directives:  cloneDirectiveUses(field.Directives),
			})
		}
	}
	sort.SliceStable(operations, func(i, j int) bool {
		if operations[i].Kind != operations[j].Kind {
			return operationKindOrder(operations[i].Kind) < operationKindOrder(operations[j].Kind)
		}
		if operations[i].RootType != operations[j].RootType {
			return operations[i].RootType < operations[j].RootType
		}
		return operations[i].FieldName < operations[j].FieldName
	})
	m.Operations = operations
}

func operationKindOrder(kind string) int {
	switch kind {
	case "query":
		return 0
	case "mutation":
		return 1
	case "subscription":
		return 2
	default:
		return 3
	}
}

func cloneTypeRef(ref *TypeRef) *TypeRef {
	if ref == nil {
		return nil
	}
	return &TypeRef{
		NamedType: ref.NamedType,
		Elem:      cloneTypeRef(ref.Elem),
		NonNull:   ref.NonNull,
	}
}

func cloneInputValues(values []*InputValueDefinition) []*InputValueDefinition {
	if len(values) == 0 {
		return nil
	}
	out := make([]*InputValueDefinition, 0, len(values))
	for _, value := range values {
		if value == nil {
			continue
		}
		out = append(out, &InputValueDefinition{
			Name:         value.Name,
			Description:  value.Description,
			Type:         cloneTypeRef(value.Type),
			DefaultValue: value.DefaultValue,
			Directives:   cloneDirectiveUses(value.Directives),
		})
	}
	return out
}

func cloneDirectiveUses(values []*DirectiveUse) []*DirectiveUse {
	if len(values) == 0 {
		return nil
	}
	out := make([]*DirectiveUse, 0, len(values))
	for _, value := range values {
		if value == nil {
			continue
		}
		clone := &DirectiveUse{Name: value.Name}
		if len(value.Arguments) > 0 {
			clone.Arguments = map[string]string{}
			for k, v := range value.Arguments {
				clone.Arguments[k] = v
			}
		}
		out = append(out, clone)
	}
	return out
}

func jsonPointer(parts ...string) string {
	escaped := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.ReplaceAll(part, "~", "~0")
		part = strings.ReplaceAll(part, "/", "~1")
		escaped = append(escaped, part)
	}
	return "#/" + strings.Join(escaped, "/")
}

func decodeJSONPointer(selector string) ([]string, bool) {
	if !strings.HasPrefix(selector, "#/") {
		return nil, false
	}
	rawParts := strings.Split(strings.TrimPrefix(selector, "#/"), "/")
	parts := make([]string, 0, len(rawParts))
	for _, raw := range rawParts {
		unescaped, err := url.PathUnescape(raw)
		if err != nil {
			return nil, false
		}
		unescaped = strings.ReplaceAll(unescaped, "~1", "/")
		unescaped = strings.ReplaceAll(unescaped, "~0", "~")
		parts = append(parts, unescaped)
	}
	return parts, true
}

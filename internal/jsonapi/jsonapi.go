package jsonapi

import (
	"encoding/json"
	"fmt"
)

type Document struct {
	Data     json.RawMessage `json:"data,omitempty"`
	Errors   []Error         `json:"errors,omitempty"`
	Included json.RawMessage `json:"included,omitempty"`
	Links    json.RawMessage `json:"links,omitempty"`
	Meta     json.RawMessage `json:"meta,omitempty"`
}

type Resource struct {
	ID            string                     `json:"id"`
	Type          string                     `json:"type"`
	Attributes    map[string]json.RawMessage `json:"attributes,omitempty"`
	Relationships map[string]json.RawMessage `json:"relationships,omitempty"`
	Links         json.RawMessage            `json:"links,omitempty"`
	Meta          json.RawMessage            `json:"meta,omitempty"`
}

type Relationship struct {
	Data  json.RawMessage `json:"data,omitempty"`
	Links json.RawMessage `json:"links,omitempty"`
	Meta  json.RawMessage `json:"meta,omitempty"`
}

type ResourceIdentifier struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

type Error struct {
	ID     string          `json:"id,omitempty"`
	Status string          `json:"status,omitempty"`
	Code   string          `json:"code,omitempty"`
	Title  string          `json:"title,omitempty"`
	Detail string          `json:"detail,omitempty"`
	Source json.RawMessage `json:"source,omitempty"`
	Meta   json.RawMessage `json:"meta,omitempty"`
}

func DecodeResource(data json.RawMessage) (Resource, error) {
	var resource Resource
	err := json.Unmarshal(data, &resource)
	return resource, err
}

func DecodeResourceCollection(data json.RawMessage) ([]Resource, error) {
	var resources []Resource
	err := json.Unmarshal(data, &resources)
	return resources, err
}

func DecodeIncluded(data json.RawMessage) ([]Resource, error) {
	if len(data) == 0 {
		return nil, nil
	}
	return DecodeResourceCollection(data)
}

func DecodeRelationship(data json.RawMessage) (Relationship, error) {
	var relationship Relationship
	err := json.Unmarshal(data, &relationship)
	return relationship, err
}

func DecodeResourceIdentifiers(data json.RawMessage) ([]ResourceIdentifier, error) {
	if len(data) == 0 || string(data) == "null" {
		return nil, nil
	}

	var identifiers []ResourceIdentifier
	if err := json.Unmarshal(data, &identifiers); err == nil {
		return identifiers, nil
	}

	var identifier ResourceIdentifier
	if err := json.Unmarshal(data, &identifier); err != nil {
		return nil, err
	}
	return []ResourceIdentifier{identifier}, nil
}

func RelationshipIdentifier(resource Resource, relationshipName string) (ResourceIdentifier, error) {
	relationshipData, ok := resource.Relationships[relationshipName]
	if !ok {
		return ResourceIdentifier{}, fmt.Errorf("relationship %q is missing from %s %s", relationshipName, resource.Type, resource.ID)
	}

	relationship, err := DecodeRelationship(relationshipData)
	if err != nil {
		return ResourceIdentifier{}, fmt.Errorf("decode relationship %q: %w", relationshipName, err)
	}

	identifiers, err := DecodeResourceIdentifiers(relationship.Data)
	if err != nil {
		return ResourceIdentifier{}, fmt.Errorf("decode relationship %q identifier: %w", relationshipName, err)
	}
	if len(identifiers) != 1 {
		return ResourceIdentifier{}, fmt.Errorf("relationship %q must contain one resource", relationshipName)
	}

	return identifiers[0], nil
}

func FindIncluded(document Document, identifier ResourceIdentifier) (Resource, error) {
	includedResources, err := DecodeIncluded(document.Included)
	if err != nil {
		return Resource{}, fmt.Errorf("decode included resources: %w", err)
	}

	for _, resource := range includedResources {
		if resource.ID == identifier.ID && resource.Type == identifier.Type {
			return resource, nil
		}
	}

	return Resource{}, fmt.Errorf("included %s %s was not returned", identifier.Type, identifier.ID)
}

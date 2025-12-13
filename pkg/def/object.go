package def

import (
	_ "embed"
	"fmt"
)

type Object struct {
	Mib                string           `yaml:"mib,omitempty" json:"mib,omitempty"`
	ObjectName         string           `yaml:"object,omitempty" json:"object,omitempty"`
	Type               string           `yaml:"type,omitempty" json:"type,omitempty"`
	Index              []ObjectIndex    `yaml:"index,omitempty" json:"index,omitempty"`
	Augments           string           `yaml:"augments,omitempty" json:"augments,omitempty"`
	DiscoveryAttribute string           `yaml:"discovery_attribute,omitempty" json:"discovery_attribute,omitempty"`
	Attributes         ObjectAttributes `yaml:"attributes,omitempty" json:"attributes,omitempty"`
}

func (o Object) Validate() error {
	if err := validate(o, "schemas/object.json"); err != nil {
		return err
	}
	// Validate discovery_attribute exists in attributes
	if o.DiscoveryAttribute != "" {
		if _, ok := o.Attributes[o.DiscoveryAttribute]; !ok {
			return fmt.Errorf("discovery_attribute %q not found in attributes", o.DiscoveryAttribute)
		}
	}
	return nil
}

func (o Object) Kind() string {
	return "object"
}

func (o Object) GetDiscoveryOid() (string, error) {
	for key, attribute := range o.Attributes {
		if key == o.DiscoveryAttribute {
			return attribute.Oid, nil
		}
	}
	return "", fmt.Errorf("discovery attribute not present for object %q", o.ObjectName)
}

func (o Object) GetOverrideOids() []string {
	overrideOids := make([]string, 0)
	for _, attribute := range o.Attributes {
		if attribute.Overrides.Attribute != "" && attribute.Overrides.Object != "" {
			overrideOids = append(overrideOids, attribute.Oid)
		}
	}
	return overrideOids
}

type ObjectIndex struct {
	Type   string `yaml:"type,omitempty" json:"type,omitempty"`
	Oid    string `yaml:"oid,omitempty" json:"oid,omitempty"`
	Name   string `yaml:"name,omitempty" json:"name,omitempty"`
	Syntax string `yaml:"syntax,omitempty" json:"syntax,omitempty"`
}

type ObjectAttributes map[string]ObjectAttribute

type ObjectAttribute struct {
	Oid        string                  `yaml:"oid,omitempty" json:"oid,omitempty"`
	Tag        bool                    `yaml:"tag,omitempty" json:"tag,omitempty"`
	Name       string                  `yaml:"name,omitempty" json:"name,omitempty"`
	Syntax     string                  `yaml:"syntax,omitempty" json:"syntax,omitempty"`
	Semantic   string                  `yaml:"semantic,omitempty" json:"semantics,omitempty"`
	Metric     string                  `yaml:"metric,omitempty" json:"metric,omitempty"`
	Rediscover string                  `yaml:"rediscover,omitempty" json:"rediscover,omitempty"`
	Overrides  ObjectAttributeOverride `yaml:"overrides,omitempty" json:"overrides,omitempty"`
}

type ObjectAttributeOverride struct {
	Object    string `yaml:"object,omitempty" json:"object,omitempty"`
	Attribute string `yaml:"attribute,omitempty" json:"attribute,omitempty"`
}

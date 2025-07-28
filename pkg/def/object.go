package def

import _ "embed"

type Object struct {
	Mib                string           `yaml:"mib,omitempty" json:"mib,omitempty"`
	ObjectName         string           `yaml:"object,omitempty" json:"object,omitempty"`
	Index              []ObjectIndex    `yaml:"index,omitempty" json:"index,omitempty"`
	Augments           string           `yaml:"augments,omitempty" json:"augments,omitempty"`
	DiscoveryAttribute string           `yaml:"discovery_attribute,omitempty" json:"discovery_attribute,omitempty"`
	Attributes         ObjectAttributes `yaml:"attributes,omitempty" json:"attributes,omitempty"`
}

func (o Object) Validate() error {
	return validate(o, "schemas/object.json")
}

func (o Object) Type() string {
	return "object"
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
	Semantics  string                  `yaml:"semantics,omitempty" json:"semantics,omitempty"`
	Metric     string                  `yaml:"metric,omitempty" json:"metric,omitempty"`
	Rediscover string                  `yaml:"rediscover,omitempty" json:"rediscover,omitempty"`
	Overrides  ObjectAttributeOverride `yaml:"overrides,omitempty" json:"overrides,omitempty"`
}

type ObjectAttributeOverride struct {
	Object    string `yaml:"object,omitempty" json:"object,omitempty"`
	Attribute string `yaml:"attribute,omitempty" json:"attribute,omitempty"`
}

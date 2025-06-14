package def

type Device struct {
	IP                 string         `yaml:"ip,omitempty" json:"ip,omitempty"`
	Port               uint16         `yaml:"port,omitempty" json:"port,omitempty"`
	Timeout            uint64         `yaml:"timeout,omitempty" json:"timeout,omitempty"`
	Retries            int            `yaml:"retries,omitempty" json:"retries,omitempty"`
	ExponentialTimeout bool           `yaml:"exponential_timeout,omitempty" json:"exponential_timeout,omitempty"`
	Version            string         `yaml:"version,omitempty" json:"version,omitempty"`
	Communities        []string       `yaml:"communities,omitempty" json:"communities,omitempty"`
	DeviceGroups       []string       `yaml:"device_groups,omitempty" json:"device_groups,omitempty"`
	PollInterval       uint64         `yaml:"poll_interval,omitempty" json:"poll_interval,omitempty"`
	MaxOIDs            uint64         `yaml:"max_oids,omitempty" json:"max_oids,omitempty"`
	V3Credentials      []V3Credential `yaml:"v3_credentials,omitempty" json:"v3_credentials,omitempty"`
	MaxRepetitions     uint32         `yaml:"max_repetitions,omitempty" json:"max_repetitions,omitempty"`
	MaxConcurrentPolls uint32         `yaml:"max_concurrent_polls,omitempty" json:"max_concurrent_polls,omitempty"`
	CiscoQosEnabled    bool           `yaml:"cisco_qos_enabled,omitempty" json:"cisco_qos_enabled,omitempty"`
}

func (d Device) Validate() error {
	return nil
}

func (d Device) Type() string {
	return "device"
}

type V3Credential struct {
	AuthoritativeEngineID    string `yaml:"authoritative_engine_id,omitempty" json:"authoritative_engine_id,omitempty"`
	AuthoritativeEngineBoots uint32 `yaml:"authoritative_engine_boots,omitempty" json:"authoritative_engine_boots,omitempty"`
	AuthoritativeEngineTime  uint32 `yaml:"authoritative_engine_time,omitempty" json:"authoritative_engine_time,omitempty"`

	Username                 string `yaml:"username,omitempty" json:"username,omitempty"`
	AuthenticationParameters string `yaml:"authentication_parameters,omitempty" json:"authentication_parameters,omitempty"`
	PrivacyParameters        string `yaml:"privacy_parameters,omitempty" json:"privacy_parameters,omitempty"`

	// string value will be mapped to gosnmp.SnmpV3AuthProtocol
	AuthenticationProtocol string `yaml:"authentication_protocol,omitempty" json:"authentication_protocol,omitempty"`
	// string value will be mapped to gosnmp.SnmpV3PrivProtocol
	PrivacyProtocol string `yaml:"privacy_protocol,omitempty" json:"privacy_protocol,omitempty"`

	AuthenticationPassphrase string `yaml:"authentication_passphrase,omitempty" json:"authentication_passphrase,omitempty"`
	PrivacyPassphrase        string `yaml:"privacy_passphrase,omitempty" json:"privacy_passphrase,omitempty"`

	SecretKey  string `yaml:"secret_key,omitempty" json:"secret_key,omitempty"`
	PrivacyKey string `yaml:"privacy_key,omitempty" json:"privacy_key,omitempty"`
}

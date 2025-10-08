package def

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

func ValidateDevices(
	devices map[string]Device,
	deviceGroups map[string]DeviceGroup,
	objectTypes map[string]ObjectType,
) error {
	invalidDefinitions := make([]string, 0)
	deviceIPs := make(map[string]string)

	for deviceName, device := range devices {
		// Ensure the device has valid device groups
		for _, dg := range device.DeviceGroups {
			if _, ok := deviceGroups[dg]; !ok {
				invalidDefinitions = append(invalidDefinitions,
					fmt.Sprintf("device %q references an undefined device group: %q", deviceName, dg),
				)
			}
		}

		// Ensure the device has a unique IP address.
		if existingDeviceName, ok := deviceIPs[device.IP]; ok {
			invalidDefinitions = append(invalidDefinitions,
				fmt.Sprintf("device %q has the same IP address (%s) as device %q", deviceName, device.IP, existingDeviceName),
			)
		} else {
			deviceIPs[device.IP] = deviceName
		}

		// Ensure the device has valid object types
		for objectType := range device.PollIntervals {
			if _, ok := objectTypes[objectType]; !ok {
				invalidDefinitions = append(invalidDefinitions,
					fmt.Sprintf("device %q references an undefined object type: %q", deviceName, objectType),
				)
			}
		}
	}

	if len(invalidDefinitions) > 0 {
		return fmt.Errorf("found %d device definition errors:\n%s",
			len(invalidDefinitions),
			strings.Join(invalidDefinitions, "\n"),
		)
	}

	return nil
}

type Device struct {
	IP                 string            `yaml:"ip,omitempty" json:"ip,omitempty"`
	Port               uint16            `yaml:"port,omitempty" json:"port,omitempty"`
	Timeout            uint64            `yaml:"timeout,omitempty" json:"timeout,omitempty"`
	Retries            int               `yaml:"retries,omitempty" json:"retries,omitempty"`
	ExponentialTimeout bool              `yaml:"exponential_timeout,omitempty" json:"exponential_timeout,omitempty"`
	Version            string            `yaml:"version,omitempty" json:"version,omitempty"`
	Communities        []string          `yaml:"communities,omitempty" json:"communities,omitempty"`
	DeviceGroups       []string          `yaml:"device_groups,omitempty" json:"device_groups,omitempty"`
	PollInterval       uint64            `yaml:"poll_interval,omitempty" json:"poll_interval,omitempty"`
	PollIntervals      map[string]uint64 `yaml:"poll_intervals,omitempty" json:"poll_intervals,omitempty"`
	MaxOIDs            uint64            `yaml:"max_oids,omitempty" json:"max_oids,omitempty"`
	V3Credentials      []V3Credential    `yaml:"v3_credentials,omitempty" json:"v3_credentials,omitempty"`
	MaxRepetitions     uint32            `yaml:"max_repetitions,omitempty" json:"max_repetitions,omitempty"`
	MaxConcurrentPolls uint32            `yaml:"max_concurrent_polls,omitempty" json:"max_concurrent_polls,omitempty"`
	CiscoQosEnabled    bool              `yaml:"cisco_qos_enabled,omitempty" json:"cisco_qos_enabled,omitempty"`
}

func (d Device) Validate() error {
	return nil
}

func (d Device) Kind() string {
	return "device"
}

// Id returns the device's unique id
func (d Device) Id() string {
	return d.IP + ":" + strconv.FormatUint(uint64(d.Port), 10)
}

// Unmarshal data into the struct
func (d Device) Unmarshal(data []byte, dataType string) error {
	if strings.ToLower(dataType) == "yml" || strings.ToLower(dataType) == "yaml" {
		return yaml.Unmarshal(data, d)
	}
	return json.Unmarshal(data, &d)
}

// Marshal marshals the struct to YAML or JSON
func (d Device) Marshal(dataType string) ([]byte, error) {
	if strings.ToLower(dataType) == "yml" || strings.ToLower(dataType) == "yaml" {
		return yaml.Marshal(d)
	}
	return json.Marshal(d)
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

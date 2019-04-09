package model

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	apiv1 "k8s.io/api/core/v1"
	resource "k8s.io/apimachinery/pkg/api/resource"
)

type serialMount struct {
	SendOnly bool   `json:"sendonly,omitempty" yaml:"sendonly,omitempty"`
	Source   string `json:"source,omitempty" yaml:"source,omitempty"`
	Path     string `json:"path" yaml:"path,omitempty"`
	Target   string `json:"target,omitempty" yaml:"target,omitempty"` //TODO: decrecated
	Size     string `json:"size,omitempty" yaml:"size,omitempty"`
}

// UnmarshalYAML Implements the Unmarshaler interface of the yaml pkg.
func (e *EnvVar) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var raw string
	err := unmarshal(&raw)
	if err != nil {
		return err
	}

	parts := strings.SplitN(raw, "=", 2)
	e.Name = parts[0]
	if len(parts) == 2 {
		if strings.HasPrefix(parts[1], "$") {
			e.Value = os.ExpandEnv(parts[1])
			return nil
		}

		e.Value = parts[1]
		return nil
	}

	val := os.ExpandEnv(parts[0])
	if val != parts[0] {
		e.Value = val
	}

	return nil
}

// MarshalYAML Implements the marshaler interface of the yaml pkg.
func (e *EnvVar) MarshalYAML() (interface{}, error) {
	return e.Name + "=" + e.Value, nil
}

// UnmarshalYAML Implements the Unmarshaler interface of the yaml pkg.
func (m *Mount) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var sMount serialMount
	var raw string
	err := unmarshal(&raw)
	if err == nil {
		m.Path = raw
	} else {
		err := unmarshal(&sMount)
		if err != nil {
			return err
		}
		m.SendOnly = sMount.SendOnly
		m.Source = sMount.Source
		m.Path = sMount.Path
		m.Target = sMount.Target
		m.Size = sMount.Size
	}
	return nil
}

// MarshalYAML Implements the marshaler interface of the yaml pkg.
func (m *Mount) MarshalYAML() (interface{}, error) {
	if !m.SendOnly && m.Source == "" && m.Target == "" && m.Size == "" {
		return m.Path, nil
	}
	return &serialMount{
		SendOnly: m.SendOnly,
		Source:   m.Source,
		Path:     m.Path,
		Target:   m.Target,
		Size:     m.Size,
	}, nil
}

// UnmarshalYAML Implements the Unmarshaler interface of the yaml pkg.
func (f *Forward) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var raw string
	err := unmarshal(&raw)
	if err != nil {
		return err
	}

	parts := strings.SplitN(raw, ":", 2)
	if len(parts) != 2 {
		return fmt.Errorf("Wrong port-forward syntax '%s', must be of the form 'localPort:RemotePort'", raw)
	}
	localPort, err := strconv.Atoi(parts[0])
	if err != nil {
		return fmt.Errorf("Cannot convert remote port '%s' in port-forward '%s'", parts[0], raw)
	}
	remotePort, err := strconv.Atoi(parts[1])
	if err != nil {
		return fmt.Errorf("Cannot convert remote port '%s' in port-forward '%s'", parts[1], raw)
	}
	f.Local = localPort
	f.Remote = remotePort
	return nil
}

// MarshalYAML Implements the marshaler interface of the yaml pkg.
func (f Forward) MarshalYAML() (interface{}, error) {
	return fmt.Sprintf("%d:%d", f.Local, f.Remote), nil
}

// UnmarshalYAML Implements the Unmarshaler interface of the yaml pkg.
func (r *ResourceList) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var raw map[apiv1.ResourceName]string
	err := unmarshal(&raw)
	if err != nil {
		return err
	}

	for k, v := range raw {
		parsed, err := resource.ParseQuantity(v)
		if err != nil {
			return err
		}

		(*r)[k] = parsed
	}

	return nil
}

// MarshalYAML Implements the marshaler interface of the yaml pkg.
func (r ResourceList) MarshalYAML() (interface{}, error) {
	m := make(map[apiv1.ResourceName]string, 0)
	for k, v := range r {
		m[k] = v.String()
	}

	return m, nil
}
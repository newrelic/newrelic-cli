package configuration

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

type ConfigProvider struct {
	cfg            []byte
	fileName       string
	scope          string
	dirty          bool
	mu             sync.Mutex
	values         []FieldDefinition
	explicitValues bool
}

type FieldDefinition struct {
	EnvVar         string
	Key            string
	Default        interface{}
	CaseSensitive  bool
	ValidationFunc ConfigValueValidationFunc
}

type ConfigValueValidationFunc func(key string, value interface{}) error

func IntGreaterThan(greaterThan int) func(key string, value interface{}) error {
	return func(key string, value interface{}) error {
		var s int
		var ok bool
		if s, ok = value.(int); !ok {
			return fmt.Errorf("%v is not an int", value)
		}

		if s > greaterThan {
			return nil
		}

		return fmt.Errorf("value %d is not greater than %d", s, greaterThan)
	}
}

func StringInStrings(caseSensitive bool, allowedValues ...string) func(key string, value interface{}) error {
	return func(key string, value interface{}) error {
		var s string
		var ok bool
		if s, ok = value.(string); !ok {
			return fmt.Errorf("%v is not a string", value)
		}

		for _, v := range allowedValues {
			if caseSensitive && s == v {
				return nil
			}

			if !caseSensitive && strings.EqualFold(s, v) {
				return nil
			}
		}

		return fmt.Errorf("value %v not in allowed values: %s", s, allowedValues)
	}
}

type ConfigProviderOption func(*ConfigProvider) error

func NewConfigProvider(opts ...ConfigProviderOption) (*ConfigProvider, error) {
	p := &ConfigProvider{}

	for _, fn := range opts {
		if fn == nil {
			continue
		}
		if err := fn(p); err != nil {
			return nil, err
		}
	}

	return p, nil
}

func WithExplicitValues() ConfigProviderOption {
	return func(p *ConfigProvider) error {
		p.explicitValues = true
		return nil
	}
}

func WithFieldDefinitions(definitions ...FieldDefinition) ConfigProviderOption {
	return func(p *ConfigProvider) error {
		for _, d := range definitions {
			if !d.CaseSensitive {
				for _, k := range p.getConfigValueKeys() {
					if strings.EqualFold(k, d.Key) {
						return fmt.Errorf("unable to add case-insensitive field definition for %s, another field already defined with matching case-folded key", d.Key)
					}
				}
			}
			p.values = append(p.values, d)
		}
		return nil
	}
}

func WithFilePersistence(fileName string) ConfigProviderOption {
	return func(p *ConfigProvider) error {
		p.fileName = fileName
		return nil
	}
}

func WithScope(scope string) ConfigProviderOption {
	return func(p *ConfigProvider) error {
		p.scope = scope
		return nil
	}
}

func (p *ConfigProvider) GetInt(key string) (int64, error) {
	d := p.getFieldDefinition(key)

	if d != nil {
		if e, ok := os.LookupEnv(d.EnvVar); ok {
			i, err := strconv.Atoi(e)
			if err != nil {
				return 0, err
			}

			return int64(i), nil
		}
	}

	res, err := p.getFromConfig(key)
	if err != nil {
		return 0, err
	}

	return res.Int(), nil
}

func (p *ConfigProvider) GetString(key string) (string, error) {
	d := p.getFieldDefinition(key)

	if d != nil {
		if e, ok := os.LookupEnv(d.EnvVar); ok {
			return e, nil
		}
	}

	res, err := p.getFromConfig(key)
	if err != nil {
		return "", err
	}

	return res.String(), nil
}

func (p *ConfigProvider) Set(key string, value interface{}) error {
	v := p.getFieldDefinition(key)

	if v != nil {
		if v.ValidationFunc != nil {
			if err := v.ValidationFunc(key, value); err != nil {
				return err
			}
		}

		if !v.CaseSensitive {
			// use the case convention from the field definition
			key = v.Key
		}
	} else if p.explicitValues {
		return fmt.Errorf("key '%s' is not valid, valid keys are: %v", key, p.getConfigValueKeys())
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	cfg, err := sjson.Set(p.getConfig(), fmt.Sprintf("%s.%s", p.getEscapedScope(), key), value)
	if err != nil {
		return err
	}

	p.writeConfig(cfg)

	return nil
}

func (p *ConfigProvider) getFromConfig(key string) (*gjson.Result, error) {
	v := p.getFieldDefinition(key)

	if v != nil && !v.CaseSensitive {
		// use the case convention from the field definition
		key = v.Key
	}

	res := gjson.Get(p.getConfig(), fmt.Sprintf("%s.%s", p.scope, key))

	if !res.Exists() {
		return nil, fmt.Errorf("key %s is not defined", key)
	}

	return &res, nil
}

func (p *ConfigProvider) writeConfig(json string) {
	p.dirty = true
	p.cfg = []byte(json)

	if p.fileName != "" {
		os.WriteFile(p.fileName, p.cfg, 0644)
	}
}

func (p *ConfigProvider) getConfig() string {
	if p.cfg == nil || p.dirty {

		if p.fileName != "" {
			p.setConfigFromFile()
		}

		if p.cfg == nil || len(p.cfg) == 0 {
			p.setDefaultConfig()
		}

		p.dirty = false
	}

	return string(p.cfg)
}

func (p *ConfigProvider) setConfigFromFile() {
	data, err := os.ReadFile(p.fileName)
	if err != nil {
		return
	}

	p.cfg = data
}

func (p *ConfigProvider) setDefaultConfig() {
	for _, v := range p.values {
		p.setValue(v.Key, v.Default)
	}
}

func (p *ConfigProvider) setValue(key string, value interface{}) {
	cfg, err := sjson.Set(p.getConfig(), fmt.Sprintf("%s.%s", p.getEscapedScope(), key), value)
	if err != nil {
		log.Fatal(err)
	}

	p.cfg = []byte(cfg)
}

// Escape wildcard characters, as required by sjson
func (p *ConfigProvider) getEscapedScope() string {
	re := regexp.MustCompile(`([*?])`)
	return re.ReplaceAllString(p.scope, "\\$1")
}

func (p *ConfigProvider) getFieldDefinition(key string) *FieldDefinition {
	for _, v := range p.values {
		if !v.CaseSensitive && strings.EqualFold(key, v.Key) {
			return &v
		}

		if v.CaseSensitive && key == v.Key {
			return &v
		}
	}

	return nil
}

func (p *ConfigProvider) getConfigValueKeys() []string {
	var keys []string
	for _, v := range p.values {
		keys = append(keys, v.Key)
	}

	return keys
}

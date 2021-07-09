package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

type FieldKey string

type Store struct {
	cfg            []byte
	fields         []FieldDefinition
	fileName       string
	scope          string
	mu             sync.Mutex
	explicitValues bool
}

type FieldDefinition struct {
	EnvVar            string
	Key               FieldKey
	Default           interface{}
	CaseSensitive     bool
	Sensitive         bool
	SetValidationFunc FieldValueValidationFunc
}

type FieldValueValidationFunc func(key FieldKey, value interface{}) error

func IntGreaterThan(greaterThan int) func(key FieldKey, value interface{}) error {
	return func(key FieldKey, value interface{}) error {
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

func IsTernary() func(key FieldKey, value interface{}) error {
	return func(key FieldKey, value interface{}) error {
		switch v := value.(type) {
		case string:
			return Ternary(v).Valid()
		case Ternary:
			return v.Valid()
		default:
			return fmt.Errorf("value %s for key %s is not valid", value, key)
		}

	}
}

func StringInStrings(caseSensitive bool, allowedValues ...string) func(key FieldKey, value interface{}) error {
	return func(key FieldKey, value interface{}) error {
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

type StoreOption func(*Store) error

func NewStore(opts ...StoreOption) (*Store, error) {
	p := &Store{}

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

func EnforceStrictFields() StoreOption {
	return func(p *Store) error {
		p.explicitValues = true
		return nil
	}
}

func ConfigureFields(definitions ...FieldDefinition) StoreOption {
	return func(p *Store) error {
		for _, d := range definitions {
			if !d.CaseSensitive {
				for _, k := range p.getConfigValueKeys() {
					if strings.EqualFold(string(k), string(d.Key)) {
						return fmt.Errorf("unable to add case-insensitive field definition for %s, another field already defined with matching case-folded key", d.Key)
					}
				}
			}
			p.fields = append(p.fields, d)
		}
		return nil
	}
}

func PersistToFile(fileName string) StoreOption {
	return func(p *Store) error {
		p.fileName = fileName
		return nil
	}
}

func UseGlobalScope(scope string) StoreOption {
	return func(p *Store) error {
		p.scope = scope
		return nil
	}
}

func (p *Store) GetInt(key FieldKey) (int64, error) {
	return p.GetIntWithScope("", key)
}

func (p *Store) GetString(key FieldKey) (string, error) {
	return p.GetStringWithScopeAndOverride("", key, nil)
}

func (p *Store) GetStringWithOverride(key FieldKey, override *string) (string, error) {
	return p.GetStringWithScopeAndOverride("", key, override)
}

func (p *Store) GetTernary(key FieldKey) (Ternary, error) {
	return p.GetTernaryWithScope("", key)
}

func (p *Store) GetTernaryWithScope(scope string, key FieldKey) (Ternary, error) {
	v, err := p.GetWithScope(scope, key)
	if err != nil {
		return Ternary(""), err
	}

	switch v := v.(type) {
	case string:
		return Ternary(v), nil
	case Ternary:
		return v, nil
	}

	return Ternary(""), fmt.Errorf("value %v for key %s is not a ternary", v, key)
}

func (p *Store) GetIntWithScope(scope string, key FieldKey) (int64, error) {
	return p.GetIntWithScopeAndOverride(scope, key, nil)
}

func (p *Store) GetIntWithScopeAndOverride(scope string, key FieldKey, override *int64) (int64, error) {
	var v interface{}
	var err error

	if override == nil {
		v, err = p.GetWithScopeAndOverride(scope, key, nil)
	} else {
		v, err = p.GetWithScopeAndOverride(scope, key, *override)
	}

	if err != nil {
		return 0, err
	}

	switch v := v.(type) {
	case int64:
		return v, nil
	case int32:
		return int64(v), nil
	case float64:
		return int64(v), nil
	case float32:
		return int64(v), nil
	case int:
		return int64(v), nil
	case string:
		i, err := strconv.Atoi(v)
		if err != nil {
			return 0, fmt.Errorf("value %v for key %s is not an int", v, key)
		}
		return int64(i), nil
	}

	return 0, fmt.Errorf("value %v for key %s is not an int", v, key)
}

func (p *Store) GetStringWithScope(scope string, key FieldKey) (string, error) {
	return p.GetStringWithScopeAndOverride(scope, key, nil)
}

func (p *Store) GetStringWithScopeAndOverride(scope string, key FieldKey, override *string) (string, error) {
	var v interface{}
	var err error

	if override == nil || *override == "" {
		v, err = p.GetWithScope(scope, key)
	} else {
		v, err = p.GetWithScopeAndOverride(scope, key, *override)
	}

	if err != nil {
		return "", err
	}

	switch v := v.(type) {
	case int64:
		return strconv.Itoa(int(v)), nil
	case int32:
		return strconv.Itoa(int(v)), nil
	case float64:
		return strconv.Itoa(int(v)), nil
	case float32:
		return strconv.Itoa(int(v)), nil
	case int:
		return strconv.Itoa(v), nil
	case string:
		return v, nil
	}

	return "", fmt.Errorf("value %v for key %s is not a string", v, key)
}

func (p *Store) Get(key FieldKey) (interface{}, error) {
	return p.GetWithScope("", key)
}

func (p *Store) GetWithScope(scope string, key FieldKey) (interface{}, error) {
	return p.GetWithScopeAndOverride(scope, key, nil)
}

func (p *Store) GetWithScopeAndOverride(scope string, key FieldKey, override interface{}) (interface{}, error) {
	d := p.getFieldDefinition(key)

	if d != nil {
		if e, ok := os.LookupEnv(d.EnvVar); ok {
			return e, nil
		}

		if !d.CaseSensitive {
			// use the case convention from the field definition
			key = d.Key
		}
	}

	if override != nil {
		return override, nil
	}

	res, err := p.getFromConfig(p.getPath(scope, key))
	if err != nil {
		if d != nil && d.Default != nil {
			return d.Default, nil
		}

		return nil, err
	}

	return res.Value(), nil
}

func (p *Store) Set(key FieldKey, value interface{}) error {
	return p.SetWithScope("", key, value)
}

func (p *Store) SetWithScope(scope string, key FieldKey, value interface{}) error {
	v := p.getFieldDefinition(key)

	if v != nil {
		if v.SetValidationFunc != nil {
			if err := v.SetValidationFunc(key, value); err != nil {
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

	cfg := p.getConfig()

	p.mu.Lock()
	defer p.mu.Unlock()

	cfg, err := sjson.Set(cfg, escapeWildcards(p.getPath(scope, key)), value)
	if err != nil {
		return err
	}

	if err := p.writeConfig(cfg); err != nil {
		return err
	}

	return nil
}

func (p *Store) RemoveScope(scope string) error {
	path := scope
	if p.scope != "" {
		path = fmt.Sprintf("%s.%s", p.scope, scope)
	}

	return p.deletePath(path)
}

func (p *Store) DeleteKey(key FieldKey) error {
	return p.DeleteKeyWithScope("", key)
}

func (p *Store) DeleteKeyWithScope(scope string, key FieldKey) error {
	return p.deletePath(p.getPath(scope, key))
}

func (p *Store) getPath(scope string, key FieldKey) string {
	path := string(key)
	if scope != "" {
		path = fmt.Sprintf("%s.%s", scope, key)
	}

	if p.scope != "" {
		path = fmt.Sprintf("%s.%s", p.scope, path)
	}

	return path
}

func (p *Store) deletePath(path string) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	cfg, err := sjson.Delete(p.getConfig(), escapeWildcards(path))
	if err != nil {
		return err
	}

	if err := p.writeConfig(cfg); err != nil {
		return err
	}

	return nil
}

func (p *Store) VisitAllFields(fn func(d FieldDefinition)) {
	p.VisitAllFieldsWithScope("", fn)
}

func (p *Store) VisitAllFieldsWithScope(scope string, fn func(d FieldDefinition)) {
	for _, f := range p.fields {
		fn(f)
	}
}

func (p *Store) GetScopes() []string {
	s := []string{}
	result := gjson.Get(p.getConfig(), "@this")
	result.ForEach(func(key, value gjson.Result) bool {
		s = append(s, key.String())
		return true
	})

	return s
}

func (p *Store) getFromConfig(path string) (*gjson.Result, error) {
	res := gjson.Get(p.getConfig(), path)

	if !res.Exists() {
		return nil, fmt.Errorf("no value found at path %s", path)
	}

	return &res, nil
}

func (p *Store) writeConfig(json string) error {
	p.cfg = []byte(json)

	if p.fileName != "" {
		dir := filepath.Dir(p.fileName)
		_, err := os.Stat(dir)
		if err != nil {
			err = os.Mkdir(dir, 0750)
			if err != nil {
				return err
			}
		}

		if err := ioutil.WriteFile(p.fileName, p.cfg, 0640); err != nil {
			return err
		}
	}

	return nil
}

func (p *Store) getConfig() string {
	if p.cfg == nil {
		if p.fileName != "" {
			p.setConfigFromFile()
		}
	}

	return string(p.cfg)
}

func (p *Store) setConfigFromFile() {
	data, err := ioutil.ReadFile(p.fileName)
	if err != nil {
		return
	}

	p.cfg = data
}

// Escape wildcard characters, as required by sjson
func escapeWildcards(key string) string {
	re := regexp.MustCompile(`([*?])`)
	return re.ReplaceAllString(key, "\\$1")
}

func (p *Store) getFieldDefinition(key FieldKey) *FieldDefinition {
	for _, v := range p.fields {
		if !v.CaseSensitive && strings.EqualFold(string(key), string(v.Key)) {
			return &v
		}

		if v.CaseSensitive && key == v.Key {
			return &v
		}
	}

	return nil
}

func (p *Store) getConfigValueKeys() []FieldKey {
	var keys []FieldKey
	for _, v := range p.fields {
		keys = append(keys, v.Key)
	}

	return keys
}

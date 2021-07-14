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

// JSONStore is a configurable json-backed configuration store.
type JSONStore struct {
	cfg            []byte
	fields         []FieldDefinition
	fileName       string
	scope          string
	mu             sync.Mutex
	explicitValues bool
}

// FieldKey is the key of a config field.
type FieldKey string

// FieldDefinition contains the information required to describe a configuration field.
type FieldDefinition struct {
	// EnvVar provides an environment variable override
	EnvVar string

	// Key represents the key that will be used to store the underlying value.
	Key FieldKey

	// EnvVar provides a default value to be returned during a get operation
	// if no value can be found for this field.
	Default interface{}

	// CaseSensitive determines whether this config key will be treated as case-sensitive
	// or not. When false, keys passed to get and set operations will be performed
	// with the canonical casing of the Key value provided in the field definition.
	CaseSensitive bool

	// Sensitive marks the underlying value as sensitive.  When true, the underlying
	// field's value will be obfuscated when printed to the console during the
	// execution of various commands.
	Sensitive bool

	// SetValidationFunc is a validation func that is run when a set operation
	// is performed for the underlying value. If the func returns an error, the
	// set operation will not succeed.
	SetValidationFunc FieldValueValidationFunc

	// SetValueFunc is a translation func that is run when a set operation
	// is performed for the underlying value. The value provided will be run
	// through the func provided and the resulting value will be set.
	SetValueFunc FieldValueTranslationFunc
}

// FieldValueValidationFunc is a configurable validation func that will ensure a field
// value conforms to some constraints before being set.
type FieldValueValidationFunc func(key FieldKey, value interface{}) error

// FieldValueTranslationFunc is a configurable translation func that will modify
// a value before setting it in the underlying config instance.
type FieldValueTranslationFunc func(key FieldKey, value interface{}) (interface{}, error)

// IntGreaterThan is a FieldValueValidationFunc ins a validation func that ensures
// the field value is an integer greater than 0.
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

// IsTernary is a FieldValueValidationFunc ins a validation func that ensures
// the field value is a valid Ternary.
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

// StringInStrings is a FieldValueValidationFunc ins a validation func that ensures
// the field value appears in the given collection of allowed values.
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

// ToLower is a FieldValueTranslationFunc translation func that ensures the provided
// value is case-folded to lowercase before writing to the underlying config.
func ToLower() func(key FieldKey, value interface{}) (interface{}, error) {
	return func(key FieldKey, value interface{}) (interface{}, error) {
		if s, ok := value.(string); ok {
			return strings.ToLower(s), nil
		}

		return nil, fmt.Errorf("the value %s provided for %s is not a string", value, key)
	}
}

// JSONStoreOption is a func for supplying options when creating a new JSONStore.
type JSONStoreOption func(*JSONStore) error

// NewJSONStore creates a new instance of JSONStore.
func NewJSONStore(opts ...JSONStoreOption) (*JSONStore, error) {
	p := &JSONStore{}

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

// EnforceStrictFields is a JSONStoreOption func that ensures that every field accessed
// is backed by a FieldDefinition.
func EnforceStrictFields() JSONStoreOption {
	return func(p *JSONStore) error {
		p.explicitValues = true
		return nil
	}
}

// ConfigureFields is a JSONStoreOption func that allows the caller to describe the
// fields stored in this config instance with one or more field definitions.
func ConfigureFields(definitions ...FieldDefinition) JSONStoreOption {
	return func(p *JSONStore) error {
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

// PersistToFile is a JSONStoreOption func that ensures all writes to this config
// instance are persisted to disk.
func PersistToFile(fileName string) JSONStoreOption {
	return func(p *JSONStore) error {
		p.fileName = fileName
		return nil
	}
}

// UseGlobalScope is a JSONStoreOption func that ensures all config fields are stored
// under a global object scope with the passed scope string as a key.
func UseGlobalScope(scope string) JSONStoreOption {
	return func(p *JSONStore) error {
		p.scope = scope
		return nil
	}
}

// GetString retrieves a string from this config instance for the given key.  An
// attempt will be made to convert the underlying type of this field's value to
// a string. If the value cannot be retrieved, a zero value will be returned.
func (p *JSONStore) GetString(key FieldKey) (string, error) {
	return p.GetStringWithScopeAndOverride("", key, nil)
}

// GetStringWithOverride retrieves a string from this config instance for the given
// key, overriding with the provided value if is not nil.  An attempt will be made
// to convert the underlying type of this field's value to a string. If the value
// cannot be retrieved, a zero value will be returned.
func (p *JSONStore) GetStringWithOverride(key FieldKey, override *string) (string, error) {
	return p.GetStringWithScopeAndOverride("", key, override)
}

// GetString retrieves a string from this config instance for the given key, prefixing
// the key's path with the given scope.  An attempt will be made to convert the underlying
// type of this field's value to a string. If the value cannot be retrieved, a zero
// value will be returned.
func (p *JSONStore) GetStringWithScope(scope string, key FieldKey) (string, error) {
	return p.GetStringWithScopeAndOverride(scope, key, nil)
}

// GetStringWithOverride retrieves a string from this config instance for the given
// key, prefixing the key's path with the given scope and overriding with the provided
// value if it is not nil.  An attempt will be made to convert the underlying type
// of this field's value to a string. If the value cannot be retrieved, a zero value
// will be returned.
func (p *JSONStore) GetStringWithScopeAndOverride(scope string, key FieldKey, override *string) (string, error) {
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

// GetInt retrieves an int64 from this config instance for the given key.  An attempt
// will be made to convert the underlying type of this field's value to an int64.
// If the value cannot be retrieved, a zero value will be returned.
func (p *JSONStore) GetInt(key FieldKey) (int64, error) {
	return p.GetIntWithScope("", key)
}

// GetIntWithScope retrieves an int64 from this config instance for the given key, prefixing
// the key's path with the given scope.  An attempt will be made to convert the underlying
// type of this field's value to an int64. If the value cannot be retrieved, a zero
// value will be returned.
func (p *JSONStore) GetIntWithScope(scope string, key FieldKey) (int64, error) {
	return p.GetIntWithScopeAndOverride(scope, key, nil)
}

// GetIntWithOverride retrieves an int64 from this config instance for the given
// key, prefixing the key's path with the given scope and overriding with the provided
// value if it is not nil.  An attempt will be made to convert the underlying type
// of this field's value to an int64. If the value cannot be retrieved, a zero value
// will be returned.
func (p *JSONStore) GetIntWithScopeAndOverride(scope string, key FieldKey, override *int64) (int64, error) {
	var v interface{}
	var err error

	if override == nil || *override == 0 {
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

// GetInt retrieves a Ternary from this config instance for the given key.  An attempt
// will be made to convert the underlying type of this field's value to a Ternary.
// If the value cannot be retrieved, a zero value will be returned.
func (p *JSONStore) GetTernary(key FieldKey) (Ternary, error) {
	return p.GetTernaryWithScope("", key)
}

// GetTernaryWithScope retrieves a Ternary from this config instance for the given
// key, prefixing the key's path with the given scope.  An attempt will be made to
// convert the underlying type of this field's value to a Ternary. If the value cannot
// be retrieved, a zero value will be returned.
func (p *JSONStore) GetTernaryWithScope(scope string, key FieldKey) (Ternary, error) {
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

// Get retrieves a value from this config instance for the given key.
func (p *JSONStore) Get(key FieldKey) (interface{}, error) {
	return p.GetWithScope("", key)
}

// Get retrieves a value from this config instance for the given key, prefixing
// the key's path with the given scope.
func (p *JSONStore) GetWithScope(scope string, key FieldKey) (interface{}, error) {
	return p.GetWithScopeAndOverride(scope, key, nil)
}

// Get retrieves a value from this config instance for the given key, prefixing
// the key's path with the given scope and overriding with the provided value if
// it is not nil.
func (p *JSONStore) GetWithScopeAndOverride(scope string, key FieldKey, override interface{}) (interface{}, error) {
	d := p.GetFieldDefinition(key)

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

// Set sets a value within this config instance for the given key.  The resulting
// config will be persisted to disk if PersistToDisk has been used.
func (p *JSONStore) Set(key FieldKey, value interface{}) error {
	return p.SetWithScope("", key, value)
}

// SetWithScope sets a value within this config instance for the given key, prefixing
// the key's path with the given scope. The resulting config will be persisted to
// disk if PersistToDisk has been used.
func (p *JSONStore) SetWithScope(scope string, key FieldKey, value interface{}) error {
	v := p.GetFieldDefinition(key)

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

		if v.SetValueFunc != nil {
			var err error
			value, err = v.SetValueFunc(key, value)
			if err != nil {
				return err
			}
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

	return p.writeConfig(cfg)
}

// Remove scope removes an entire scope from this config instance, including all
// the fields that appear underneath it. The resulting config will be persisted to
// disk if PersistToDisk has been used.
func (p *JSONStore) RemoveScope(scope string) error {
	path := scope
	if p.scope != "" {
		path = fmt.Sprintf("%s.%s", p.scope, scope)
	}

	return p.deletePath(path)
}

// Remove scope removes the provided key from this config instance. The resulting
// config will be persisted to disk if PersistToDisk has been used.
func (p *JSONStore) DeleteKey(key FieldKey) error {
	return p.DeleteKeyWithScope("", key)
}

// Remove scope removes the provided key from this config instance, prefixing
// the key's path with the given scope. The resulting config will be persisted to disk if PersistToDisk has been used.
func (p *JSONStore) DeleteKeyWithScope(scope string, key FieldKey) error {
	return p.deletePath(p.getPath(scope, key))
}

// ForEachFieldDefinition iterates through the defined fields for this config instance,
// yielding each to the func provided.
func (p *JSONStore) ForEachFieldDefinition(fn func(d FieldDefinition)) {
	for _, f := range p.fields {
		fn(f)
	}
}

// GetScopes returns a slice of all scopes defined within this config instance.
func (p *JSONStore) GetScopes() []string {
	s := []string{}
	result := gjson.Get(p.getConfig(), "@this")
	result.ForEach(func(key, value gjson.Result) bool {
		s = append(s, key.String())
		return true
	})

	return s
}

// GetFieldDefinition returns a field definition for the given key if one exists.
func (p *JSONStore) GetFieldDefinition(key FieldKey) *FieldDefinition {
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

func (p *JSONStore) getPath(scope string, key FieldKey) string {
	path := string(key)
	if scope != "" {
		path = fmt.Sprintf("%s.%s", scope, key)
	}

	if p.scope != "" {
		path = fmt.Sprintf("%s.%s", p.scope, path)
	}

	return path
}

func (p *JSONStore) deletePath(path string) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	cfg, err := sjson.Delete(p.getConfig(), escapeWildcards(path))
	if err != nil {
		return err
	}

	return p.writeConfig(cfg)
}

func (p *JSONStore) getFromConfig(path string) (*gjson.Result, error) {
	res := gjson.Get(p.getConfig(), path)

	if !res.Exists() {
		return nil, fmt.Errorf("no value found at path %s", path)
	}

	return &res, nil
}

func (p *JSONStore) writeConfig(json string) error {
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

func (p *JSONStore) getConfig() string {
	if p.cfg == nil {
		if p.fileName != "" {
			p.setConfigFromFile()
		}
	}

	return string(p.cfg)
}

func (p *JSONStore) setConfigFromFile() {
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

func (p *JSONStore) getConfigValueKeys() []FieldKey {
	var keys []FieldKey
	for _, v := range p.fields {
		keys = append(keys, v.Key)
	}

	return keys
}

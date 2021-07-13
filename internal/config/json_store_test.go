//+build integration

package config

import (
	"errors"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	testCfg = `{
		"*":{
			"loglevel":"debug",
			"plugindir": "/Users/ctrombley/.newrelic/plugins",
			"prereleasefeatures": "NOT_ASKED",
			"sendusagedata": "NOT_ASKED",
			"testInt": 42,
			"testString": "value1",
			"teststring": "value2"
			"caseInsensitiveTest": "value"
		}
	}`
)

func TestStore_Ctor_NilOption(t *testing.T) {
	_, err := NewJSONStore(nil)
	require.NoError(t, err)
}

func TestStore_Ctor_OptionError(t *testing.T) {
	_, err := NewJSONStore(func(*JSONStore) error { return errors.New("") })
	require.Error(t, err)
}

func TestStore_Ctor_CaseInsensitiveKeyCollision(t *testing.T) {
	_, err := NewJSONStore(
		ConfigureFields(
			FieldDefinition{Key: "asdf"},
			FieldDefinition{Key: "ASDF"},
		),
	)
	require.Error(t, err)
}

func TestStore_Ctor_CaseSensitiveKeyOverlap(t *testing.T) {
	_, err := NewJSONStore(
		ConfigureFields(
			FieldDefinition{
				Key:           "asdf",
				CaseSensitive: true,
			},
		),
	)
	require.NoError(t, err)
}

func TestStore_GetString(t *testing.T) {
	f, err := ioutil.TempFile("", "newrelic-cli.config_provider_test.*.json")
	require.NoError(t, err)
	defer f.Close()

	_, err = f.WriteString(testCfg)
	require.NoError(t, err)

	p, err := NewJSONStore(
		PersistToFile(f.Name()),
		UseGlobalScope("*"),
	)
	require.NoError(t, err)

	actual, err := p.GetString("loglevel")
	require.NoError(t, err)
	require.Equal(t, "debug", actual)
}

func TestStore_GetStringWithOverride(t *testing.T) {
	f, err := ioutil.TempFile("", "newrelic-cli.config_provider_test.*.json")
	require.NoError(t, err)
	defer f.Close()

	_, err = f.WriteString(testCfg)
	require.NoError(t, err)

	p, err := NewJSONStore(
		PersistToFile(f.Name()),
		UseGlobalScope("*"),
	)
	require.NoError(t, err)

	override := "trace"
	actual, err := p.GetStringWithOverride("loglevel", &override)
	require.NoError(t, err)
	require.Equal(t, "trace", actual)
}

func TestStore_GetStringWithScope(t *testing.T) {
	f, err := ioutil.TempFile("", "newrelic-cli.config_provider_test.*.json")
	require.NoError(t, err)
	defer f.Close()

	_, err = f.WriteString(testCfg)
	require.NoError(t, err)

	p, err := NewJSONStore(
		PersistToFile(f.Name()),
	)
	require.NoError(t, err)

	actual, err := p.GetStringWithScope("*", LogLevel)
	require.NoError(t, err)
	require.Equal(t, "debug", actual)
}

func TestStore_Get(t *testing.T) {
	f, err := ioutil.TempFile("", "newrelic-cli.config_provider_test.*.json")
	require.NoError(t, err)
	defer f.Close()

	_, err = f.WriteString(testCfg)
	require.NoError(t, err)

	p, err := NewJSONStore(
		PersistToFile(f.Name()),
		UseGlobalScope("*"),
	)
	require.NoError(t, err)

	actual, err := p.Get(LogLevel)
	require.NoError(t, err)
	require.Equal(t, "debug", actual.(string))
}

func TestStore_GetString_CaseSensitive(t *testing.T) {
	f, err := ioutil.TempFile("", "newrelic-cli.config_provider_test.*.json")
	require.NoError(t, err)
	defer f.Close()

	_, err = f.WriteString(testCfg)
	require.NoError(t, err)

	p, err := NewJSONStore(
		ConfigureFields(FieldDefinition{
			Key:           "testString",
			CaseSensitive: true,
		}),
		PersistToFile(f.Name()),
		UseGlobalScope("*"),
	)
	require.NoError(t, err)

	actual, err := p.GetString("testString")
	require.NoError(t, err)
	require.Equal(t, "value1", actual)

	actual, err = p.GetString("teststring")
	require.NoError(t, err)
	require.Equal(t, "value2", actual)
}

func TestStore_GetString_CaseInsensitive(t *testing.T) {
	f, err := ioutil.TempFile("", "newrelic-cli.config_provider_test.*.json")
	require.NoError(t, err)
	defer f.Close()

	_, err = f.WriteString(testCfg)
	require.NoError(t, err)

	p, err := NewJSONStore(
		ConfigureFields(FieldDefinition{
			Key:           "caseInsensitiveTest",
			CaseSensitive: false,
		}),
		PersistToFile(f.Name()),
		UseGlobalScope("*"),
	)
	require.NoError(t, err)

	actual, err := p.GetString("caseinsensitivetest")
	require.NoError(t, err)
	require.Equal(t, "value", actual)
}

func TestStore_GetString_NotDefined(t *testing.T) {
	f, err := ioutil.TempFile("", "newrelic-cli.config_provider_test.*.json")
	require.NoError(t, err)
	defer f.Close()

	_, err = f.WriteString(testCfg)
	require.NoError(t, err)

	p, err := NewJSONStore(
		PersistToFile(f.Name()),
		UseGlobalScope("*"),
	)
	require.NoError(t, err)

	_, err = p.GetString("undefined")
	require.Error(t, err)
	require.Contains(t, err.Error(), "no value found")
}

func TestStore_GetString_DefaultValue(t *testing.T) {
	f, err := ioutil.TempFile("", "newrelic-cli.config_provider_test.*.json")
	require.NoError(t, err)
	defer f.Close()

	p, err := NewJSONStore(
		ConfigureFields(FieldDefinition{
			Key:     "prereleasefeatures",
			Default: "NOT_ASKED",
		}),
		PersistToFile(f.Name()),
		UseGlobalScope("*"),
	)
	require.NoError(t, err)

	actual, err := p.GetString("prereleasefeatures")
	require.NoError(t, err)
	require.Equal(t, "NOT_ASKED", actual)
}

func TestStore_GetString_EnvVarOverride(t *testing.T) {
	f, err := ioutil.TempFile("", "newrelic-cli.config_provider_test.*.json")
	require.NoError(t, err)
	defer f.Close()

	p, err := NewJSONStore(
		ConfigureFields(FieldDefinition{
			Key:     "prereleasefeatures",
			Default: "NOT_ASKED",
			EnvVar:  "NEW_RELIC_CLI_PRERELEASE",
		}),
		PersistToFile(f.Name()),
		UseGlobalScope("*"),
	)
	require.NoError(t, err)

	err = os.Setenv("NEW_RELIC_CLI_PRERELEASE", "testValue")
	require.NoError(t, err)

	actual, err := p.GetString("prereleasefeatures")
	require.NoError(t, err)
	require.Equal(t, "testValue", actual)
}

func TestStore_GetInt(t *testing.T) {
	f, err := ioutil.TempFile("", "newrelic-cli.config_provider_test.*.json")
	require.NoError(t, err)
	defer f.Close()

	_, err = f.WriteString(testCfg)
	require.NoError(t, err)

	p, err := NewJSONStore(
		PersistToFile(f.Name()),
		UseGlobalScope("*"),
	)
	require.NoError(t, err)

	actual, err := p.GetInt("testInt")
	require.NoError(t, err)
	require.Equal(t, int64(42), actual)
}

func TestStore_GetInt_NotDefined(t *testing.T) {
	f, err := ioutil.TempFile("", "newrelic-cli.config_provider_test.*.json")
	require.NoError(t, err)
	defer f.Close()

	_, err = f.WriteString(testCfg)
	require.NoError(t, err)

	p, err := NewJSONStore(
		PersistToFile(f.Name()),
		UseGlobalScope("*"),
	)
	require.NoError(t, err)

	_, err = p.GetInt("undefined")
	require.Error(t, err)
	require.Contains(t, err.Error(), "no value found")
}

func TestStore_GetInt_DefaultValue(t *testing.T) {
	f, err := ioutil.TempFile("", "newrelic-cli.config_provider_test.*.json")
	require.NoError(t, err)
	defer f.Close()

	p, err := NewJSONStore(
		ConfigureFields(FieldDefinition{
			Key:     "testInt",
			Default: 42,
		}),
		PersistToFile(f.Name()),
		UseGlobalScope("*"),
	)
	require.NoError(t, err)

	actual, err := p.GetInt("testInt")
	require.NoError(t, err)
	require.Equal(t, int64(42), actual)
}

func TestStore_GetInt_EnvVarOverride(t *testing.T) {
	f, err := ioutil.TempFile("", "newrelic-cli.config_provider_test.*.json")
	require.NoError(t, err)
	defer f.Close()

	p, err := NewJSONStore(
		ConfigureFields(FieldDefinition{
			Key:    "testInt",
			EnvVar: "TEST_INT",
		}),
		PersistToFile(f.Name()),
		UseGlobalScope("*"),
	)
	require.NoError(t, err)

	err = os.Setenv("TEST_INT", "42")
	require.NoError(t, err)

	actual, err := p.GetInt("testInt")
	require.NoError(t, err)
	require.Equal(t, int64(42), actual)
}

func TestStore_GetInt_EnvVarOverride_WrongType(t *testing.T) {
	f, err := ioutil.TempFile("", "newrelic-cli.config_provider_test.*.json")
	require.NoError(t, err)
	defer f.Close()

	p, err := NewJSONStore(
		ConfigureFields(FieldDefinition{
			Key:    "testInt",
			EnvVar: "TEST_INT",
		}),
		PersistToFile(f.Name()),
		UseGlobalScope("*"),
	)
	require.NoError(t, err)

	err = os.Setenv("TEST_INT", "TEST_VALUE")
	require.NoError(t, err)

	_, err = p.GetInt("testInt")
	require.Error(t, err)
}

func TestStore_Set(t *testing.T) {
	f, err := ioutil.TempFile("", "newrelic-cli.config_provider_test.*.json")
	require.NoError(t, err)
	defer f.Close()

	_, err = f.WriteString(testCfg)
	require.NoError(t, err)

	p, err := NewJSONStore(
		PersistToFile(f.Name()),
		UseGlobalScope("*"),
	)
	require.NoError(t, err)

	err = p.Set("loglevel", "trace")
	require.NoError(t, err)

	actual, err := p.GetString("loglevel")
	require.NoError(t, err)
	require.Equal(t, "trace", actual)
}

func TestStore_SetTernary(t *testing.T) {
	f, err := ioutil.TempFile("", "newrelic-cli.config_provider_test.*.json")
	require.NoError(t, err)
	defer f.Close()

	_, err = f.WriteString(testCfg)
	require.NoError(t, err)

	p, err := NewJSONStore(
		PersistToFile(f.Name()),
		UseGlobalScope("*"),
	)
	require.NoError(t, err)

	err = p.Set("testTernary", TernaryValues.Allow)
	require.NoError(t, err)

	actual, err := p.GetTernary("testTernary")
	require.NoError(t, err)
	require.Equal(t, TernaryValues.Allow, actual)
}

func TestStore_SetTernary_Invalid(t *testing.T) {
	f, err := ioutil.TempFile("", "newrelic-cli.config_provider_test.*.json")
	require.NoError(t, err)
	defer f.Close()

	_, err = f.WriteString(testCfg)
	require.NoError(t, err)

	p, err := NewJSONStore(
		ConfigureFields(FieldDefinition{
			Key:               "testTernary",
			SetValidationFunc: IsTernary(),
		}),
		PersistToFile(f.Name()),
		UseGlobalScope("*"),
	)
	require.NoError(t, err)

	err = p.Set("testTernary", Ternary("invalid"))
	require.Error(t, err)

	err = p.Set("anotherTestTernary", "invalid")
	require.NoError(t, err)

	actual, err := p.GetTernary("anotherTestTernary")
	require.NoError(t, err)
	require.False(t, actual.Bool())
}

func TestStore_Set_CaseSensitive(t *testing.T) {
	f, err := ioutil.TempFile("", "newrelic-cli.config_provider_test.*.json")
	require.NoError(t, err)
	defer f.Close()

	_, err = f.WriteString(testCfg)
	require.NoError(t, err)

	p, err := NewJSONStore(
		EnforceStrictFields(),
		ConfigureFields(FieldDefinition{
			Key:           "loglevel",
			CaseSensitive: true,
		}),
		PersistToFile(f.Name()),
		UseGlobalScope("*"),
	)
	require.NoError(t, err)

	err = p.Set("loglevel", "trace")
	require.NoError(t, err)

	err = p.Set("logLevel", "info")
	require.Error(t, err)
}

func TestStore_Set_CaseInsensitive(t *testing.T) {
	f, err := ioutil.TempFile("", "newrelic-cli.config_provider_test.*.json")
	require.NoError(t, err)
	defer f.Close()

	_, err = f.WriteString(testCfg)
	require.NoError(t, err)

	p, err := NewJSONStore(
		EnforceStrictFields(),
		ConfigureFields(FieldDefinition{
			Key:           "loglevel",
			CaseSensitive: false,
		}),
		PersistToFile(f.Name()),
		UseGlobalScope("*"),
	)
	require.NoError(t, err)

	err = p.Set("loglevel", "trace")
	require.NoError(t, err)

	err = p.Set("LOGLEVEL", "info")
	require.NoError(t, err)

	actual, err := p.GetString("loglevel")
	require.NoError(t, err)
	require.Equal(t, "info", actual)
}

func TestStore_Set_FileDoesNotExist(t *testing.T) {
	p, err := NewJSONStore(
		UseGlobalScope("*"),
	)
	require.NoError(t, err)

	err = p.Set("loglevel", "trace")
	require.NoError(t, err)

	actual, err := p.GetString("loglevel")
	require.NoError(t, err)
	require.Equal(t, "trace", actual)
}

func TestStore_Set_ExplicitValues_CaseInsensitive(t *testing.T) {
	f, err := ioutil.TempFile("", "newrelic-cli.config_provider_test.*.json")
	require.NoError(t, err)
	defer f.Close()

	p, err := NewJSONStore(
		EnforceStrictFields(),
		ConfigureFields(FieldDefinition{
			Key: "allowed",
		}),
		PersistToFile(f.Name()),
		UseGlobalScope("*"),
	)
	require.NoError(t, err)

	err = p.Set("loglevel", "trace")
	require.Error(t, err)

	err = p.Set("ALLOWED", "testValue")
	require.NoError(t, err)
}

func TestStore_Set_ExplicitValues_CaseSensitive(t *testing.T) {
	f, err := ioutil.TempFile("", "newrelic-cli.config_provider_test.*.json")
	require.NoError(t, err)
	defer f.Close()

	p, err := NewJSONStore(
		EnforceStrictFields(),
		ConfigureFields(FieldDefinition{
			Key:           "allowed",
			CaseSensitive: true,
		}),
		PersistToFile(f.Name()),
		UseGlobalScope("*"),
	)
	require.NoError(t, err)

	err = p.Set("loglevel", "trace")
	require.Error(t, err)

	err = p.Set("ALLOWED", "testValue")
	require.Error(t, err)

	err = p.Set("allowed", "testValue")
	require.NoError(t, err)
}

func TestStore_Set_ValidationFunc_IntGreaterThan(t *testing.T) {
	f, err := ioutil.TempFile("", "newrelic-cli.config_provider_test.*.json")
	require.NoError(t, err)
	defer f.Close()

	p, err := NewJSONStore(
		ConfigureFields(FieldDefinition{
			Key:               "loglevel",
			SetValidationFunc: IntGreaterThan(0),
		}),
		PersistToFile(f.Name()),
		UseGlobalScope("*"),
	)
	require.NoError(t, err)

	err = p.Set("loglevel", 0)
	require.Error(t, err)

	err = p.Set("loglevel", 1)
	require.NoError(t, err)
}

func TestStore_Set_ValidationFunc_IntGreaterThan_WrongType(t *testing.T) {
	f, err := ioutil.TempFile("", "newrelic-cli.config_provider_test.*.json")
	require.NoError(t, err)
	defer f.Close()

	p, err := NewJSONStore(
		ConfigureFields(FieldDefinition{
			Key:               "loglevel",
			SetValidationFunc: IntGreaterThan(0),
		}),
		PersistToFile(f.Name()),
		UseGlobalScope("*"),
	)
	require.NoError(t, err)

	err = p.Set("loglevel", "debug")
	require.Error(t, err)

	err = p.Set("loglevel", 1)
	require.NoError(t, err)
}

func TestStore_Set_ValidationFunc_StringInStrings_CaseSensitive(t *testing.T) {
	f, err := ioutil.TempFile("", "newrelic-cli.config_provider_test.*.json")
	require.NoError(t, err)
	defer f.Close()

	p, err := NewJSONStore(
		ConfigureFields(FieldDefinition{
			Key:               "loglevel",
			SetValidationFunc: StringInStrings(true, "valid", "alsoValid"),
		}),
		PersistToFile(f.Name()),
		UseGlobalScope("*"),
	)
	require.NoError(t, err)

	err = p.Set("loglevel", "trace")
	require.Error(t, err)

	err = p.Set("loglevel", "valid")
	require.NoError(t, err)
}

func TestStore_Set_ValidationFunc_StringInStrings_CaseInsensitive(t *testing.T) {
	f, err := ioutil.TempFile("", "newrelic-cli.config_provider_test.*.json")
	require.NoError(t, err)
	defer f.Close()

	p, err := NewJSONStore(
		ConfigureFields(FieldDefinition{
			Key:               "loglevel",
			SetValidationFunc: StringInStrings(false, "valid", "alsoValid"),
		}),
		PersistToFile(f.Name()),
		UseGlobalScope("*"),
	)
	require.NoError(t, err)

	err = p.Set("loglevel", "VALID")
	require.NoError(t, err)

	err = p.Set("loglevel", "ALSOVALID")
	require.NoError(t, err)
}

func TestStore_Set_ValidationFunc_StringInStrings_WrongType(t *testing.T) {
	f, err := ioutil.TempFile("", "newrelic-cli.config_provider_test.*.json")
	require.NoError(t, err)
	defer f.Close()

	p, err := NewJSONStore(
		ConfigureFields(FieldDefinition{
			Key:               "testInt",
			SetValidationFunc: StringInStrings(false, "valid", "alsoValid"),
		}),
		PersistToFile(f.Name()),
		UseGlobalScope("*"),
	)
	require.NoError(t, err)

	err = p.Set("testInt", 42)
	require.Error(t, err)
	require.Contains(t, err.Error(), "is not a string")
}

func TestStore_Set_ValueFunc_ToLower(t *testing.T) {
	f, err := ioutil.TempFile("", "newrelic-cli.config_provider_test.*.json")
	require.NoError(t, err)
	defer f.Close()

	p, err := NewJSONStore(
		ConfigureFields(FieldDefinition{
			Key:          "testString",
			SetValueFunc: ToLower(),
		}),
		PersistToFile(f.Name()),
	)
	require.NoError(t, err)

	err = p.Set("testString", "TEST_VALUE")
	require.NoError(t, err)

	v, err := p.GetString("testString")
	require.NoError(t, err)

	require.Equal(t, "test_value", v)
}

func TestStore_RemoveScope(t *testing.T) {
	f, err := ioutil.TempFile("", "newrelic-cli.config_provider_test.*.json")
	require.NoError(t, err)
	defer f.Close()

	p, err := NewJSONStore(
		ConfigureFields(FieldDefinition{
			Key: "testString",
		}),
	)
	require.NoError(t, err)

	err = p.RemoveScope("testString")
	require.NoError(t, err)

	_, err = p.GetString("testString")
	require.Error(t, err)
}

func TestStore_RemoveScope_GlobalScope(t *testing.T) {
	f, err := ioutil.TempFile("", "newrelic-cli.config_provider_test.*.json")
	require.NoError(t, err)
	defer f.Close()

	p, err := NewJSONStore(
		ConfigureFields(FieldDefinition{
			Key: "testString",
		}),
		UseGlobalScope("*"),
	)
	require.NoError(t, err)

	err = p.SetWithScope("scope", FieldKey("testString"), "testValue")
	require.NoError(t, err)

	err = p.RemoveScope("scope")
	require.NoError(t, err)

	_, err = p.GetStringWithScope("scope", "testString")
	require.Error(t, err)
}

func TestStore_DeleteKey(t *testing.T) {
	f, err := ioutil.TempFile("", "newrelic-cli.config_provider_test.*.json")
	require.NoError(t, err)
	defer f.Close()

	p, err := NewJSONStore(
		ConfigureFields(FieldDefinition{
			Key: "testString",
		}),
	)
	require.NoError(t, err)

	err = p.DeleteKey("testString")
	require.NoError(t, err)

	_, err = p.GetString("testString")
	require.Error(t, err)
}

func TestStore_ForEachFieldDefinition(t *testing.T) {
	f, err := ioutil.TempFile("", "newrelic-cli.config_provider_test.*.json")
	require.NoError(t, err)
	defer f.Close()

	p, err := NewJSONStore(
		ConfigureFields(
			FieldDefinition{Key: "1"},
			FieldDefinition{Key: "2"},
			FieldDefinition{Key: "3"},
		),
	)
	require.NoError(t, err)

	count := 0
	fn := func(fd FieldDefinition) { count++ }
	p.ForEachFieldDefinition(fn)
	require.Equal(t, 3, count)
}

func TestStore_GetScopes(t *testing.T) {
	f, err := ioutil.TempFile("", "newrelic-cli.config_provider_test.*.json")
	require.NoError(t, err)
	defer f.Close()

	p, err := NewJSONStore()
	require.NoError(t, err)

	err = p.SetWithScope("scope1", FieldKey("testKey"), "testValue")
	require.NoError(t, err)

	err = p.SetWithScope("scope2", FieldKey("testKey"), "testValue")
	require.NoError(t, err)

	s := p.GetScopes()
	require.Equal(t, 2, len(s))
}

# Testing YAML Validation

## Important: Rebuild After YAML Changes!

**The YAML files are embedded into the binary at compile time using `//go:embed`.**
If you change a YAML file, you MUST rebuild the binary for the changes to take effect.

## Note: Clean JSON Output

All commands now return JSON with status/error fields. Validation errors will appear in the error field:

```json
{
  "status": "failed",
  "error": "invalid value 'KUBERNETESCLUSTER' for flag --managed-entity-type: must be one of [HOST, KUBRNETESCLUSTER]"
}
```

## Test Steps

### 1. Verify Current State

First, check what's currently in the YAML:

```bash
cat internal/fleetcontrol/configs/create.yaml | grep -A 3 "managed-entity-type"
```

### 2. Rebuild the Binary

```bash
# From the repository root
go build -o ./bin/darwin/newrelic ./cmd/newrelic
```

### 3. Test with Valid Value

```bash
# This should work (assuming YAML has correct values)
./bin/darwin/newrelic fleetcontrol fleet create \
  --name "Test Fleet" \
  --managed-entity-type KUBERNETESCLUSTER
```

**Expected result**: Command proceeds (validation passes silently)

### 4. Test with Typo in YAML

Edit `internal/fleetcontrol/configs/create.yaml` line 26:

```yaml
# Change from:
allowed_values: ["HOST", "KUBERNETESCLUSTER"]

# To (intentional typo):
allowed_values: ["HOST", "KUBRNETESCLUSTER"]
```

**REBUILD** (critical step!):

```bash
go build -o ./bin/darwin/newrelic ./cmd/newrelic
```

Now test with the CORRECT spelling:

```bash
./bin/darwin/newrelic fleetcontrol fleet create \
  --name "Test Fleet" \
  --managed-entity-type KUBERNETESCLUSTER
```

**Expected result**: Should FAIL with validation error in JSON format:
```json
{
  "status": "failed",
  "error": "invalid value 'KUBERNETESCLUSTER' for flag --managed-entity-type: must be one of [HOST, KUBRNETESCLUSTER]"
}
```

This proves validation is working! The correct value was rejected because the YAML has the typo.

**Using jq to check validation:**
```bash
./bin/darwin/newrelic fleetcontrol fleet create \
  --name "Test Fleet" \
  --managed-entity-type KUBERNETESCLUSTER | jq -r '.status'
# Output: "failed"

./bin/darwin/newrelic fleetcontrol fleet create \
  --name "Test Fleet" \
  --managed-entity-type KUBERNETESCLUSTER | jq -r '.error'
# Output: "invalid value 'KUBERNETESCLUSTER' for flag --managed-entity-type: must be one of [HOST, KUBRNETESCLUSTER]"
```

### 5. Test with Typo Value

Test with the TYPO value (that's in the YAML):

```bash
./bin/darwin/newrelic fleetcontrol fleet create \
  --name "Test Fleet" \
  --managed-entity-type KUBRNETESCLUSTER
```

**Expected result**: Should PASS validation (because it matches the YAML), but FAIL in the mapper:
```json
{
  "status": "failed",
  "error": "unrecognized managed entity type 'KUBRNETESCLUSTER' - YAML validation may have failed"
}
```

This shows the two-layer validation:
1. **Framework validation** (YAML rules) - passed because "KUBRNETESCLUSTER" is in allowed_values
2. **Mapper validation** (business logic) - failed because the mapper doesn't recognize the typo

**Check which validation failed:**
```bash
# Framework validation error will mention "must be one of"
./bin/darwin/newrelic fleetcontrol fleet create ... | jq -r '.error' | grep "must be one of"

# Mapper validation error will mention "unrecognized"
./bin/darwin/newrelic fleetcontrol fleet create ... | jq -r '.error' | grep "unrecognized"
```

## How Validation Works

The validation happens in two stages:

1. **YAML Framework Validation** (`command_framework.go`):
   - Checks if the provided value is in `allowed_values`
   - Runs automatically before the command handler
   - Produces error: `invalid value 'X' for flag --Y: must be one of [...]`

2. **Mapper Validation** (`helpers.go`):
   - Converts validated strings to client library types
   - Acts as a safety check
   - Produces error: `unrecognized X type 'Y' - YAML validation may have failed`

If mapper validation fails but YAML validation passed, it means the YAML `allowed_values`
don't match what the mapper expects - they need to be synchronized.

## Common Issues

### Issue: "Validation passes when it should fail"

**Cause**: You didn't rebuild after changing YAML

**Solution**:
```bash
go build -o ./bin/darwin/newrelic ./cmd/newrelic
```

### Issue: "Binary seems to use old YAML values"

**Cause**: Binary still has old embedded YAML

**Solution**: Clean and rebuild:
```bash
rm ./bin/darwin/newrelic
go clean -cache
go build -o ./bin/darwin/newrelic ./cmd/newrelic
```

### Issue: "Mapper error but validation passed"

**Cause**: YAML `allowed_values` and mapper function are out of sync

**Solution**: Update YAML to match the mapper, or vice versa. They must be consistent.

## Verification Checklist

- [ ] YAML file edited with typo: `allowed_values: ["HOST", "KUBRNETESCLUSTER"]`
- [ ] Binary rebuilt after YAML change: `go build -o ./bin/darwin/newrelic ./cmd/newrelic`
- [ ] Correct spelling rejected with error: `Error: invalid value 'KUBERNETESCLUSTER' for flag --managed-entity-type: must be one of [HOST, KUBRNETESCLUSTER]`
- [ ] Typo spelling passes framework validation but mapper rejects it: `Error: unrecognized managed entity type 'KUBRNETESCLUSTER'`
- [ ] After fixing YAML and rebuilding, correct spelling works again

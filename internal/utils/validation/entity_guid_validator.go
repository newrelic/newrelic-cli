package validation

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const (
	DELIMITER = "|"
)

var (
	invalidEntityGuidErr = errors.New("invalid entity GUID format")
	emptyEntityTypeErr   = errors.New("entity type is required")
	emptyDomainIdErr     = errors.New("domain ID is required")
)

// GenericEntity represents the decoded entity information
type GenericEntity struct {
	AccountId  int64  `json:"accountId"`
	Domain     string `json:"domain"`
	EntityType string `json:"entityType"`
	DomainId   string `json:"domainId"`
}

// DecodeEntityGuid decodes a string representation of an entity GUID and returns an GenericEntity (replaced with struct)
func DecodeEntityGuid(entityGuid string) (*GenericEntity, error) {
	decodedGuid, err := base64.StdEncoding.DecodeString(entityGuid)
	if err != nil {
		return nil, invalidEntityGuidErr
	}

	parts := strings.Split(string(decodedGuid), "|")
	if len(parts) < 4 {
		return nil, fmt.Errorf("invalid entity GUID format: expected at least 4 parts delimited by '%s': %s", DELIMITER, entityGuid)
	}

	accountId, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid account ID: %w", err)
	}

	domain := parts[1]
	entityType := parts[2]
	domainId := parts[3]

	if entityType == "" {
		return nil, emptyEntityTypeErr
	}

	if domainId == "" {
		return nil, emptyDomainIdErr
	}

	return &GenericEntity{
		AccountId:  accountId,
		Domain:     domain,
		EntityType: entityType,
		DomainId:   domainId,
	}, nil
}

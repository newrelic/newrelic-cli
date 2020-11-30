package install

type nerdGraphClient interface {
	QueryWithResponse(query string, variables map[string]interface{}, respBody interface{}) error
}

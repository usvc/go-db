package db

import (
	"fmt"
	"net/url"
	"strings"
)

// generateDSN returns the data source name that is represented in a
// structured format by the connection :options parameter
func generateDSN(options Options) string {
	var params url.URL
	query := params.Query()
	query.Add("parseTime", "true")
	for key, value := range options.Params {
		query.Add(key, value)
	}
	connectionParams := query.Encode()
	return strings.Trim(fmt.Sprintf(
		"%s:%s@tcp(%s:%v)/%s?%s",
		options.Username,
		options.Password,
		options.Hostname,
		options.Port,
		options.Database,
		connectionParams,
	), "/")
}

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
	var dsnFormat string
	switch options.Driver {
	case DriverPostgreSQL:
		dsnFormat = "postgresql://%s:%s@%s:%v/%s?%s"
	case DriverMSSQL:
		dsnFormat = "sqlserver://%s:%s@%s:%v/%s?%s"
	case DriverMySQL:
		fallthrough
	default:
		dsnFormat = "%s:%s@tcp(%s:%v)/%s?%s"
	}
	return strings.Trim(fmt.Sprintf(
		dsnFormat,
		options.Username,
		options.Password,
		options.Hostname,
		options.Port,
		options.Database,
		connectionParams,
	), "/")
}

func stringNotSet(test string) bool {
	var zeroValue string
	return test == zeroValue
}

func uint16NotSet(test uint16) bool {
	var zeroValue uint16
	return test == zeroValue
}

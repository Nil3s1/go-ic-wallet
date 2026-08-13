package kurrentdbstore

import "github.com/kurrent-io/KurrentDB-Client-Go/kurrentdb"

func NewClient(connectionString string) (*kurrentdb.Client, error) {
	settings, err := kurrentdb.ParseConnectionString(connectionString)
	if err != nil {
		return nil, err
	}

	return kurrentdb.NewClient(settings)
}

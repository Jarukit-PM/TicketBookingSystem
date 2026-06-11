package db

import (
	"net/url"
	"strings"

	"go.mongodb.org/mongo-driver/mongo"
)

const defaultDatabaseName = "tbs"

// DatabaseNameFromURI extracts the database name from a MongoDB URI path, or returns the default.
func DatabaseNameFromURI(uri string) string {
	u, err := url.Parse(uri)
	if err != nil || u.Path == "" || u.Path == "/" {
		return defaultDatabaseName
	}
	name := strings.TrimPrefix(u.Path, "/")
	if idx := strings.Index(name, "?"); idx >= 0 {
		name = name[:idx]
	}
	if name == "" {
		return defaultDatabaseName
	}
	return name
}

// Database returns the configured MongoDB database handle.
func Database(client *mongo.Client, uri string) *mongo.Database {
	return client.Database(DatabaseNameFromURI(uri))
}

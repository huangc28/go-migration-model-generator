package main

import "github.com/huangc28/go-migration-model-generator/internal/genmodel"

// This package uses SQL parser from `sqlc` to generate go code models. It reads the content from all migrations in `db/migrations`. All migration files are prefixed with current datetime stamp as the version number. For instance, we got the most up to date version of migration from DB:
//
//   version: 7
//   dirty: false
//
// when we run:
//   go run cmd/main.go gen
//
// we will collect content from migration files from 1 ~ 7 merge them in a master file `db/schema.sql` in which we will generate our go code from.
func main() {
	genmodel.Execute()
}

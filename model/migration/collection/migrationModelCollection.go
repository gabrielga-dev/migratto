package migration_collection_model

import (
	"fmt"

	migration_model "github.com/gabrielga-dev/migratto/model/migration"
)

type MigrationModelCollection struct {
	Migrations []migration_model.MigrationModel
}

func (m MigrationModelCollection) GetMigrationByTag(tag string) (migration_model.MigrationModel, error) {
	for _, migration := range m.Migrations {
		if migration.Tag == tag {
			return migration, nil
		}
	}
	return migration_model.MigrationModel{}, fmt.Errorf("migration with tag %s not found", tag)
}

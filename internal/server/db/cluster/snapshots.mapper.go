//go:build linux && cgo && !agent

package cluster

// The code below was generated by incus-generate - DO NOT EDIT!

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/lxc/incus/v6/internal/server/db/query"
	"github.com/lxc/incus/v6/shared/api"
)

var _ = api.ServerEnvironment{}

var instanceSnapshotObjects = RegisterStmt(`
SELECT instances_snapshots.id, projects.name AS project, instances.name AS instance, instances_snapshots.name, instances_snapshots.creation_date, instances_snapshots.stateful, coalesce(instances_snapshots.description, ''), instances_snapshots.expiry_date
  FROM instances_snapshots
  JOIN projects ON instances.project_id = projects.id
  JOIN instances ON instances_snapshots.instance_id = instances.id
  ORDER BY projects.id, instances.id, instances_snapshots.name
`)

var instanceSnapshotObjectsByID = RegisterStmt(`
SELECT instances_snapshots.id, projects.name AS project, instances.name AS instance, instances_snapshots.name, instances_snapshots.creation_date, instances_snapshots.stateful, coalesce(instances_snapshots.description, ''), instances_snapshots.expiry_date
  FROM instances_snapshots
  JOIN projects ON instances.project_id = projects.id
  JOIN instances ON instances_snapshots.instance_id = instances.id
  WHERE ( instances_snapshots.id = ? )
  ORDER BY projects.id, instances.id, instances_snapshots.name
`)

var instanceSnapshotObjectsByProjectAndInstance = RegisterStmt(`
SELECT instances_snapshots.id, projects.name AS project, instances.name AS instance, instances_snapshots.name, instances_snapshots.creation_date, instances_snapshots.stateful, coalesce(instances_snapshots.description, ''), instances_snapshots.expiry_date
  FROM instances_snapshots
  JOIN projects ON instances.project_id = projects.id
  JOIN instances ON instances_snapshots.instance_id = instances.id
  WHERE ( project = ? AND instance = ? )
  ORDER BY projects.id, instances.id, instances_snapshots.name
`)

var instanceSnapshotObjectsByProjectAndInstanceAndName = RegisterStmt(`
SELECT instances_snapshots.id, projects.name AS project, instances.name AS instance, instances_snapshots.name, instances_snapshots.creation_date, instances_snapshots.stateful, coalesce(instances_snapshots.description, ''), instances_snapshots.expiry_date
  FROM instances_snapshots
  JOIN projects ON instances.project_id = projects.id
  JOIN instances ON instances_snapshots.instance_id = instances.id
  WHERE ( project = ? AND instance = ? AND instances_snapshots.name = ? )
  ORDER BY projects.id, instances.id, instances_snapshots.name
`)

var instanceSnapshotID = RegisterStmt(`
SELECT instances_snapshots.id FROM instances_snapshots
  JOIN projects ON instances.project_id = projects.id
  JOIN instances ON instances_snapshots.instance_id = instances.id
  WHERE projects.name = ? AND instances.name = ? AND instances_snapshots.name = ?
`)

var instanceSnapshotCreate = RegisterStmt(`
INSERT INTO instances_snapshots (instance_id, name, creation_date, stateful, description, expiry_date)
  VALUES ((SELECT instances.id FROM instances JOIN projects ON instances.project_id = projects.id WHERE projects.name = ? AND instances.name = ?), ?, ?, ?, ?, ?)
`)

var instanceSnapshotRename = RegisterStmt(`
UPDATE instances_snapshots SET name = ? WHERE instance_id = (SELECT instances.id FROM instances JOIN projects ON instances.project_id = projects.id WHERE projects.name = ? AND instances.name = ?) AND name = ?
`)

var instanceSnapshotDeleteByProjectAndInstanceAndName = RegisterStmt(`
DELETE FROM instances_snapshots WHERE instance_id = (SELECT instances.id FROM instances JOIN projects ON instances.project_id = projects.id WHERE projects.name = ? AND instances.name = ?) AND name = ?
`)

// instanceSnapshotColumns returns a string of column names to be used with a SELECT statement for the entity.
// Use this function when building statements to retrieve database entries matching the InstanceSnapshot entity.
func instanceSnapshotColumns() string {
	return "instances_snapshots.id, projects.name AS project, instances.name AS instance, instances_snapshots.name, instances_snapshots.creation_date, instances_snapshots.stateful, coalesce(instances_snapshots.description, ''), instances_snapshots.expiry_date"
}

// getInstanceSnapshots can be used to run handwritten sql.Stmts to return a slice of objects.
func getInstanceSnapshots(ctx context.Context, stmt *sql.Stmt, args ...any) ([]InstanceSnapshot, error) {
	objects := make([]InstanceSnapshot, 0)

	dest := func(scan func(dest ...any) error) error {
		i := InstanceSnapshot{}
		err := scan(&i.ID, &i.Project, &i.Instance, &i.Name, &i.CreationDate, &i.Stateful, &i.Description, &i.ExpiryDate)
		if err != nil {
			return err
		}

		objects = append(objects, i)

		return nil
	}

	err := query.SelectObjects(ctx, stmt, dest, args...)
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch from \"instances_snapshots\" table: %w", err)
	}

	return objects, nil
}

// getInstanceSnapshotsRaw can be used to run handwritten query strings to return a slice of objects.
func getInstanceSnapshotsRaw(ctx context.Context, tx *sql.Tx, sql string, args ...any) ([]InstanceSnapshot, error) {
	objects := make([]InstanceSnapshot, 0)

	dest := func(scan func(dest ...any) error) error {
		i := InstanceSnapshot{}
		err := scan(&i.ID, &i.Project, &i.Instance, &i.Name, &i.CreationDate, &i.Stateful, &i.Description, &i.ExpiryDate)
		if err != nil {
			return err
		}

		objects = append(objects, i)

		return nil
	}

	err := query.Scan(ctx, tx, sql, dest, args...)
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch from \"instances_snapshots\" table: %w", err)
	}

	return objects, nil
}

// GetInstanceSnapshots returns all available instance_snapshots.
// generator: instance_snapshot GetMany
func GetInstanceSnapshots(ctx context.Context, tx *sql.Tx, filters ...InstanceSnapshotFilter) ([]InstanceSnapshot, error) {
	var err error

	// Result slice.
	objects := make([]InstanceSnapshot, 0)

	// Pick the prepared statement and arguments to use based on active criteria.
	var sqlStmt *sql.Stmt
	args := []any{}
	queryParts := [2]string{}

	if len(filters) == 0 {
		sqlStmt, err = Stmt(tx, instanceSnapshotObjects)
		if err != nil {
			return nil, fmt.Errorf("Failed to get \"instanceSnapshotObjects\" prepared statement: %w", err)
		}
	}

	for i, filter := range filters {
		if filter.Project != nil && filter.Instance != nil && filter.Name != nil && filter.ID == nil {
			args = append(args, []any{filter.Project, filter.Instance, filter.Name}...)
			if len(filters) == 1 {
				sqlStmt, err = Stmt(tx, instanceSnapshotObjectsByProjectAndInstanceAndName)
				if err != nil {
					return nil, fmt.Errorf("Failed to get \"instanceSnapshotObjectsByProjectAndInstanceAndName\" prepared statement: %w", err)
				}

				break
			}

			query, err := StmtString(instanceSnapshotObjectsByProjectAndInstanceAndName)
			if err != nil {
				return nil, fmt.Errorf("Failed to get \"instanceSnapshotObjects\" prepared statement: %w", err)
			}

			parts := strings.SplitN(query, "ORDER BY", 2)
			if i == 0 {
				copy(queryParts[:], parts)
				continue
			}

			_, where, _ := strings.Cut(parts[0], "WHERE")
			queryParts[0] += "OR" + where
		} else if filter.Project != nil && filter.Instance != nil && filter.ID == nil && filter.Name == nil {
			args = append(args, []any{filter.Project, filter.Instance}...)
			if len(filters) == 1 {
				sqlStmt, err = Stmt(tx, instanceSnapshotObjectsByProjectAndInstance)
				if err != nil {
					return nil, fmt.Errorf("Failed to get \"instanceSnapshotObjectsByProjectAndInstance\" prepared statement: %w", err)
				}

				break
			}

			query, err := StmtString(instanceSnapshotObjectsByProjectAndInstance)
			if err != nil {
				return nil, fmt.Errorf("Failed to get \"instanceSnapshotObjects\" prepared statement: %w", err)
			}

			parts := strings.SplitN(query, "ORDER BY", 2)
			if i == 0 {
				copy(queryParts[:], parts)
				continue
			}

			_, where, _ := strings.Cut(parts[0], "WHERE")
			queryParts[0] += "OR" + where
		} else if filter.ID != nil && filter.Project == nil && filter.Instance == nil && filter.Name == nil {
			args = append(args, []any{filter.ID}...)
			if len(filters) == 1 {
				sqlStmt, err = Stmt(tx, instanceSnapshotObjectsByID)
				if err != nil {
					return nil, fmt.Errorf("Failed to get \"instanceSnapshotObjectsByID\" prepared statement: %w", err)
				}

				break
			}

			query, err := StmtString(instanceSnapshotObjectsByID)
			if err != nil {
				return nil, fmt.Errorf("Failed to get \"instanceSnapshotObjects\" prepared statement: %w", err)
			}

			parts := strings.SplitN(query, "ORDER BY", 2)
			if i == 0 {
				copy(queryParts[:], parts)
				continue
			}

			_, where, _ := strings.Cut(parts[0], "WHERE")
			queryParts[0] += "OR" + where
		} else if filter.ID == nil && filter.Project == nil && filter.Instance == nil && filter.Name == nil {
			return nil, fmt.Errorf("Cannot filter on empty InstanceSnapshotFilter")
		} else {
			return nil, fmt.Errorf("No statement exists for the given Filter")
		}
	}

	// Select.
	if sqlStmt != nil {
		objects, err = getInstanceSnapshots(ctx, sqlStmt, args...)
	} else {
		queryStr := strings.Join(queryParts[:], "ORDER BY")
		objects, err = getInstanceSnapshotsRaw(ctx, tx, queryStr, args...)
	}

	if err != nil {
		return nil, fmt.Errorf("Failed to fetch from \"instances_snapshots\" table: %w", err)
	}

	return objects, nil
}

// GetInstanceSnapshotDevices returns all available InstanceSnapshot Devices
// generator: instance_snapshot GetMany
func GetInstanceSnapshotDevices(ctx context.Context, tx *sql.Tx, instanceSnapshotID int, filters ...DeviceFilter) (map[string]Device, error) {
	instanceSnapshotDevices, err := GetDevices(ctx, tx, "instance_snapshot", filters...)
	if err != nil {
		return nil, err
	}

	devices := map[string]Device{}
	for _, ref := range instanceSnapshotDevices[instanceSnapshotID] {
		_, ok := devices[ref.Name]
		if !ok {
			devices[ref.Name] = ref
		} else {
			return nil, fmt.Errorf("Found duplicate Device with name %q", ref.Name)
		}
	}

	return devices, nil
}

// GetInstanceSnapshotConfig returns all available InstanceSnapshot Config
// generator: instance_snapshot GetMany
func GetInstanceSnapshotConfig(ctx context.Context, tx *sql.Tx, instanceSnapshotID int, filters ...ConfigFilter) (map[string]string, error) {
	instanceSnapshotConfig, err := GetConfig(ctx, tx, "instance_snapshot", filters...)
	if err != nil {
		return nil, err
	}

	config, ok := instanceSnapshotConfig[instanceSnapshotID]
	if !ok {
		config = map[string]string{}
	}

	return config, nil
}

// GetInstanceSnapshot returns the instance_snapshot with the given key.
// generator: instance_snapshot GetOne
func GetInstanceSnapshot(ctx context.Context, tx *sql.Tx, project string, instance string, name string) (*InstanceSnapshot, error) {
	filter := InstanceSnapshotFilter{}
	filter.Project = &project
	filter.Instance = &instance
	filter.Name = &name

	objects, err := GetInstanceSnapshots(ctx, tx, filter)
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch from \"instances_snapshots\" table: %w", err)
	}

	switch len(objects) {
	case 0:
		return nil, api.StatusErrorf(http.StatusNotFound, "InstanceSnapshot not found")
	case 1:
		return &objects[0], nil
	default:
		return nil, fmt.Errorf("More than one \"instances_snapshots\" entry matches")
	}
}

// GetInstanceSnapshotID return the ID of the instance_snapshot with the given key.
// generator: instance_snapshot ID
func GetInstanceSnapshotID(ctx context.Context, tx *sql.Tx, project string, instance string, name string) (int64, error) {
	stmt, err := Stmt(tx, instanceSnapshotID)
	if err != nil {
		return -1, fmt.Errorf("Failed to get \"instanceSnapshotID\" prepared statement: %w", err)
	}

	row := stmt.QueryRowContext(ctx, project, instance, name)
	var id int64
	err = row.Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		return -1, api.StatusErrorf(http.StatusNotFound, "InstanceSnapshot not found")
	}

	if err != nil {
		return -1, fmt.Errorf("Failed to get \"instances_snapshots\" ID: %w", err)
	}

	return id, nil
}

// InstanceSnapshotExists checks if a instance_snapshot with the given key exists.
// generator: instance_snapshot Exists
func InstanceSnapshotExists(ctx context.Context, tx *sql.Tx, project string, instance string, name string) (bool, error) {
	_, err := GetInstanceSnapshotID(ctx, tx, project, instance, name)
	if err != nil {
		if api.StatusErrorCheck(err, http.StatusNotFound) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

// CreateInstanceSnapshot adds a new instance_snapshot to the database.
// generator: instance_snapshot Create
func CreateInstanceSnapshot(ctx context.Context, tx *sql.Tx, object InstanceSnapshot) (int64, error) {
	// Check if a instance_snapshot with the same key exists.
	exists, err := InstanceSnapshotExists(ctx, tx, object.Project, object.Instance, object.Name)
	if err != nil {
		return -1, fmt.Errorf("Failed to check for duplicates: %w", err)
	}

	if exists {
		return -1, api.StatusErrorf(http.StatusConflict, "This \"instances_snapshots\" entry already exists")
	}

	args := make([]any, 7)

	// Populate the statement arguments.
	args[0] = object.Project
	args[1] = object.Instance
	args[2] = object.Name
	args[3] = object.CreationDate
	args[4] = object.Stateful
	args[5] = object.Description
	args[6] = object.ExpiryDate

	// Prepared statement to use.
	stmt, err := Stmt(tx, instanceSnapshotCreate)
	if err != nil {
		return -1, fmt.Errorf("Failed to get \"instanceSnapshotCreate\" prepared statement: %w", err)
	}

	// Execute the statement.
	result, err := stmt.Exec(args...)
	if err != nil {
		return -1, fmt.Errorf("Failed to create \"instances_snapshots\" entry: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return -1, fmt.Errorf("Failed to fetch \"instances_snapshots\" entry ID: %w", err)
	}

	return id, nil
}

// CreateInstanceSnapshotDevices adds new instance_snapshot Devices to the database.
// generator: instance_snapshot Create
func CreateInstanceSnapshotDevices(ctx context.Context, tx *sql.Tx, instanceSnapshotID int64, devices map[string]Device) error {
	for key, device := range devices {
		device.ReferenceID = int(instanceSnapshotID)
		devices[key] = device
	}

	err := CreateDevices(ctx, tx, "instance_snapshot", devices)
	if err != nil {
		return fmt.Errorf("Insert Device failed for InstanceSnapshot: %w", err)
	}

	return nil
}

// CreateInstanceSnapshotConfig adds new instance_snapshot Config to the database.
// generator: instance_snapshot Create
func CreateInstanceSnapshotConfig(ctx context.Context, tx *sql.Tx, instanceSnapshotID int64, config map[string]string) error {
	referenceID := int(instanceSnapshotID)
	for key, value := range config {
		insert := Config{
			ReferenceID: referenceID,
			Key:         key,
			Value:       value,
		}

		err := CreateConfig(ctx, tx, "instance_snapshot", insert)
		if err != nil {
			return fmt.Errorf("Insert Config failed for InstanceSnapshot: %w", err)
		}

	}

	return nil
}

// RenameInstanceSnapshot renames the instance_snapshot matching the given key parameters.
// generator: instance_snapshot Rename
func RenameInstanceSnapshot(ctx context.Context, tx *sql.Tx, project string, instance string, name string, to string) error {
	stmt, err := Stmt(tx, instanceSnapshotRename)
	if err != nil {
		return fmt.Errorf("Failed to get \"instanceSnapshotRename\" prepared statement: %w", err)
	}

	result, err := stmt.Exec(to, project, instance, name)
	if err != nil {
		return fmt.Errorf("Rename InstanceSnapshot failed: %w", err)
	}

	n, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("Fetch affected rows failed: %w", err)
	}

	if n != 1 {
		return fmt.Errorf("Query affected %d rows instead of 1", n)
	}

	return nil
}

// DeleteInstanceSnapshot deletes the instance_snapshot matching the given key parameters.
// generator: instance_snapshot DeleteOne-by-Project-and-Instance-and-Name
func DeleteInstanceSnapshot(ctx context.Context, tx *sql.Tx, project string, instance string, name string) error {
	stmt, err := Stmt(tx, instanceSnapshotDeleteByProjectAndInstanceAndName)
	if err != nil {
		return fmt.Errorf("Failed to get \"instanceSnapshotDeleteByProjectAndInstanceAndName\" prepared statement: %w", err)
	}

	result, err := stmt.Exec(project, instance, name)
	if err != nil {
		return fmt.Errorf("Delete \"instances_snapshots\": %w", err)
	}

	n, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("Fetch affected rows: %w", err)
	}

	if n == 0 {
		return api.StatusErrorf(http.StatusNotFound, "InstanceSnapshot not found")
	} else if n > 1 {
		return fmt.Errorf("Query deleted %d InstanceSnapshot rows instead of 1", n)
	}

	return nil
}

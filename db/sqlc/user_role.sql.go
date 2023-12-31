// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0
// source: user_role.sql

package db

import (
	"context"

	"github.com/google/uuid"
)

const createRole = `-- name: CreateRole :one
INSERT INTO user_roles (name, user_id)
VALUES ($1, $2)
RETURNING id, name, user_id, created_at, updated_at
`

type CreateRoleParams struct {
	Name   string    `json:"name"`
	UserID uuid.UUID `json:"user_id"`
}

func (q *Queries) CreateRole(ctx context.Context, arg CreateRoleParams) (UserRole, error) {
	row := q.db.QueryRowContext(ctx, createRole, arg.Name, arg.UserID)
	var i UserRole
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.UserID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteRole = `-- name: DeleteRole :exec
DELETE FROM user_roles
WHERE id = $1
`

func (q *Queries) DeleteRole(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, deleteRole, id)
	return err
}

const getRole = `-- name: GetRole :one
SELECT id, name, user_id, created_at, updated_at
FROM user_roles
WHERE id = $1
`

func (q *Queries) GetRole(ctx context.Context, id uuid.UUID) (UserRole, error) {
	row := q.db.QueryRowContext(ctx, getRole, id)
	var i UserRole
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.UserID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getRoles = `-- name: GetRoles :many
SELECT id, name, user_id, created_at, updated_at
FROM user_roles
WHERE user_id = $1
`

func (q *Queries) GetRoles(ctx context.Context, userID uuid.UUID) ([]UserRole, error) {
	rows, err := q.db.QueryContext(ctx, getRoles, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []UserRole{}
	for rows.Next() {
		var i UserRole
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.UserID,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getRolesByName = `-- name: GetRolesByName :many
SELECT id, name, user_id, created_at, updated_at
FROM user_roles
WHERE name = $1
`

func (q *Queries) GetRolesByName(ctx context.Context, name string) ([]UserRole, error) {
	rows, err := q.db.QueryContext(ctx, getRolesByName, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []UserRole{}
	for rows.Next() {
		var i UserRole
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.UserID,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listRoles = `-- name: ListRoles :many
SELECT id, name, user_id, created_at, updated_at
FROM user_roles
LIMIT $1 OFFSET $2
`

type ListRolesParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListRoles(ctx context.Context, arg ListRolesParams) ([]UserRole, error) {
	rows, err := q.db.QueryContext(ctx, listRoles, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []UserRole{}
	for rows.Next() {
		var i UserRole
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.UserID,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

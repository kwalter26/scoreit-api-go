package db

import "context"

// CreateUserTxParams contains the input parameters of the CreateUser transaction
type CreateUserTxParams struct {
	CreateUserParams CreateUserParams
	CreateRoleParams CreateRoleParams
	AfterCreate      func(user User) error
}

// CreateUserTxResult is the result of the CreateUser transaction
type CreateUserTxResult struct {
	User      User
	UserRoles []UserRole
}

// CreateUserTx performs a money CreateUser from one account to the other.
// It creates a CreateUser record, add account entries, and update accounts' balance within a single database transaction.
// If any of the operations fail, it will rollback the transaction and return an error.
func (store *SQLStore) CreateUserTx(ctx context.Context, arg CreateUserTxParams) (CreateUserTxResult, error) {
	var result CreateUserTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.User, err = q.CreateUser(ctx, arg.CreateUserParams)
		if err != nil {
			return err
		}

		arg.CreateRoleParams.UserID = result.User.ID

		role, err := q.CreateRole(ctx, arg.CreateRoleParams)
		if err != nil {
			return err
		}

		result.UserRoles = []UserRole{role}

		if arg.AfterCreate != nil {
			return arg.AfterCreate(result.User)
		}

		return nil
	})

	return result, err
}

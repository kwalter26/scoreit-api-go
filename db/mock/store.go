// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/kwalter26/scoreit-api-go/db/sqlc (interfaces: Store)

// Package mockdb is a generated GoMock package.
package mockdb

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
	db "github.com/kwalter26/scoreit-api-go/db/sqlc"
)

// MockStore is a mock of Store interface.
type MockStore struct {
	ctrl     *gomock.Controller
	recorder *MockStoreMockRecorder
}

// MockStoreMockRecorder is the mock recorder for MockStore.
type MockStoreMockRecorder struct {
	mock *MockStore
}

// NewMockStore creates a new mock instance.
func NewMockStore(ctrl *gomock.Controller) *MockStore {
	mock := &MockStore{ctrl: ctrl}
	mock.recorder = &MockStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStore) EXPECT() *MockStoreMockRecorder {
	return m.recorder
}

// AddTeamMember mocks base method.
func (m *MockStore) AddTeamMember(arg0 context.Context, arg1 db.AddTeamMemberParams) (db.TeamMember, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddTeamMember", arg0, arg1)
	ret0, _ := ret[0].(db.TeamMember)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddTeamMember indicates an expected call of AddTeamMember.
func (mr *MockStoreMockRecorder) AddTeamMember(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddTeamMember", reflect.TypeOf((*MockStore)(nil).AddTeamMember), arg0, arg1)
}

// CreateGame mocks base method.
func (m *MockStore) CreateGame(arg0 context.Context, arg1 db.CreateGameParams) (db.Game, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateGame", arg0, arg1)
	ret0, _ := ret[0].(db.Game)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateGame indicates an expected call of CreateGame.
func (mr *MockStoreMockRecorder) CreateGame(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateGame", reflect.TypeOf((*MockStore)(nil).CreateGame), arg0, arg1)
}

// CreateGameParticipant mocks base method.
func (m *MockStore) CreateGameParticipant(arg0 context.Context, arg1 db.CreateGameParticipantParams) (db.GameParticipant, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateGameParticipant", arg0, arg1)
	ret0, _ := ret[0].(db.GameParticipant)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateGameParticipant indicates an expected call of CreateGameParticipant.
func (mr *MockStoreMockRecorder) CreateGameParticipant(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateGameParticipant", reflect.TypeOf((*MockStore)(nil).CreateGameParticipant), arg0, arg1)
}

// CreateRole mocks base method.
func (m *MockStore) CreateRole(arg0 context.Context, arg1 db.CreateRoleParams) (db.UserRole, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateRole", arg0, arg1)
	ret0, _ := ret[0].(db.UserRole)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateRole indicates an expected call of CreateRole.
func (mr *MockStoreMockRecorder) CreateRole(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateRole", reflect.TypeOf((*MockStore)(nil).CreateRole), arg0, arg1)
}

// CreateSession mocks base method.
func (m *MockStore) CreateSession(arg0 context.Context, arg1 db.CreateSessionParams) (db.Session, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateSession", arg0, arg1)
	ret0, _ := ret[0].(db.Session)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateSession indicates an expected call of CreateSession.
func (mr *MockStoreMockRecorder) CreateSession(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateSession", reflect.TypeOf((*MockStore)(nil).CreateSession), arg0, arg1)
}

// CreateTeam mocks base method.
func (m *MockStore) CreateTeam(arg0 context.Context, arg1 string) (db.Team, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateTeam", arg0, arg1)
	ret0, _ := ret[0].(db.Team)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateTeam indicates an expected call of CreateTeam.
func (mr *MockStoreMockRecorder) CreateTeam(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateTeam", reflect.TypeOf((*MockStore)(nil).CreateTeam), arg0, arg1)
}

// CreateUser mocks base method.
func (m *MockStore) CreateUser(arg0 context.Context, arg1 db.CreateUserParams) (db.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", arg0, arg1)
	ret0, _ := ret[0].(db.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateUser indicates an expected call of CreateUser.
func (mr *MockStoreMockRecorder) CreateUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockStore)(nil).CreateUser), arg0, arg1)
}

// CreateUserTx mocks base method.
func (m *MockStore) CreateUserTx(arg0 context.Context, arg1 db.CreateUserTxParams) (db.CreateUserTxResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUserTx", arg0, arg1)
	ret0, _ := ret[0].(db.CreateUserTxResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateUserTx indicates an expected call of CreateUserTx.
func (mr *MockStoreMockRecorder) CreateUserTx(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUserTx", reflect.TypeOf((*MockStore)(nil).CreateUserTx), arg0, arg1)
}

// DeleteGameParticipant mocks base method.
func (m *MockStore) DeleteGameParticipant(arg0 context.Context, arg1 uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteGameParticipant", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteGameParticipant indicates an expected call of DeleteGameParticipant.
func (mr *MockStoreMockRecorder) DeleteGameParticipant(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteGameParticipant", reflect.TypeOf((*MockStore)(nil).DeleteGameParticipant), arg0, arg1)
}

// DeleteRole mocks base method.
func (m *MockStore) DeleteRole(arg0 context.Context, arg1 uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteRole", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteRole indicates an expected call of DeleteRole.
func (mr *MockStoreMockRecorder) DeleteRole(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteRole", reflect.TypeOf((*MockStore)(nil).DeleteRole), arg0, arg1)
}

// DeleteTeam mocks base method.
func (m *MockStore) DeleteTeam(arg0 context.Context, arg1 uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteTeam", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteTeam indicates an expected call of DeleteTeam.
func (mr *MockStoreMockRecorder) DeleteTeam(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteTeam", reflect.TypeOf((*MockStore)(nil).DeleteTeam), arg0, arg1)
}

// DeleteUser mocks base method.
func (m *MockStore) DeleteUser(arg0 context.Context, arg1 uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteUser", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteUser indicates an expected call of DeleteUser.
func (mr *MockStoreMockRecorder) DeleteUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteUser", reflect.TypeOf((*MockStore)(nil).DeleteUser), arg0, arg1)
}

// GetGame mocks base method.
func (m *MockStore) GetGame(arg0 context.Context, arg1 uuid.UUID) (db.Game, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetGame", arg0, arg1)
	ret0, _ := ret[0].(db.Game)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetGame indicates an expected call of GetGame.
func (mr *MockStoreMockRecorder) GetGame(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGame", reflect.TypeOf((*MockStore)(nil).GetGame), arg0, arg1)
}

// GetGameParticipant mocks base method.
func (m *MockStore) GetGameParticipant(arg0 context.Context, arg1 uuid.UUID) (db.GameParticipant, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetGameParticipant", arg0, arg1)
	ret0, _ := ret[0].(db.GameParticipant)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetGameParticipant indicates an expected call of GetGameParticipant.
func (mr *MockStoreMockRecorder) GetGameParticipant(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGameParticipant", reflect.TypeOf((*MockStore)(nil).GetGameParticipant), arg0, arg1)
}

// GetRole mocks base method.
func (m *MockStore) GetRole(arg0 context.Context, arg1 uuid.UUID) (db.UserRole, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRole", arg0, arg1)
	ret0, _ := ret[0].(db.UserRole)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRole indicates an expected call of GetRole.
func (mr *MockStoreMockRecorder) GetRole(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRole", reflect.TypeOf((*MockStore)(nil).GetRole), arg0, arg1)
}

// GetRoles mocks base method.
func (m *MockStore) GetRoles(arg0 context.Context, arg1 uuid.UUID) ([]db.UserRole, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRoles", arg0, arg1)
	ret0, _ := ret[0].([]db.UserRole)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRoles indicates an expected call of GetRoles.
func (mr *MockStoreMockRecorder) GetRoles(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRoles", reflect.TypeOf((*MockStore)(nil).GetRoles), arg0, arg1)
}

// GetRolesByName mocks base method.
func (m *MockStore) GetRolesByName(arg0 context.Context, arg1 string) ([]db.UserRole, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRolesByName", arg0, arg1)
	ret0, _ := ret[0].([]db.UserRole)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRolesByName indicates an expected call of GetRolesByName.
func (mr *MockStoreMockRecorder) GetRolesByName(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRolesByName", reflect.TypeOf((*MockStore)(nil).GetRolesByName), arg0, arg1)
}

// GetSession mocks base method.
func (m *MockStore) GetSession(arg0 context.Context, arg1 uuid.UUID) (db.Session, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSession", arg0, arg1)
	ret0, _ := ret[0].(db.Session)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSession indicates an expected call of GetSession.
func (mr *MockStoreMockRecorder) GetSession(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSession", reflect.TypeOf((*MockStore)(nil).GetSession), arg0, arg1)
}

// GetTeam mocks base method.
func (m *MockStore) GetTeam(arg0 context.Context, arg1 uuid.UUID) (db.Team, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTeam", arg0, arg1)
	ret0, _ := ret[0].(db.Team)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTeam indicates an expected call of GetTeam.
func (mr *MockStoreMockRecorder) GetTeam(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTeam", reflect.TypeOf((*MockStore)(nil).GetTeam), arg0, arg1)
}

// GetUser mocks base method.
func (m *MockStore) GetUser(arg0 context.Context, arg1 uuid.UUID) (db.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUser", arg0, arg1)
	ret0, _ := ret[0].(db.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUser indicates an expected call of GetUser.
func (mr *MockStoreMockRecorder) GetUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUser", reflect.TypeOf((*MockStore)(nil).GetUser), arg0, arg1)
}

// GetUserByUsername mocks base method.
func (m *MockStore) GetUserByUsername(arg0 context.Context, arg1 string) (db.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByUsername", arg0, arg1)
	ret0, _ := ret[0].(db.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByUsername indicates an expected call of GetUserByUsername.
func (mr *MockStoreMockRecorder) GetUserByUsername(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByUsername", reflect.TypeOf((*MockStore)(nil).GetUserByUsername), arg0, arg1)
}

// ListGameParticipants mocks base method.
func (m *MockStore) ListGameParticipants(arg0 context.Context, arg1 db.ListGameParticipantsParams) ([]db.GameParticipant, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListGameParticipants", arg0, arg1)
	ret0, _ := ret[0].([]db.GameParticipant)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListGameParticipants indicates an expected call of ListGameParticipants.
func (mr *MockStoreMockRecorder) ListGameParticipants(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListGameParticipants", reflect.TypeOf((*MockStore)(nil).ListGameParticipants), arg0, arg1)
}

// ListGameParticipantsForPlayer mocks base method.
func (m *MockStore) ListGameParticipantsForPlayer(arg0 context.Context, arg1 db.ListGameParticipantsForPlayerParams) ([]db.GameParticipant, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListGameParticipantsForPlayer", arg0, arg1)
	ret0, _ := ret[0].([]db.GameParticipant)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListGameParticipantsForPlayer indicates an expected call of ListGameParticipantsForPlayer.
func (mr *MockStoreMockRecorder) ListGameParticipantsForPlayer(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListGameParticipantsForPlayer", reflect.TypeOf((*MockStore)(nil).ListGameParticipantsForPlayer), arg0, arg1)
}

// ListGames mocks base method.
func (m *MockStore) ListGames(arg0 context.Context, arg1 db.ListGamesParams) ([]db.Game, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListGames", arg0, arg1)
	ret0, _ := ret[0].([]db.Game)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListGames indicates an expected call of ListGames.
func (mr *MockStoreMockRecorder) ListGames(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListGames", reflect.TypeOf((*MockStore)(nil).ListGames), arg0, arg1)
}

// ListRoles mocks base method.
func (m *MockStore) ListRoles(arg0 context.Context, arg1 db.ListRolesParams) ([]db.UserRole, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListRoles", arg0, arg1)
	ret0, _ := ret[0].([]db.UserRole)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListRoles indicates an expected call of ListRoles.
func (mr *MockStoreMockRecorder) ListRoles(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListRoles", reflect.TypeOf((*MockStore)(nil).ListRoles), arg0, arg1)
}

// ListTeamMembers mocks base method.
func (m *MockStore) ListTeamMembers(arg0 context.Context, arg1 db.ListTeamMembersParams) ([]db.ListTeamMembersRow, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListTeamMembers", arg0, arg1)
	ret0, _ := ret[0].([]db.ListTeamMembersRow)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListTeamMembers indicates an expected call of ListTeamMembers.
func (mr *MockStoreMockRecorder) ListTeamMembers(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListTeamMembers", reflect.TypeOf((*MockStore)(nil).ListTeamMembers), arg0, arg1)
}

// ListTeams mocks base method.
func (m *MockStore) ListTeams(arg0 context.Context, arg1 db.ListTeamsParams) ([]db.Team, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListTeams", arg0, arg1)
	ret0, _ := ret[0].([]db.Team)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListTeams indicates an expected call of ListTeams.
func (mr *MockStoreMockRecorder) ListTeams(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListTeams", reflect.TypeOf((*MockStore)(nil).ListTeams), arg0, arg1)
}

// ListTeamsOfUser mocks base method.
func (m *MockStore) ListTeamsOfUser(arg0 context.Context, arg1 db.ListTeamsOfUserParams) ([]db.ListTeamsOfUserRow, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListTeamsOfUser", arg0, arg1)
	ret0, _ := ret[0].([]db.ListTeamsOfUserRow)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListTeamsOfUser indicates an expected call of ListTeamsOfUser.
func (mr *MockStoreMockRecorder) ListTeamsOfUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListTeamsOfUser", reflect.TypeOf((*MockStore)(nil).ListTeamsOfUser), arg0, arg1)
}

// ListUsers mocks base method.
func (m *MockStore) ListUsers(arg0 context.Context, arg1 db.ListUsersParams) ([]db.ListUsersRow, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListUsers", arg0, arg1)
	ret0, _ := ret[0].([]db.ListUsersRow)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListUsers indicates an expected call of ListUsers.
func (mr *MockStoreMockRecorder) ListUsers(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListUsers", reflect.TypeOf((*MockStore)(nil).ListUsers), arg0, arg1)
}

// UpdateGame mocks base method.
func (m *MockStore) UpdateGame(arg0 context.Context, arg1 db.UpdateGameParams) (db.Game, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateGame", arg0, arg1)
	ret0, _ := ret[0].(db.Game)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateGame indicates an expected call of UpdateGame.
func (mr *MockStoreMockRecorder) UpdateGame(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateGame", reflect.TypeOf((*MockStore)(nil).UpdateGame), arg0, arg1)
}

// UpdateGameParticipant mocks base method.
func (m *MockStore) UpdateGameParticipant(arg0 context.Context, arg1 db.UpdateGameParticipantParams) (db.GameParticipant, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateGameParticipant", arg0, arg1)
	ret0, _ := ret[0].(db.GameParticipant)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateGameParticipant indicates an expected call of UpdateGameParticipant.
func (mr *MockStoreMockRecorder) UpdateGameParticipant(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateGameParticipant", reflect.TypeOf((*MockStore)(nil).UpdateGameParticipant), arg0, arg1)
}

// UpdateSession mocks base method.
func (m *MockStore) UpdateSession(arg0 context.Context, arg1 db.UpdateSessionParams) (db.Session, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateSession", arg0, arg1)
	ret0, _ := ret[0].(db.Session)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateSession indicates an expected call of UpdateSession.
func (mr *MockStoreMockRecorder) UpdateSession(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateSession", reflect.TypeOf((*MockStore)(nil).UpdateSession), arg0, arg1)
}

// UpdateTeam mocks base method.
func (m *MockStore) UpdateTeam(arg0 context.Context, arg1 db.UpdateTeamParams) (db.Team, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateTeam", arg0, arg1)
	ret0, _ := ret[0].(db.Team)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateTeam indicates an expected call of UpdateTeam.
func (mr *MockStoreMockRecorder) UpdateTeam(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateTeam", reflect.TypeOf((*MockStore)(nil).UpdateTeam), arg0, arg1)
}

// UpdateUser mocks base method.
func (m *MockStore) UpdateUser(arg0 context.Context, arg1 db.UpdateUserParams) (db.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUser", arg0, arg1)
	ret0, _ := ret[0].(db.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateUser indicates an expected call of UpdateUser.
func (mr *MockStoreMockRecorder) UpdateUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUser", reflect.TypeOf((*MockStore)(nil).UpdateUser), arg0, arg1)
}

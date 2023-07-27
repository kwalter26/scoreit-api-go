package security

type Role string

var UserRole Role = "user"
var AdminRole Role = "admin"
var UserRoles = []Role{UserRole}
var AdminRoles = []Role{AdminRole}

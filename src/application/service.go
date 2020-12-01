package application

/*
	User Service
*/
const (
	Auth_NoVerify = 0
	Auth_BasicVerify = 1
	Auth_AdvancedVerify = 2
)

// Client Service name
// 客户端服务名
const (
	Service_AttachmentService = "AttachmentService"
)

// Relationship name
// 关系名
const (
	Relationship_User_Avatar = "user-avatar"
)

// Cache variable
// 缓存变量
var (
	Cache_MenuGetAll = "MenuGetAll"
	Cache_UserGetCurrent = "UserGetCurrent"
	Cache_PermissionGetAll = "PermissionGetAll"
	Cache_UserGetAll = "UserGetAll"
	Cache_RoleGetAll = "RoleGetAll"
)

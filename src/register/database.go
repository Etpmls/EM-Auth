package register

import (
	"github.com/Etpmls/EM-Auth/src/application"
	"github.com/Etpmls/EM-Auth/src/application/database"
	em "github.com/Etpmls/Etpmls-Micro"

	"gorm.io/gorm"
)

func InsertBasicDataToDatabase()  {
	// Create Role
	role := database.Role{
		Name:        "Administrator",
		Remark: "System Administrator",
	}
	if err := em.DB.Debug().Create(&role).Error; err != nil {
		em.LogError.Output(em.MessageWithLineNum(err.Error()))
	}


	// Create User
	user := database.User{
		Username: "admin",
		Password: "$2a$10$yNoJrsN7mrtHzUyvm6s8KOwHrnkkGmqcRJvcieQKItIfQNwyzqfMy",
		Roles: []database.Role{
			{
				Model:       gorm.Model{ID:1},
			},
		},
	}
	if err := em.DB.Debug().Set("gorm:association_autoupdate", false).Set("gorm:association_autocreate", false).Create(&user).Error; err != nil {
		em.LogError.Output(em.MessageWithLineNum(err.Error()))
	}

	// Create Permission
	permission := []database.Permission{
		{
			Name: "User/Login",
			Method: "POST",
			Path: "/api/auth/*/user/login",
			Auth: application.Auth_NoVerify,
		},
		{
			Name: "User/Logout",
			Method: "POST",
			Path: "/api/auth/*/user/logout",
			Auth: application.Auth_NoVerify,
		},
		{
			Name: "User/UpdateInformation",
			Method: "PUT",
			Path: "/api/auth/*/user/updateInformation",
			Auth: application.Auth_BasicVerify,
		},
		{
			Name: "User/Get Current",
			Method: "GET",
			Path: "/api/auth/*/user/getCurrent",
			Auth: application.Auth_BasicVerify,
		},
		{
			Name: "User/View",
			Method: "GET",
			Path: "/api/auth/*/user/getAll",
			Auth: application.Auth_AdvancedVerify,
			Remark: "View user list",
			Roles: []database.Role{
				{
					Model:       gorm.Model{ID:1},
				},
			},
		},
		{
			Name: "User/Create",
			Method: "POST",
			Path: "/api/auth/*/user/create",
			Auth: application.Auth_AdvancedVerify,
			Roles: []database.Role{
				{
					Model:       gorm.Model{ID:1},
				},
			},
		},
		{
			Name: "User/Edit",
			Method: "PUT",
			Path: "/api/auth/*/user/edit",
			Auth: application.Auth_AdvancedVerify,
			Roles: []database.Role{
				{
					Model:       gorm.Model{ID:1},
				},
			},
		},
		{
			Name: "User/Delete",
			Method: "DELETE",
			Path: "/api/auth/*/user/delete",
			Auth: application.Auth_AdvancedVerify,
			Roles: []database.Role{
				{
					Model:       gorm.Model{ID:1},
				},
			},
		},
		{
			Name: "Role/View",
			Method: "GET",
			Path: "/api/auth/*/role/getAll",
			Remark: "View role list",
			Auth: application.Auth_AdvancedVerify,
			Roles: []database.Role{
				{
					Model:       gorm.Model{ID:1},
				},
			},
		},
		{
			Name: "Role/Create",
			Method: "POST",
			Path: "/api/auth/*/role/create",
			Auth: application.Auth_AdvancedVerify,
			Roles: []database.Role{
				{
					Model:       gorm.Model{ID:1},
				},
			},
		},
		{
			Name: "Role/Edit",
			Method: "PUT",
			Path: "/api/auth/*/role/edit",
			Auth: application.Auth_AdvancedVerify,
			Roles: []database.Role{
				{
					Model:       gorm.Model{ID:1},
				},
			},
		},
		{
			Name: "Role/Delete",
			Method: "DELETE",
			Path: "/api/auth/*/role/delete",
			Auth: application.Auth_AdvancedVerify,
			Roles: []database.Role{
				{
					Model:       gorm.Model{ID:1},
				},
			},
		},
		{
			Name: "Permission/GetAdvancedVerify",
			Method: "GET",
			Path: "/api/auth/*/permission/getAdvancedVerify",
			Auth: application.Auth_AdvancedVerify,
			Remark: "View permission list",
			Roles: []database.Role{
				{
					Model:       gorm.Model{ID:1},
				},
			},
		},
		{
			Name: "Permission/View",
			Method: "GET",
			Path: "/api/auth/*/permission/getAll",
			Auth: application.Auth_AdvancedVerify,
			Remark: "View permission list",
			Roles: []database.Role{
				{
					Model:       gorm.Model{ID:1},
				},
			},
		},
		{
			Name: "Permission/Create",
			Method: "POST",
			Path: "/api/auth/*/permission/create",
			Auth: application.Auth_AdvancedVerify,
			Roles: []database.Role{
				{
					Model:       gorm.Model{ID:1},
				},
			},
		},
		{
			Name: "Permission/Edit",
			Method: "PUT",
			Path: "/api/auth/*/permission/edit",
			Auth: application.Auth_AdvancedVerify,
			Roles: []database.Role{
				{
					Model:       gorm.Model{ID:1},
				},
			},
		},
		{
			Name: "Permission/Delete",
			Method: "DELETE",
			Path: "/api/auth/*/permission/delete",
			Auth: application.Auth_AdvancedVerify,
			Roles: []database.Role{
				{
					Model:       gorm.Model{ID:1},
				},
			},
		},
		{
			Name: "Menu/Create/Edit",
			Method: "POST",
			Path: "/api/auth/*/menu/create",
			Auth: application.Auth_AdvancedVerify,
			Roles: []database.Role{
				{
					Model:       gorm.Model{ID:1},
				},
			},
		},
		{
			Name: "Menu/Get All",
			Method: "GET",
			Path: "/api/auth/*/menu/getAll",
			Auth: application.Auth_BasicVerify,
		},
		{
			Name: "Setting/Cache Clear",
			Method: "GET",
			Path: "/api/auth/*/setting/cacheClear",
			Auth: application.Auth_AdvancedVerify,
			Roles: []database.Role{
				{
					Model:       gorm.Model{ID:1},
				},
			},
		},
		{
			Name: "Setting/Disk Clean Up",
			Method: "GET",
			Path: "/api/auth/*/setting/diskCleanUp",
			Auth: application.Auth_AdvancedVerify,
			Roles: []database.Role{
				{
					Model:       gorm.Model{ID:1},
				},
			},
		},
		{
			Name: "Attachment/Get One",
			Method: "GET",
			Path: "/api/attachment/*/attachment/getOne",
			Auth: application.Auth_NoVerify,
		},
		{
			Name: "Attachment/Create",
			Method: "GET",
			Path: "/api/attachment/*/attachment/diskCleanUp",
			Auth: application.Auth_AdvancedVerify,
			Roles: []database.Role{
				{
					Model:       gorm.Model{ID:1},
				},
			},
		},

	}
	if err := em.DB.Debug().Set("gorm:association_autoupdate", false).Set("gorm:association_autocreate", false).Create(&permission).Error; err != nil {
		em.LogError.Output(em.MessageWithLineNum(err.Error()))
	}
}

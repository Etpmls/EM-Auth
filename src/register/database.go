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
			Path: "User/Login",
			Auth: application.Auth_NoVerify,
		},
		{
			Name: "User/Logout",
			Path: "User/Logout",
			Auth: application.Auth_NoVerify,
		},
		{
			Name: "User/UpdateInformation",
			Path: "User/UpdateInformation",
			Auth: application.Auth_BasicVerify,
		},
		{
			Name: "User/Get Current",
			Path: "User/GetCurrent",
			Auth: application.Auth_BasicVerify,
		},
		{
			Name: "User/View",
			Path: "User/GetAll",
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
			Path: "User/Create",
			Auth: application.Auth_AdvancedVerify,
			Roles: []database.Role{
				{
					Model:       gorm.Model{ID:1},
				},
			},
		},
		{
			Name: "User/Edit",
			Path: "User/Edit",
			Auth: application.Auth_AdvancedVerify,
			Roles: []database.Role{
				{
					Model:       gorm.Model{ID:1},
				},
			},
		},
		{
			Name: "User/Delete",
			Path: "User/Delete",
			Auth: application.Auth_AdvancedVerify,
			Roles: []database.Role{
				{
					Model:       gorm.Model{ID:1},
				},
			},
		},
		{
			Name: "Role/View",
			Path: "Role/GetAll",
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
			Path: "Role/Create",
			Auth: application.Auth_AdvancedVerify,
			Roles: []database.Role{
				{
					Model:       gorm.Model{ID:1},
				},
			},
		},
		{
			Name: "Role/Edit",
			Path: "Role/Edit",
			Auth: application.Auth_AdvancedVerify,
			Roles: []database.Role{
				{
					Model:       gorm.Model{ID:1},
				},
			},
		},
		{
			Name: "Role/Delete",
			Path: "Role/Delete",
			Auth: application.Auth_AdvancedVerify,
			Roles: []database.Role{
				{
					Model:       gorm.Model{ID:1},
				},
			},
		},
		{
			Name: "Permission/GetAdvancedVerify",
			Path: "Permission/GetAdvancedVerify",
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
			Path: "Permission/GetAll",
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
			Path: "Permission/Create",
			Auth: application.Auth_AdvancedVerify,
			Roles: []database.Role{
				{
					Model:       gorm.Model{ID:1},
				},
			},
		},
		{
			Name: "Permission/Edit",
			Path: "Permission/Edit",
			Auth: application.Auth_AdvancedVerify,
			Roles: []database.Role{
				{
					Model:       gorm.Model{ID:1},
				},
			},
		},
		{
			Name: "Permission/Delete",
			Path: "Permission/Delete",
			Auth: application.Auth_AdvancedVerify,
			Roles: []database.Role{
				{
					Model:       gorm.Model{ID:1},
				},
			},
		},
		{
			Name: "Menu/Create/Edit",
			Path: "Menu/Create",
			Auth: application.Auth_AdvancedVerify,
			Roles: []database.Role{
				{
					Model:       gorm.Model{ID:1},
				},
			},
		},
		{
			Name: "Menu/Get All",
			Path: "Menu/GetAll",
			Auth: application.Auth_BasicVerify,
		},
		{
			Name: "Setting/Cache Clear",
			Path: "Setting/CacheClear",
			Auth: application.Auth_AdvancedVerify,
			Roles: []database.Role{
				{
					Model:       gorm.Model{ID:1},
				},
			},
		},
		{
			Name: "Setting/Disk Clean Up",
			Path: "Setting/DiskCleanUp",
			Auth: application.Auth_AdvancedVerify,
			Roles: []database.Role{
				{
					Model:       gorm.Model{ID:1},
				},
			},
		},
		{
			Name: "Setting/Disk Clean Up",
			Path: "Setting/DiskCleanUp",
			Auth: application.Auth_AdvancedVerify,
			Roles: []database.Role{
				{
					Model:       gorm.Model{ID:1},
				},
			},
		},
		{
			Name: "Attachment/Get One",
			Path: "Attachment/GetOne",
			Auth: application.Auth_NoVerify,
		},
		{
			Name: "Attachment/Create",
			Path: "Attachment/Create",
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

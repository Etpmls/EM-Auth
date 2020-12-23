package main

import (
	"github.com/Etpmls/EM-Auth/src/application"
	"github.com/Etpmls/EM-Auth/src/application/database"
	"github.com/Etpmls/EM-Auth/src/register"
	"github.com/Etpmls/Etpmls-Micro"
)


func main() {
	var reg = em.Register{
		Version_Service: 		map[string]string{"EM-Auth Version": application.Version_Service},
		GrpcServiceFunc:    	register.RegisterRpcService,
		GrpcMiddleware:     	register.RegisterGrpcMiddleware,
		HttpServiceFunc:    	register.RegisterHttpService,
		RouteFunc:          	register.RegisterRoute,
		DatabaseMigrate:		[]interface{}{
			&database.User{},
			&database.Role{},
			&database.Permission{},
		},
		InsertDatabaseInitialData: []func(){register.InsertBasicDataToDatabase},
		CustomConfiguration: struct {
			Path       string
			DebugPath  string
			StructAddr interface{}
		}{Path: "storage/config/auth.yaml", DebugPath: "storage/config/auth_debug.yaml", StructAddr: &application.ServiceConfig},
	}
	reg.Init()
	reg.Run()
}

package main

import (
	"github.com/Etpmls/EM-Auth/src/application/database"
	"github.com/Etpmls/EM-Auth/src/register"
	"github.com/Etpmls/Etpmls-Micro"
	em_library "github.com/Etpmls/Etpmls-Micro/library"
)


func main() {
	var reg = em.Register{
		Config:             	em_library.Config,
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
	}

	reg.Run()
}

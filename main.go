package main

import (
	"github.com/Etpmls/EM-Auth/src/application"
	"github.com/Etpmls/EM-Auth/src/application/database"
	"github.com/Etpmls/EM-Auth/src/register"
	em "github.com/Etpmls/Etpmls-Micro/v2"
)


func main() {
	var reg = em.Register{
		AppVersion: 		map[string]string{"EM-Auth Version": application.Version_Service},
		AppEnabledFeatureName:		[]string{em.EnableCaptcha, em.EnableCircuitBreaker, em.EnableDatabase, em.EnableI18n, em.EnableServiceDiscovery, em.EnableValidator},
		RpcServiceFunc:    	register.RegisterRpcService,
		RpcMiddleware:     	register.RegisterGrpcMiddleware,
		HttpServiceFunc:    	register.RegisterHttpService,
		HttpRouteFunc:          	register.RegisterRoute,
		DatabaseMigrate:		[]interface{}{
			&database.User{},
			&database.Role{},
			&database.Permission{},
		},
		DatabaseInsertInitialData: []func(){register.InsertBasicDataToDatabase},
	}
	reg.Init()
	reg.Run()
}

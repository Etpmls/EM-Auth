package register

import (
	"github.com/Etpmls/EM-Auth/src/application/service"
	em "github.com/Etpmls/Etpmls-Micro/v2"
	"github.com/Etpmls/Etpmls-Micro/v2/define"
	em_library "github.com/Etpmls/Etpmls-Micro/v2/library"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"net/http"
)

// Register Route
func RegisterRoute(mux *runtime.ServeMux)  {
	e, _ := em.Kv.ReadKey(define.MakeServiceConfField(em_library.Config.Service.RpcId, define.KvServiceCheckUrl))

	mux.HandlePath("GET", e, func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		w.Write([]byte("hello"))
	})
	mux.HandlePath("GET", "/api/checkAuth", service.ServiceAuth{}.Check)

}
package engine

import (
	"deny"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
)

func (d *DbEngine) TokenAuth(next http.Handler) http.Handler {
	//权限验证
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		poid, err := primitive.ObjectIDFromHex(auth)
		if auth != "" && err == nil && deny.InjectionPass([]byte(auth)) {
			log.Println(poid)
		}
		r.Header.Set("role", "test role")
		next.ServeHTTP(w, r)
		//// Request Basic Authentication otherwise
		//w.Header().Set("WWW-Authenticate", "Basic realm=Restricted")
		//http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
	})
}

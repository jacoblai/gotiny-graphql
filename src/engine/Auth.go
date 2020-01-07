package engine

import (
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"strings"
)

func (d *DbEngine) TokenAuth(next http.Handler) http.Handler {
	//权限验证
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if strings.HasPrefix(auth, "Bearer ") {
			auth = strings.Replace(auth, "Bearer ", "", -1)
		}
		if auth != "" {
			rtoken, err := jwt.ParseWithClaims(auth, &CoolpyClaims{}, func(token *jwt.Token) (interface{}, error) {
				return []byte("Coolpy7yeah"), nil
			})
			if err == nil {
				if tk, ok := rtoken.Claims.(*CoolpyClaims); ok && rtoken.Valid && tk.Issuer == "coolpy7_api" {
					r.Header.Set("user_id", tk.UserId)
					r.Header.Set("user_name", tk.Uid)
					r.Header.Set("rule", tk.Rule)
					//r.Header.Set("client", tk.Client)
					//r.Header.Set("access", strconv.Itoa(int(tk.Access)))
					next.ServeHTTP(w, r)
					return
				}
			}
		}

		// Request Basic Authentication otherwise
		w.Header().Set("WWW-Authenticate", "Bearer realm=Restricted")
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
	})
}

type CoolpyClaims struct {
	UserId string `json:"user_id"`
	Uid    string `json:"uid"`
	Rule   string `json:"rule"`
	Client string `json:"client" bson:"client,omitempty"`
	Access byte   `json:"access" bson:"access,omitempty"`
	jwt.StandardClaims
}

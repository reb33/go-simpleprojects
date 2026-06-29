package middleware

import (
	"context"
	"demo-store/configs"
	"demo-store/pkg/jwt"
	"net/http"
	"strings"
)

type key string

const (
	ContextPhoneKey key = "ContextPhoneKey"
)

func IsAuthed(next http.Handler, config *configs.Configs) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		if !strings.HasPrefix(authHeader, "Bearer") { // проверка наличия хедера
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		token := strings.TrimPrefix(authHeader, "Bearer ")           // проверка наличия токена
		isValid, data := jwt.NewJWT(config.Auth.Secret).Parse(token) // проверка валидности токена и получение данных из токена, если токен валидный. Если токен не валидный, то возвращается ошибка 401
		if !isValid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), ContextPhoneKey, data.Phone) // добавление данных в контекст
		req := r.WithContext(ctx)
		next.ServeHTTP(w, req)
	})
}

package auth_test

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/ShukinDmitriy/GophKeeper/internal/server/config"
	"github.com/ShukinDmitriy/GophKeeper/internal/server/models/responses"
	mockRepositories "github.com/ShukinDmitriy/GophKeeper/mocks/internal_/server/repositories"

	"github.com/ShukinDmitriy/GophKeeper/internal/server/auth"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestGetAccessTokenCookieName(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "positive test #1",
			want: "access-token",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := auth.GetAccessTokenCookieName(); got != tt.want {
				t.Errorf("GetAccessTokenCookieName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetJWTSecret(t *testing.T) {
	tests := []struct {
		name string
		want bool
	}{
		{
			name: "positive test #1",
			want: true,
		},
	}

	_, err := config.NewConfig()
	if err != nil {
		t.Fatal(err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := auth.GetJWTSecret(); got == "" {
				t.Errorf("GetJWTSecret() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetSigningMethod(t *testing.T) {
	tests := []struct {
		name string
		want *jwt.SigningMethodHMAC
	}{
		{
			name: "positive test #1",
			want: jwt.SigningMethodHS256,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := auth.GetSigningMethod(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetSigningMethod() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGenerateTokensAndSetCookies(t *testing.T) {
	type args struct {
		user *responses.UserInfo
	}
	type want struct {
		error bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "positive test #1",
			args: args{
				user: &responses.UserInfo{
					ID: 123,
				},
			},
			want: want{
				error: false,
			},
		},
	}

	_, err := config.NewConfig()
	if err != nil {
		t.Fatal(err)
	}

	userRepo := new(mockRepositories.UserRepositoryInterface)
	authUser := auth.NewAuthUser(userRepo)
	authService := auth.NewAuthService(*authUser)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			var err error
			// Устанавливаем куку
			assert.Nil(t, func(c echo.Context) error {
				err = authService.GenerateTokensAndSetCookies(c, tt.args.user)

				return c.NoContent(http.StatusOK)
			}(c))
			assert.Nil(t, err)

			res := rec.Result()
			defer res.Body.Close()

			// Проверяем, что ответ сервера OK
			assert.Equal(t, http.StatusOK, res.StatusCode)
			headers := res.Header.Values("Set-Cookie")
			hasAccessCookie := false
			for _, h := range headers {
				tmp := strings.Split(h, "=")
				if tmp[0] == auth.GetAccessTokenCookieName() {
					hasAccessCookie = true
					break
				}
			}
			assert.True(t, hasAccessCookie)
		})
	}
}

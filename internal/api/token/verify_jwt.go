package token

import (
	"os"
	"time"

	"github.com/form3tech-oss/jwt-go"

	"github.com/ashah360/nyte-auth/internal/api/cerror"
)

type AccessTokenPayload struct {
	ID    string `json:"id"`
	Token string `json:"token"`
}

func ValidateJWT(token string) (*AccessTokenPayload, error) {
	t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, cerror.ErrMalformedToken
		}

		return []byte(os.Getenv("NYTE_SECRET")), nil
	})
	if err != nil {
		return nil, cerror.ErrMalformedToken
	}

	if !t.Valid {
		return nil, cerror.ErrInvalidToken
	}

	claims := t.Claims.(jwt.MapClaims)

	exp, ok := claims["exp"]
	if !ok {
		return nil, cerror.ErrMalformedToken
	}

	if time.Since(time.Unix(int64(exp.(float64)), 0)) > 0 {
		return nil, cerror.ErrInvalidToken
	}

	uid, ok := claims["id"]
	if !ok {
		return nil, cerror.ErrMalformedToken
	}

	ot, ok := claims["token"]
	if !ok {
		return nil, cerror.ErrMalformedToken
	}

	return &AccessTokenPayload{
		ID:    uid.(string),
		Token: ot.(string),
	}, nil
}

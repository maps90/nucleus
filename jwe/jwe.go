package jwe

import (
	"time"

	"github.com/maps90/nucleus/util"
	"gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"
)

const (
	// default duration 5 minutes
	defaultDuration = 300
	// default JWE subject
	defaultSubject = "Jwt Encrypted AUTH"
	// default JWE Issuer
	defaultIssuer = "XAUTH"
	// default authorization bearer name
	defaultBearer = "Bearer"
)

type Config struct {
	signature    string
	data         *StandardClaims
	refreshToken string
	expired      *time.Time
	bearer       string
}

func New(key string) *Config {
	return &Config{
		signature: key,
	}
}

func (s *Config) SetData(claims *StandardClaims) *Config {
	expiresAt := TimeFunc().Add(time.Duration(defaultDuration) * time.Second)
	refreshToken := util.RandomString(6)

	s.refreshToken = refreshToken
	s.expired = &expiresAt

	claims.RefreshToken = s.RefreshToken()
	claims.ExpiresOn = s.ExpiresAt().Format(time.RFC3339)
	claims.jwtID = util.RandomNumber(10)
	claims.Subject = defaultSubject
	claims.Issuer = defaultIssuer
	claims.IssuedAt = TimeFunc().Unix()
	claims.NotBefore = TimeFunc().Unix()

	// added expired tolerance (5 second)
	claims.ExpiresAt = s.ExpiresAt().Add(time.Duration(5) * time.Second).Unix()

	s.data = claims

	return s
}

func (s *Config) SetBearer(bearer string) {
	s.bearer = bearer
}

func (s *Config) GetBearer() string {
	if s.bearer == "" {
		return defaultBearer
	}
	return s.bearer
}

func (s *Config) ExpiresAt() *time.Time {
	return s.expired
}

func (s *Config) RefreshToken() string {
	return s.refreshToken
}

func (s *Config) Compact() (token string, err error) {
	sig, err := jose.NewEncrypter(
		jose.A128GCM,
		jose.Recipient{Algorithm: jose.DIRECT, Key: []byte(s.signature)},
		(&jose.EncrypterOptions{}).WithType("JWT"),
	)
	if err != nil {
		return token, err
	}

	raw, err := jwt.Encrypted(sig).Claims(s.data).CompactSerialize()
	if err != nil {
		return token, err
	}

	return raw, nil
}

func (c *Config) Decrypt(accessToken string) (result *StandardClaims, err error) {
	tok, err := jwt.ParseEncrypted(accessToken)
	if err != nil {
		return result, err
	}

	cl := StandardClaims{}
	if err = tok.Claims([]byte(c.signature), &cl); err != nil {
		return result, err
	}

	return &cl, nil
}

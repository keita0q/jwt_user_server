package jwtAuth

import (
	"github.com/dgrijalva/jwt-go"
	"crypto/rsa"
	"io/ioutil"
	"fmt"
	"github.com/keita0q/user_server/database/sequreDB"
	"github.com/keita0q/user_server/auth"
)

type JWTAuth struct {
	database       sequreDB.SequreDB
	publicKeyPath  string
	privateKeyPath string
}
type Config struct {
	DB             sequreDB.SequreDB
	PublicKeyPath  string
	PrivateKeyPath string
}

func New(aConfig *Config) *JWTAuth {
	return &JWTAuth{
		database: aConfig.DB,
		publicKeyPath: aConfig.PublicKeyPath,
		privateKeyPath: aConfig.PrivateKeyPath,
	}
}

type Claim struct {
	ID string `json:"id"`
	jwt.StandardClaims
}

func (aClaim *Claim)GetUserID() string {
	return aClaim.ID
}

func (aAuth *JWTAuth) CreateToken(aID string, aPassword string) (string, error) {
	if (!aAuth.database.Exist(aID, aPassword)) {
		return "", &auth.NotFoundError{Message: aID + "は存在しません"}
	}

	tToken := jwt.NewWithClaims(jwt.SigningMethodRS256, &Claim{ID: aID})
	tKey, tError := lookupPrivateKey(aAuth.privateKeyPath)
	if tError != nil {
		return "", tError
	}

	tTokenString, tError := tToken.SignedString(tKey)
	if tError != nil {
		fmt.Println(tError)
		return "", tError
	}

	return tTokenString, nil
}

func (aAuth *JWTAuth)Authenticate(aToken string) (auth.Claim, bool, error) {
	tClaims := &Claim{}
	tToken, tError := jwt.ParseWithClaims(aToken, tClaims, func(tToken *jwt.Token) (interface{}, error) {
		if _, ok := tToken.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", tToken.Header["alg"])
		}
		return lookupPublicKey(aAuth.publicKeyPath)
	})

	if tError != nil || !tToken.Valid {
		return nil, false, tError
	}

	return tClaims, true, nil
}

func lookupPrivateKey(tPath string) (*rsa.PrivateKey, error) {
	tKey, tError := ioutil.ReadFile(tPath)
	if tError != nil {
		fmt.Println(tError)
	}
	parsedKey, err := jwt.ParseRSAPrivateKeyFromPEM(tKey)
	return parsedKey, err
}

func lookupPublicKey(tPath string) (*rsa.PublicKey, error) {
	key, _ := ioutil.ReadFile(tPath)
	parsedKey, err := jwt.ParseRSAPublicKeyFromPEM(key)
	return parsedKey, err
}

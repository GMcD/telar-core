package utils

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/MicahParks/keyfunc"
	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/red-gold/telar-core/config"
)

func getEnvInt(key string, fallback time.Duration) time.Duration {
	if value, ok := os.LookupEnv(key); ok {
		if i, err := strconv.Atoi(value); err != nil {
			return time.Duration(i)
		}
	}
	return fallback
}

func getEnvBool(key string, fallback bool) bool {
	if value, ok := os.LookupEnv(key); ok {
		if b, err := strconv.ParseBool(value); err != nil {
			return b
		}
	}
	return fallback
}

func UnixToTime(epoch float64) string {
	tm := time.Unix(int64(epoch), 0)
	return tm.Format(time.RFC822)
}

func PrintClaim(claims jwt.MapClaims) {
	fmt.Println()
	for key, value := range claims {
		if key == "iat" || key == "exp" || key == "auth_time" {
			fmt.Printf("%s : %v\n", key, UnixToTime(value.(float64)))
		} else {
			fmt.Printf("%s : %v\n", key, value)
		}
	}
	fmt.Println()
}

func VerifyJWT(jwtB64 string) (jwt.MapClaims, error) {

	// Get the JWKs URL from your AWS region and userPoolId.
	regionID := *config.AppConfig.AwsRegion
	userPoolID := *config.AppConfig.CognitoUserPool
	jwksURL := fmt.Sprintf("https://cognito-idp.%s.amazonaws.com/%s/.well-known/jwks.json", regionID, userPoolID)

	// Create the keyfunc options. Use an error handler that logs. Refresh the JWKs when a JWT signed by an unknown KID
	// is found or at the specified interval. Timeout the initial JWKs refresh request after
	// 10 seconds. This timeout is also used to create the initial context.Context for keyfunc.Get.
	// Get Token Refresh Parameters from Environment
	refreshInterval := time.Minute * getEnvInt("JWT_REFRESH_INTERVAL_MINS", 10)
	refreshTimeout := time.Second * getEnvInt("JWT_REFRESH_TIMEOUT_SECS", 10)
	refreshUnknownKID := getEnvBool("JWT_REFRESH_UNKNWON_KID", true)

	options := keyfunc.Options{
		RefreshErrorHandler: func(err error) {
			log.Printf("There was an error with the jwt.KeyFunc\nError:%s\n", err.Error())
		},
		RefreshInterval:   &refreshInterval,
		RefreshTimeout:    &refreshTimeout,
		RefreshUnknownKID: &refreshUnknownKID,
	}

	// Create the JWKs from the resource at the given URL.
	jwks, err := keyfunc.Get(jwksURL, options)
	if err != nil {
		return nil, fmt.Errorf("Failed to create JWKs from resource at %s.\nError:%s\n", jwksURL, err.Error())
	}

	// Parse the JWT, and extract the Claims, if we have a token (which may have expired?)
	var claims jwt.MapClaims
	token, err := jwt.Parse(jwtB64, jwks.Keyfunc)
	if token != nil {
		claims = token.Claims.(jwt.MapClaims)
	}

	if err != nil {
		PrintClaim(claims)
		return nil, fmt.Errorf("Error: %s\n", err.Error())
	}

	return claims, nil

}

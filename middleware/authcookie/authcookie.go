package authcookie

import (
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofrs/uuid"
	"github.com/red-gold/telar-core/pkg/log"
	"github.com/red-gold/telar-core/pkg/parser"
	"github.com/red-gold/telar-core/types"
	"github.com/red-gold/telar-core/utils"
)

// New creates a new middleware handler
func New(config Config) fiber.Handler {
	// Set default config
	cfg := configDefault(config)

	// Return new handler
	return func(c *fiber.Ctx) error {
		// Don't execute middleware if Next returns true
		if cfg.Next != nil && cfg.Next(c) {
			return c.Next()
		}

		var auth string

		// Get authorization from header
		auth = c.Get("Authorization", "")
		if auth == "" {
			// Get authorization from cookies
			auth = readAuthCookie(c, cfg.HeaderCookieName, cfg.PayloadCookieName, cfg.SignatureCookieName)
		}

		// Check if the jwt token contains content
		jwtToken := strings.Split(auth, ".")
		if len(jwtToken[0]) == 0 || len(jwtToken[1]) == 0 || len(jwtToken[2]) == 0 {
			log.Error("Token does not contains content %s ", auth)
			return cfg.Unauthorized(c)
		}

		// Check if the JWT secret key is not nill
		if cfg.JWTSecretKey == nil {
			log.Error("JWT secret key is not provided in config!")
			return c.SendStatus(http.StatusInternalServerError)
		}

		// Check token validation and set user context in locals
		if parsedClaim, err := cfg.Authorizer(auth); err == nil && parsedClaim != nil {

			// Get standard Claims
			userCtx := new(types.UserContext)
			parser.MarshalMap(parsedClaim["claim"], userCtx)

			// Get Additional Claims
			claims, err := utils.VerifyJWT(auth)
			if err != nil {
				log.Error("Verify JWT error : %s\n", err.Error())
				return cfg.Unauthorized(c)
			} else {
				log.Info("Claims : %s", claims)
			}

			userCtx.UserID, _ = uuid.FromString(claims["cognito:username"].(string))
			userCtx.Username = claims["email"].(string)
			userCtx.DisplayName = claims["name"].(string)

			log.Info("UserContext : %s", *userCtx)

			c.Locals(cfg.UserCtxName, *userCtx)
			return c.Next()

		} else {
			log.Error("Unauthorize user due to parsedClaim : %s\vToken : %s", err.Error(), auth)
		}

		// Authentication failed
		return cfg.Unauthorized(c)
	}
}

package authcookie

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/red-gold/telar-core/pkg/log"
)

// readAuthCookie Read cookie and get auth JWT string
func readAuthCookie(c *fiber.Ctx, HeaderCookieName, PayloadCookieName, SignatureCookieName string) string {
	log.Info("JWT Token header : %s", c.Cookies(HeaderCookieName))
	log.Info("JWT Token payload : %s", c.Cookies(PayloadCookieName))
	log.Info("JWT Token signature : %s", c.Cookies(SignatureCookieName))

	return fmt.Sprintf("%s.%s.%s", c.Cookies(HeaderCookieName), c.Cookies(PayloadCookieName), c.Cookies(SignatureCookieName))
}

func authCookiePresented(c *fiber.Ctx, HeaderCookieName, PayloadCookieName, SignatureCookieName string) {

}

package common

import (
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/jtblin/go-ldap-client"
	"gopkg.in/appleboy/gin-jwt.v2"
)

// GetAuthMiddleware returns a gin middleware for JWT with cookie based auth
func GetAuthMiddleware() *jwt.GinJWTMiddleware {
	key := os.Getenv("SESSION_KEY")

	if len(key) == 0 {
		log.Fatal("Env variable 'SESSION_KEY' must be specified")
	}

	return &jwt.GinJWTMiddleware{
		Realm:         "CLOUD_SSP",
		Key:           []byte(key),
		Timeout:       time.Hour,
		MaxRefresh:    time.Hour,
		Authenticator: ldapAuthenticator,
		Authorizator: func(userId string, c *gin.Context) bool {
			return true
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		TokenLookup: "header:Authorization",
		TimeFunc:    time.Now,
	}
}

func ldapAuthenticator(userID string, password string, c *gin.Context) (string, bool) {
	ldapHost := os.Getenv("LDAP_URL")
	ldapBind := os.Getenv("LDAP_BIND_DN")
	ldapBindPw := os.Getenv("LDAP_BIND_CRED")
	ldapFilter := os.Getenv("LDAP_FILTER")
	ldapSearchBase := os.Getenv("LDAP_SEARCH_BASE")

	client := &ldap.LDAPClient{
		Base:         ldapSearchBase,
		Host:         ldapHost,
		Port:         389,
		UseSSL:       false,
		SkipTLS:      true,
		BindDN:       ldapBind,
		BindPassword: ldapBindPw,
		UserFilter:   ldapFilter,
	}
	// It is the responsibility of the caller to close the connection
	defer client.Close()

	ok, _, err := client.Authenticate(userID, password)
	if err != nil {
		log.Printf("Error authenticating user %s: %+v", userID, err)
	}
	if !ok {
		log.Printf("Authenticating failed for user %s", userID)
	}
	return userID, ok
}


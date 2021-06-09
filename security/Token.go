package security

import (
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	uuid "github.com/satori/go.uuid"
)

// Appliance ExternalID Type
const Appliance = "Appliance"

// Session ExternalID Type
const Session = "Session"

// User ExternalID Type
const User = "User"

// JwtToken represents the parsed Token from Authentication Header
type JwtToken struct {
	// UserID is id of user matchimg the token
	UserID         uuid.UUID   `json:"user,omitempty"`
	UserName       string      `json:"name,omitempty"`
	DisplayName    string      `json:"displayName,omitempty"`
	UserGroupIDs   []uuid.UUID `json:"usergroupIds,omitempty"`
	TenantID       uuid.UUID   `json:"tenant,omitempty"`
	TenantName     string      `json:"tenantName,omitempty"`
	ExternalID     string      `json:"externalId,omitempty"`
	ExternalIDType string      `json:"externalIdType,omitempty"`
	Scopes         []string    `json:"scope,omitempty"`
	Admin          bool        `json:"admin,omitempty"`
	Raw            string      `json:"-"`
	jwt.StandardClaims
}

func (token *JwtToken) isValidForScope(allowedScopes []string) bool {
	permissiveTokenScopes := []string{}
	nonPermissiveTokenScopes := []string{}

	for _, tokenScope := range token.Scopes {
		if strings.HasPrefix(tokenScope, "-") {
			nonPermissiveTokenScopes = append(nonPermissiveTokenScopes, tokenScope[1:])
		} else {
			permissiveTokenScopes = append(permissiveTokenScopes, tokenScope)
		}
	}

	if len(nonPermissiveTokenScopes) > 0 {
		if isNegativeScopePresent(nonPermissiveTokenScopes, allowedScopes) {
			return false
		}
	}

	return isPositiveScopePresent(permissiveTokenScopes, allowedScopes)
}

func isPositiveScopePresent(scopes []string, scopeToCheck []string) bool {
	if ok, _ := inArray("*", scopes); ok {
		return true
	}
	for _, allowedScope := range scopeToCheck {
		if ok, _ := inArray(allowedScope, scopes); !ok {
			scopeParts := strings.Split(allowedScope, ":")
			if ok, _ := inArray(scopeParts[0]+":*", scopes); !ok {
				if ok, _ := inArray("*:"+scopeParts[1], scopes); !ok {
					return false
				}
			}
		}
	}
	return true
}

func isNegativeScopePresent(scopes []string, scopeToCheck []string) bool {
	for _, allowedScope := range scopeToCheck {
		if ok, _ := inArray(allowedScope, scopes); ok {
			scopeParts := strings.Split(allowedScope, ":")
			if ok, _ := inArray(scopeParts[0]+":*", scopes); ok {
				if ok, _ := inArray("*:"+scopeParts[1], scopes); ok {
					return true
				}
			}
		}
	}
	return false
}

func inArray(val string, array []string) (ok bool, i int) {
	for i = range array {
		if ok = array[i] == val; ok {
			return
		}
	}
	return
}

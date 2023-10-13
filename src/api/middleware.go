package api

import (
	"errors"
	"net/http"
	"strings"

	"github.com/beatzoid/simple-bank/token"
	"github.com/gin-gonic/gin"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "authorization_payload"
)

func authMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Get authorizationHeader header
		authorizationHeader := ctx.GetHeader(authorizationHeaderKey)

		// If the authorization header is empty (the header doesn't exist)
		if len(authorizationHeader) == 0 {
			// Abort with unauthorized status and error
			err := errors.New("authorization header is not provided")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		// Splits the header into a slice of strings, which should be ["bearer", "<authorization token>"]
		fields := strings.Fields(authorizationHeader)

		// If the header doesn't contain 2 fields (missing bearer, auth token, or both)
		if len(fields) < 2 {
			// Abort with unauthorized status and error
			err := errors.New("invalid authorization header format")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		// Get the authorization type
		authorizationType := strings.ToLower(fields[0])

		// If the authorization type is not the supported type
		if authorizationType != authorizationTypeBearer {
			// Abort with unauthorized status and error
			err := errors.New("unsupported authorization type")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		// Get the access token from the header
		accessToken := fields[1]
		// Verify the access token
		payload, err := tokenMaker.VerifyToken(accessToken)
		// If it fails to verify the token
		if err != nil {
			// Abort with unauthorized status and error
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		// Save the token payload to the context
		ctx.Set(authorizationPayloadKey, payload)
		// Forward the request to the handler function
		ctx.Next()
	}
}

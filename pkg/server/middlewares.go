package server

//const CorrelationIDKey = "CorrelationID"

//func RequestIDMiddleware(logger services.LogService) echo.MiddlewareFunc {
//	return middleware.RequestIDWithConfig(middleware.RequestIDConfig{
//		RequestIDHandler: func(ctx echo.Context, s string) {
//			cID := ctx.Request().Header.Get(echo.HeaderXRequestID)
//			if cID != "" {
//				logger.Info(fmt.Sprintf("keeping existing correlationID %s", cID))
//				ctx.Set(CorrelationIDKey, cID)
//
//				return
//			}
//			cID = uuid.NewString()
//			ctx.Request().Header.Set(echo.HeaderXRequestID, cID)
//			ctx.Set(CorrelationIDKey, cID)
//		},
//	})
//}
//
//func LoggerMiddleware() echo.MiddlewareFunc {
//	return middleware.LoggerWithConfig(middleware.LoggerConfig{
//		Format: "method=${method}, uri=${uri}, status=${status} correlation_id=${header:X-Request-Id}\n",
//	})
//}
//
//func CorsMiddleware() echo.MiddlewareFunc {
//	return middleware.CORSWithConfig(middleware.CORSConfig{
//		AllowOrigins: []string{"*"},
//		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
//		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
//	})
//}
//
//// CorrelationMiddleware injects a correlation ID into the request context
//func CorrelationMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
//	return func(c echo.Context) error {
//		correlationID := c.Request().Header.Get("X-Correlation-ID")
//		if correlationID == "" {
//			correlationID = uuid.New().String() // Generate a new UUID if missing
//		}
//
//		// Store correlation ID in request context
//		ctx := context.WithValue(c.Request().Context(), CorrelationIDKey, correlationID)
//		c.SetRequest(c.Request().WithContext(ctx))
//
//		// Add correlation ID to response header
//		c.Response().Header().Set("X-Correlation-ID", correlationID)
//
//		return next(c)
//	}
//}
//
//// GetCorrelationID retrieves the correlation ID from context
//func GetCorrelationID(ctx context.Context) string {
//	if val, ok := ctx.Value(CorrelationIDKey).(string); ok {
//		return val
//	}
//	return ""
//}

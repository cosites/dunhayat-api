package docs

import "github.com/swaggo/swag"

var SwaggerInfo = &swag.Spec{
	Version:     "1.0.0",
	Host:        "localhost:8080",
	BasePath:    "/api/v1",
	Title:       "Dunhayat Coffee Roastery API",
	Description: "Backend for Dunhayat Coffee Roastery's E-Commerce System",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}

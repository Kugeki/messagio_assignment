generate-swagger:
	swag init -g internal/ports/rest/rest.go

generate-markdown-from-swagger:
	swagger generate markdown -f docs/swagger.json --output=api.md
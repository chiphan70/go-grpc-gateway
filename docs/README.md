# Documentation

This directory contains generated documentation and API specifications.

## Files

- `*.swagger.json` - OpenAPI/Swagger documentation generated from protobuf definitions
- `*.openapi.yaml` - OpenAPI specifications (if generated)

## Viewing API Documentation

You can view the generated API documentation in several ways:

### 1. Swagger UI (Recommended)

Use any Swagger UI viewer with the generated JSON files:

```bash
# Using Docker
docker run -p 8081:8080 -e SWAGGER_JSON=/docs/user.swagger.json -v $(pwd)/docs:/docs swaggerapi/swagger-ui

# Or use online viewer
# Upload the JSON file to: https://editor.swagger.io/
```

### 2. VS Code Extension

Install the "OpenAPI (Swagger) Editor" extension and open the `.swagger.json` files.

### 3. Command Line

```bash
# View the raw JSON
cat docs/user.swagger.json | jq
```

## Regenerating Documentation

Documentation is automatically generated when you run:

```bash
make buf-generate
```

The documentation is generated from the Protocol Buffer definitions in `api/proto/` and follows the annotations defined in the `.proto` files.

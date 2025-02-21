#!/bin/bash

# Create main project directory
mkdir -p ecommerce-htmx
cd ecommerce-htmx

# Create directory structure
mkdir -p cmd/server
mkdir -p internal/{config,domain,handlers,middleware,repository/mysql,service,templates,utils}
mkdir -p internal/templates/{layouts,partials,pages}
mkdir -p assets/{css,js,images}
mkdir -p migrations
mkdir -p public

# Create main Go files
touch cmd/server/main.go
touch internal/config/config.go
touch internal/domain/models.go
touch internal/handlers/handlers.go
touch internal/middleware/middleware.go
touch internal/service/service.go
touch go.mod
touch .env
touch .gitignore

# Create base template files
touch internal/templates/layouts/base.html
touch internal/templates/partials/header.html
touch internal/templates/partials/footer.html
touch internal/templates/partials/nav.html

# Create initial asset files
touch assets/css/styles.css
touch assets/js/htmx.min.js

# Create basic .gitignore
echo "# Binaries
*.exe
*.exe~
*.dll
*.so
*.dylib

# Environment variables
.env

# IDE specific files
.idea/
.vscode/
*.swp
*.swo

# Dependencies
/vendor/

# Build output
/bin/
/dist/

# Logs
*.log

# OS specific
.DS_Store" > .gitignore

# Initialize Go module
go mod init ecommerce-htmx

# Initialize Git repository
git init

echo "Project structure created successfully!"
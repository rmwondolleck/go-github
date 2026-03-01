// +build tools

package main

// This file imports packages that are used when running go generate, or are
// otherwise not referenced from source code but still required for development.
import (
	_ "github.com/google/uuid"
	_ "github.com/stretchr/testify"
	_ "github.com/swaggo/files"
	_ "github.com/swaggo/gin-swagger"
	_ "github.com/swaggo/swag"
)

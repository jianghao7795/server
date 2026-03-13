module server/tests

go 1.25

require (
	github.com/gofiber/fiber/v2 v2.52.0
	github.com/stretchr/testify v1.9.0
	server v0.0.0
)

replace server => ../

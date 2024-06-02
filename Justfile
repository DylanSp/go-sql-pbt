go_package_name := "github.com/DylanSp/go-sql-pbt"

# Using "@just" instead of "just" avoids printing out "just --list"

default:
    @just --list

up:
    docker compose up -d

down:
    docker compose down

test-unit:
    go test {{go_package_name}}/pkg/storage -run TestAllMethodsWithExistingStudent

go_package_name := "github.com/DylanSp/go-sql-pbt"

# Using "@just" instead of "just" avoids printing out "just --list"

default:
    @just --list

up:
    docker compose up -d

down:
    docker compose down

docker-clean:
    docker compose rm -f db migrate
    docker volume rm "$(docker volume ls -q)"

test-unit:
    go test {{go_package_name}}/pkg/storage -run TestAllMethodsWithExistingStudent

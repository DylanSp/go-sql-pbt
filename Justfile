go_package_name := "github.com/DylanSp/go-sql-pbt"

# Using "@just" instead of "just" avoids printing out "just --list"
[doc]
default:
    @just --list

up:
    docker compose up -d

down:
    docker compose down

docker-clean:
    docker compose rm -f db migrate
    docker volume rm "$(docker volume ls -q)"

psql:
    docker exec -it go-sql-pbt-db-1 sh -c "psql -U postgres -d school"

# run tests with -count=1 to avoid caching
[doc]
test-unit:
    go test -count=1 {{go_package_name}}/pkg/storage -run TestBasicUsage

test-fuzz:
    go test github.com/DylanSp/go-sql-pbt/pkg/storage -fuzz=FuzzBasicUsage

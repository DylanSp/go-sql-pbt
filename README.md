# Go SQL Property-Based Testing

An experiment with using [property-based testing techniques](https://increment.com/testing/in-praise-of-property-based-testing/) to test a data access layer that connects to a SQL database, as well as adapting [Go fuzz tests](https://go.dev/doc/security/fuzz/) to more general property-based testing.

Currently, this is more of a fuzzer; the only property it tests is matching that the database stores the same data as an in-memory `map`, when the database and map are updated simultaneously.

I'd also like to explore some techniques for easily and quickly running isolated integration tests against a real Postgres database, mostly based on [this blog post](https://gajus.com/blog/setting-up-postgre-sql-for-running-integration-tests).

## Fuzz testing

The fuzz test in [`pkg/storage/student_test.go`](/pkg/storage/student_test.go#L122) creates a random sequence of `SELECT`, `INSERT`, `UPDATE`, and `DELETE` operations. The test then executes these with random parameters, calling the data access methods defined in [`pkg/storage/student.go`](/pkg/storage/student.go). As well as performing these operations against a Postgres database (running in a Docker container), the fuzz test performs equivalent operations against a `map` and verifies that the database has the same data; the `map` is used as a [test oracle](https://en.wikipedia.org/wiki/Test_oracle).

To see the effectiveness of fuzz testing, we can intentionally introduce a bug that the fuzz test catches, while the unit tests don't. The `GetStudentByID` method has an erroneous query commented-out; instead of the correct query, the erroneous query simply returns the student with the lowest-sorted ID. If we comment out the correct query and uncomment the erroneous query, the unit tests still pass, because they only test a single record at a time. However, running the fuzz test catches the bug almost immediately.

## Environment setup and running tests

Requirements:

- Go (>= 1.22)
- Docker
- The [Just](https://just.systems/man/en/) task runner

All of these are satisfied by creating a Github Codespace from this repository.

To run tests, first run `just up` to start a Postgres container and run SQL migrations to update its schema. Then, `just test-unit` runs the unit tests, while `just test-fuzz` runs the fuzz test as well.

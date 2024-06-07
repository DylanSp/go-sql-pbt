# Go SQL Property-Based Testing

An experiment with using [property-based testing techniques](https://increment.com/testing/in-praise-of-property-based-testing/) to test a data access layer that connects to a SQL database, as well as adapting [Go fuzz tests](https://go.dev/doc/security/fuzz/) to more general property-based testing.

Currently, this is more of a fuzzer, without any formally stated properties; the idea is to create a random series of SQL operations, then verify that their behavior is correct by using an in-memory `map` as a [test oracle](https://en.wikipedia.org/wiki/Test_oracle).

I'd also like to explore some techniques for easily and quickly running isolated integration tests against a real Postgres database, mostly based on [this blog post](https://gajus.com/blog/setting-up-postgre-sql-for-running-integration-tests).

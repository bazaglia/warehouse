# warehouse

An implementation in Go of a warehouse. The software has mainly two entities: articles and products. A product is a composition of articles. For example, a table is a product that contanins the articles screws, legs and table top.

## Database

The data modeling suits a SQL design and we are particular interested in an ACID compliant database, since the application has relevant use cases for the ACID properties:

(A) Atomicity - this is so that each request to update a database item is not interfered by other write requests. Even if a crash happens while updating the articles stock, either all of queries will be executed, or none of them.

(C) Consistency - by maintaining data integrity constraints, the database prevents the data to enter an ilegal state. We can, therefore, assure that the amount of an article never go bellow zero.

(I) Isolation – allows concurrency, will play a part in guaranteeing the systematic process of the concurrent transaction, which means one by one. Multiple users will never be able to succeed buying a product that has a single unit available in stock. Even if requests are made at the same time, the first transaction will succeed, and the following ones will fail.

(D) Durability – Sucessful commits will survive permanently.

### Folder structure

The folder structure was designed after the following thoughts and considerations:

Why not flat: a flat structure could work for small projects and even though there are successful open source projects using that structure, it puts everything in a global context and it can be hard to understand what is going on at a first time

Why don't group by function: even though this is a common approach for other programming languages, creating a package called models or types can lead to easy circular dependencies. A repository might use a model for a definition, and that model might use the repository to call the database.

Group by context: That is the approach this project uses. We organize packages by their functional responsibilities. The flexibility introduced on structs declaration also helps to have a clear understanding of the package needs. For instance, a product struct in the listing package should have an ID, but in the create context it wouldn't necessarily need one since it would be generated by the storage package or by the database itself.

## Running the service

### To run the tests

`make test`

### To build the binaries

`make build`

Prefixing the command with the environment variables `GOOS` and `GOARCH` lets you target the binary for another platform. The list of available ones can be seen by running `go tool dist list`.

### To run the service locally

Docker is the only requirement to run the service. It will spin up a database container for Postgres and an application container for the Go server.

`make start`

## API

After the service is running, you can use [editor.swagger.io](https://editor.swagger.io/) and paste the content of the `openapi.yml` there. The UI lets you make requests to the operations exposed through the API.

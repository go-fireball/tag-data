# Time-Series API

This is a Go API for ingesting and retrieving time-series data by tag name.

## Requirements
To use this API, you will need to have the following installed:
- Go 1.18
- PostgreSQL 10 or later

## Installation
- Clone the repository:
- Change directory to the project root
- Set up the database schema

   `psql -h <database-host> -U <database-user> -f db/schema.sql`

- Update the connection string in `environment.go`
- Build and run the API `go build`
- Run the Api `go run .`    


## Usage

### Ingesting Data

```
    POST /ingest-data HTTP/1.1
    Host: localhost:9090
    PartitionKey: 098c40e0-e988-7316-9fbe-29d8781dd988
    Content-Type: application/json
    Content-Length: 85
    
    [
        {
            "tagName" : "AAPL",
            "dt": 1672552860001,
            "val": 100.00
        },
        {
            "tagName" : "GOOG",
            "dt": 1672552860001,
            "val": 99.12
        }
    ]
```
### Retrieving Data

Request 

```http request
POST /get-tag-data HTTP/1.1
Host: localhost:9090
PartitionKey: 098c40e0-e988-7316-9fbe-29d8781dd988
Content-Type: application/json
Content-Length: 85

{
    "tagName" : "AAPL",
    "start": 1672552860000,
    "end":1677797460000

}
```

Response 

```json
    {
        "seriesId": 3,
        "dt": [
            1672552860000,
            1672552860001
        ],
        "Val": [
            16777,
            "12"
        ]
    }

```

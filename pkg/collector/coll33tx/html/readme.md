## Test data

The presumption here is that web pages will change over time, 
so unit tests for parsers will fail. Keeping them up to date
is of utmost importance, and part of this is periodically
updating the test data.

```bash
go run loader.go
```

This fetches all the html pages needed for testing and encodes their contents to base64 
so search engines won't be able to point to this repo.
All data is encoded to json and kept in `dataFile` (`data.json`).
Tests can then use the `Read` function to parse the json and decode base64.

# Simple URL Pinger

## Description

Simple REST API that can ping a list of URLs and return their status.

The API is running concurrently, so it can ping multiple URLs at the same time. Also, there is a timeout for each request, so the API won't hang on a slow request, and the ability to stop the process if there is an error on any url. The response also includes the average time of pinging all urls.  

### Usage

1. Create `.env` file and copy the content from `.env.example` into it. Set your own values for the variables.
2. To run without docker, run `make run` in the root directory of the project.
3. To run with docker, run `docker-compose up` in the root directory of the project.


### API

#### POST /api/v1/ping-urls
Parameters:
- `urls` - list of URLs to ping
  - type: `array`
  - required: `true`
- `return_on_err` - if `true`, the API will return the results as soon as it encounters an error. If `false`, the API will ping all URLs and return the results.
  - type: `boolean`
  - required: `false`
  - default: `false`
- `timeout` - timeout for each request in seconds
  - type: `integer`
  - required: `false`
  - default: `5`


### Example

Request:

```
curl -X POST http://localhost:8080/api/v1/ping-urls -d '{"urls": ["http://httpbin.org/get", "http://example.com", "https://cloudflare.com/cdn-cgi/trace", "http://www.google.com", "http://www.wikipedia.org"],"return_on_err": false}'
```

Response:

```json
{
  "results": {
    "http://example.com": "OK",
    "http://httpbin.org/get": "OK",
    "http://www.google.com": "OK",
    "http://www.wikipedia.org": "OK",
    "https://cloudflare.com/cdn-cgi/trace": "OK"
  },
  "average": 241 // in milliseconds
}
```

## Architecture

<img width="627" alt="Screenshot 2024-01-04 at 21 02 01" src="https://github.com/vadimpk/url-pinger/assets/65962115/e76e0408-51eb-477e-9aec-7cd090d25919">

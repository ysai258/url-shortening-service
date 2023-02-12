# URL Shortening Service API Documentation

## Introduction

This API provides functionality for shortening URLs and redirecting to the original URL when the short URL is accessed.

## Usage

1. Clone the repository:
    `$ git clone https://github.com/ysai258/url-shortening-service`


2. Navigate into the cloned repository:
    `$ cd url-shortening-service`


3. Run docker-compose to create the containers:
    `$ docker-compose up`

    This command will create three containers:
   - One for MySQL,
   - One for the URL shortening service, and
   - Another for the worker that deletes URLs created in the last 24 hours.

4. Verify the containers are running:
    `$ docker ps`

## Endpoints

### 1. Create Short Link

Endpoint: `/shorten`  
Method: `POST`

**Request Body**

json

`{
    "url": "https://example.com"
}`

**Response**

-   200 OK:

    `{
    "shortLink": "http://localhost:8080/abcdef"
}`

### 2. Redirect Short Link

Endpoint: `/{shortLink}`  
Method: `GET`

**Response**

-   302 Found: Redirects to the original URL.
-   400 Bad Request: If the provided short link is invalid.

    `{
    "error": "Bad Request",
    "message": "Invalid link"
}`

-   502 Bad Gateway: If there was an error fetching the long link.

    `{
    "error": "Error message",
    "message": "Error in fetching long link"
}`

## Example Usage

-   Create short link:

`curl -X POST \
  http://localhost:8080/shorten \
  -H 'Content-Type: application/json' \
  -d '{
    "longLink": "https://example.com"
}'`

-   Redirect short link:

`curl -I http://localhost:8080/abcdef`

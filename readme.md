### GO RSS Aggregator

### HTTP API

#### Create user

```http
POST http://localhost:8080/v1/users
Content-Type: application/json

{
  "name":"kiquetal2"
}

```

### Create Feed

```http
POST http://localhost:8080/v1/feeds
Authorization: ApiKey 425a7237fd6c36a76b2e7400a370e67608f619118a34f9fa9b3c79314fbe21fc

{
  "name": "Something cool",
  "url": "http://example.com/feed-kiquetal-something.xml"
}

```

### Get All Feeds

```http
GET http://localhost:8080/v1/feeds

```

### Get Post by User

```http
GET http://localhost:8080/v1/posts
Authorization: ApiKey 13d2e54f454672692aa50e4119be7873da442e3846c3f89140a281cc3b32b0fb

```

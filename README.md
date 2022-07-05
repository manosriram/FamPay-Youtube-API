### fampay task - youtube api

`Go server which pulls videos metadata from youtube and stores it in DB.`

# APIS

1. `/` Returns list of videos paginated with 5 items per page.
   > Eg: localhost:5001/?page=1
2. `/search` Returns list of videos with partially or completely match the given `search` query parameter, paginated with 5 items per page.
   > Eg: localhost:5001/search?page=1&search=ronaldo

##### Response Format

`{ "success": boolean, "videos": []string }`

# Functionalities

The server spawns a `go routine` which gets videos metadata (with predefined query, "football" in our case) from youtube every 10seconds.

User can supply `multiple API keys`, first valid API key will be used everytime a request is made.

# Database

This server uses a MongoDB service. You can create one [here](https://www.mongodb.com/atlas/database)

# Config file

Create a `config.yaml` at the root path with the following template:

```
    MONGO_USERNAME:
        username
    MONGO_PASSWORD:
        password
    YOUTUBE_API_KEY:
        - api_key_from_google1
        - api_key_from_google2
        - api_key_from_google3
```

# Running the server

1. `make run` will start the server locally on port 5001.
2. `docker build . -t youtubeapi && docker run -p 5001:5001 youtubeapi` will start a docker container on port 5001.

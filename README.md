# LasramR Backend Technical Test at Scalingo

## Summary

This is my submission for the Scalingo "Hard-skills" tests.

My design focus on "Clean architecture" by ensuring that each layer of the application (transport, business logic, data fetching, ...) works independently.

With the use of common techniques such as "Dependency Injection" and "Inversion Of Control", the code I provided is decoupled and is designed to support change : the app is relying on the GitHub REST API to fetch data, wanna change to the graphQL API ? No problem, as this design of code emphasize maintainability and extendability.

Finally, given the tests constraints, attention has been given to the application performances : concurrent processing has been implemented whenever possible and data caching has been implemented in order to respond as fast as possible.

### Functionnalities

* The projects allows to retrieve aggregated data about up to 100 public Github repositories.
* By default, result will be sorted by the date they have been pushed to github but can be sorted by other
* Results can be filtered by different parameters such as : language, license, org, user, repos 
* Results contains additionnal metadata about the request such as the next page url, the total count of matching repositories on github, ...

For details see [API](#api)

## Configuration

#### Fetching source code :

```bash
git clone https://github.com/LasramR/sclng-backend-test-lasramR.git
```

#### .env based app configuration

The project rely on the usage of a `.env` files to configures internals values. These values are used to configure the app behavior and the container infrastructure. 

To start configuring the project, at project root, create a `.env` file :

```bash
touch .env
```

Then, edit the newly create `.env` file by adding variable as follows :

```py
SOME_VARIABLE=SomeValue
```

Here are the differents variables used by the projects :

| Name | Value type | Default | Optionnal |
| --- | --- | --- | --- |
|PORT | Integer between 1024 and 49152 | 5000 | Yes |
|GITHUB_API_VERSION | String | 2022-11-28 | Yes |
|GITHUB_TOKEN | String | | Yes |
|REDIS_PORT | Integer between 1024 and 49152 | 6379 | Yes |
|REDIS_PASSWORD | String | | Yes |
|CACHE_DURATION_IN_MIN | Integer > 0 | 5 | Yes |

Note: despite all these variables being optionnal, you must set up a github authentication token otherwise the app will run in limited mode (only 60 queries / hour to the GitHub REST API). To create a Github authentication token see [Github Doc](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/managing-your-personal-access-tokens#creating-a-fine-grained-personal-access-token).

## Execution

App can be booted using :

```bash
docker compose up
```

You should be able to access the app using `http://localhost:$PORT/repos`

## [API](#api)

The app exposes a single endpoint : `/repos`

### /repos

This endpoint is used to fetch aggregated results about the last public Github repositories.

#### Success

The endpoint will respond with HTTP 200 in case of a success

#### Success Response Body

The endpoint will return a jsoned object as follow :

```json
{
  "total_count": "int", // Total number of repository on github that matched the request
  "count": "int", // Number of repository returned by the API
  "content": []"null"|{ // Aggregated data from Github, null entries means that the corresponding data aggregation failed 
    "full_name": "string", // Repository full name : owner + name
    "owner": "string", // User name or organisation owning the repository
    "description": "string", // Repository description
    "repository": "string", // Repository name
    "repository_url": "string", // URL to access the repository
    "languages": { // Map of languages used in the repository
      "[key]": { // Language name
        "Bytes": "int" // Total number of bytes of this language in the repository
      }
    },
    "license": "string", // License name, can be null
    "size":, "int", // Size of the repository in bytes
    "created_at": "string", // Creation date of the repository
    "updated_at": "string" // Date of last update to the repository
  },
  "incomplete_result": "bool", // Describes if content contains null values
  "previous": "string|null", // Url pointing to the previous paginated content
  "next": "string|null" // Url pointing to the next paginated content
}
```

#### Error

If the response code is not included between 200 and 299, an error has been responded

#### Error Body

```json
{
  "status": "int", // Error HTTP status code
  "reasons": "[]string" // String array describing error
}
```

#### Filtering

Result can be filtered by different parameters using query parameters.

Here are supported filtering parameters :
* language, the main programming language of the repos 
* license, the license name of the repos
* user, the user owning the repos
* org, the organization owning the repos
* repo, the full name of the repos

All these parameters are string parameters

Usage :
* `/repos?repo=jquery/jquery`
* `/repos?language=python`
* `/repos?org=Scalingo`

Of course they can be additionned :
Usage :
* `/repos?org=Scalingo&language=Go`

#### Sorting

Result can be sorted with the use of the **sort** query parameter.

Sorting may take one of the following values :
* updated: Sort results by updated date time (ie datetime of `git commit` command)
* forks: Sort results by number of forks
* stars: Sort results by number of stars

Results are defaultly sorted by push time (ie datetime of `git push` command)

#### Limiting

Result count can be limited between 1 and 100 with the use of the **limit** query parameter.

**limit** is an integer between 1 and 100

Usage : `/repos?limit=50

### Examples

To easely run these test requests, set up the **PORT** env var on your host machine :

```bash
export PORT=...
```

* Get the last 100 Github repositories sorted by updated datetime :

```bash
curl http://localhost:$PORT/repos?sort=updated > last100Updated.json
```

* Get the last 10 Go repos of the Scalingo organization

```bash
curl http://localhost:$PORT/repos?org=Scalingo&language=Go&limit=10 > scalingoLast10GoRepos.json
```

* Get the jquery/jquery repository

```bash
curl http://localhost:$PORT/repos?repo=jquery/jquery > jqueryRepository
```

### Project structure

The project structure is as follows :

```
/project-root
│
├── /api
│   └── (API related file providing transport)
│
├── /builder
│   └── (provide mecanism to build complex queries)
│
├── /model
│   └── (definition of structs from the app and fetched apis)
│
├── /providers
│   └── (abstraction of external mecanisms eg http, cache)
│
├── /repositories
│   └── (abstraction of data access)
│
├── /services
│   └── (business logic)
│
├── /utils
│   └── (utilities used accross the app)
│
├── config.go (env configuration)
└── main.go   (entry point and app setup)
```

Boot sequence :
- main.go will starts a new process
- config.go will be used to parse in environment variable
- main.go will created and configure the different layers
  - initializing providers
  - initializing services
  - initializing http server and defining route handlers

As a bonus, unit tests have been written, these are the *_test.go files next to the source code.

### Network structure

Running the app

```mermaid
flowchart LR
	O(<b>Outside world</b>)
    G(Github API)
    subgraph Docker network
        A(<b>App</b> app:$PORT)
        R(<b>Redis Cache</b> redis:$REDIS_PORT)
    end
    A --> |exposes port| O
    A --> |Insert and retrieve from| R
    A --> |fetch data from| G
```

### Third parties

## Design choice and operation

### API Performances

### Clean architecture
---
title: Project `NoteVook`
author: [Ulises Tirado Zatarain]
date: 2023-05-19
subject: Markdown
keywords: [golang, api, sqlite]
lang: en-gb
...

# 🧑🏽‍💻 Project: `NoteVook` 
[![License](https://img.shields.io/badge/License-BSD_3--Clause-blue.svg)](https://opensource.org/licenses/BSD-3-Clause) [![CI/CD Pipeline](https://github.com/zatarain/note-vook/actions/workflows/pipeline.yml/badge.svg)](https://github.com/zatarain/note-vook/actions/workflows/pipeline.yml) [![codecov](https://codecov.io/gh/zatarain/note-vook/branch/main/graph/badge.svg?token=bufQuVyLCi)](https://codecov.io/gh/zatarain/note-vook) [![Go Report Card](https://goreportcard.com/badge/github.com/zatarain/note-vook)](https://goreportcard.com/report/github.com/zatarain/note-vook)

This project aims to be an exercise to discuss about software engineering technical topics like software development, pair programming, testing, deployment, etcetera. More specifically, to discuss the development of an [API (Application Programming Interface)][what-is-api] to **manage annotations for videos** implemented written in [go programming language][go-lang].

## 📂 Table of content

* 📹 [Overview](#📹-overview)
  - ☑️ [Requirements](#-requirements)
  - 🤔 [Assumptions](#-assumptions)
* 📐 [Design](#-design)
  - 📊 [Data model](#-data-model)
    * 🎞️ [Video](#-video)
    * ✍🏽 [Annotation](#-annotation)
    * 👤 [User](#-user)
  - 🔀 [Workflows](#-workflows)
    * 🔀 [User sign up](#-user-sign-up)
    * 🔀 [User login](#-user-login)
    * 🔀 [Authorised requests](#-authorised-requests)
  - 🔚 [End-points](#-end-points)
* 🏗️ [Implementation details](#-implementation-details)
  - 📦 [Dependencies](#-dependencies)
  - 🗄️ [Storage](#-storage)
* ⏯️ [Running](#-running)
  - 🍏 [Development Mode](#-development-mode)
  - 🍎 [Production Mode](#-production-mode)
* ✅ [Testing](#-testing)
  - 🧪 [Manual](#-manual)
  - ♻️ [Automated](#-automated)
  - 💯 [Coverage](#-coverage)
* 📚 [References](#-references)

## 📹 Overview
This simple API aims to manage a video annotations database, this means we will have a collection of videos and we will be able to add text notes for some given interval of time of the videos. Additionally we would like to manage some basic security layer based-on [JWT (JSON Web Token)][what-is-jwt]. The application should be ready to deploy as a [Docker][docker] container, so we need to generate an image for it available to download in [Docker Hub][docker-hub].

### ☑️ Requirements
The API should be able to manage a database for the videos and each video may have many annotations related to it. An annotation allow s to capture time related information about the video. For example, we may want to create an annotation that references to a part of the video from `00:04:00` to `00:05:00`.

 Allowing the client to perform following operations:

 * **List all the vide.** It should return the list of all videos in the system.
 * **Create a vide.** It should insert a new record for video provided by the client that includes some metadata.
 * **Update a video.** It should allow the client to update the information of a given video.
 * **Delete a video.** The API should provide an end-point to delete videos from the system.
 * **View video details.** It should show the details of a single video provided by the client.
 * **Annotate a video.** It should create a new annotation for a video with `start` and `end` time.
 * **Update annotations.** The API should allow the client to update annotation details.
 * **Delete annotations.** The API should provide a mechanism to delete annotations.
 * **Security layer.** It should implement a security layer based-on JWT.

### 🤔 Assumptions
This is a small example and it's not taking care about some corner case scenarios like following:

* The environment variables and secrets (e. g. `SECRET_TOKEN_KEY` to encode sign the authorisation token) for API configuration are stored in `.env` files (see [Running section](#-running) below for more information).
* In the real world the secrets should be stored and provisioned by an external system (e. g. AWS Secret Manager). In order to test and play around with the API you can leave them as blank string in the `.env` files.
* The videos can only be annotated by the user creator.
* In order to keep things simple, there is no ACID transactions for the database. We will remove the annotations in cascade though, so if we delete a vide from database we will remove its annotations too.
* The users won't be able to edit the video ID for an annotation. If the users want to do so, it's better to remove the annotation from the video, then add a new one in the other video. 
* A video with the same link can be added multiple times by different users.
* It's been assumed that the annotation type it's some sort of category and each annotation can only be of one type.
* Users were not part of the original requirements, but I added them as makes simpler the way to explain the authorisation layer.
* The user names and passwords aren't validated properly, so the client can provide any input except an empty string.
* Users can anonymously be created in the system.
* If we would like access to the API end-point programmatically (e.g. via some automation), we would need to create a new user and their correspondent password for that client.
* Even if we added the security layer with the authorisation process, this is not secure enough, there are several flaws (e. g. non-secure cookie, non-password charset checking, lack of HTTPS certificates, etcetera), but it's implemented in this way just for didactical purposes.

## 📐 Design
The architecture will be a HTTP API for a microservice that will consume some configuration and use ORM to represent the records in the database tables and also a Model-Controller (MC) pattern design, so the controllers will contain the handlers for the API requests, while the models will represent the data. The service will be stateless, so we won't hold any state (e. g. session management) on the server side, instead we will use authorisation tokens.

### 📊 Data model
In order to store and manipulate the data needed the API will rely on the entities shown in following diagram:

```mermaid
erDiagram
  Video {
    id integer PK
    user_id integer FK
    title string
    description string
    link integer
    duration integer
    created_at datetime
    updated_at datetime
  }

  Annotation {
    id integer PK
    video_id integer FK
    type enum
    start integer
    end integer
    title string
    notes string
    created_at datetime
    updated_at datetime
  }

  User {
    id integer PK
    nickname string
    password string
    created_at datetime
    updated_at datetime
  }

  User ||--o{ Video : "may own"
  Annotation }o--|| Video: "may have"

```

As we can see in the diagram, a `User` _may own_ several `Video`s, which is made possible with the user of the foreign key `user_id` within the `Video` entity. Then, a `Video` _may have_ many `Annotation`s thanks to the foreign key `video_id`.

The API manages the persistency of the data with a 🪶 [SQLite][sqlite] database, which is a simple local storage database. So, the records for entities shown in the diagram will be stored as rows in tables. SQLite manages a [reduced set of data types][sqlite-data-types] so, we will use the actual data type (affinity) used in following subsections.

#### 🎞️ Video
This entity will represent the videos in the system and each record will be stored in the table `videos` which has following fields:

| ⏹️ | Name          |     Type    | Description                                           |
|:--:| :---          |    :----:   | :---                                                  |
| 🗝️ | `id`          | `INTEGER`   | Auto-numeric identifier for the video                 |
| ✳️ | `user_id`     | `INTEGER`   | Foreign key for the user owner of the video           |
| 🔤 | `title`       | `TEXT`      | Title of the video                                    |
| 📄 | `description` | `BLOB`      | Description for the video                             |
| 🔤 | `link`        | `TEXT`      | URL for a link of the video. Unique along user domain |
| 🔢 | `duration`    | `INTEGER`   | Duration of the video in seconds                      |
| 🗓️ | `created_at`  | `NUMERIC`   | Timestamp representing the creation time              |
| 🗓️ | `updated_at`  | `NUMERIC`   | Timestamp representing the last update time           |

#### ✍🏽 Annotation
This entity will represent the annotations for the videos in the system and each record will be stored in the table `annotations` which has following fields:

| ⏹️ | Name          |     Type    | Description                                      |
|:--:| :---          |    :----:   | :---                                             | 
| 🗝️ | `id`          | `INTEGER`   | Auto-numeric identifier for the annotation       |
| ✳️ | `video_id`    | `INTEGER`   | Foreign key for the video                        |
| 🔢 | `type`        | `INTEGER`   | Annotation type or category                      |
| 🔢 | `start`       | `INTEGER`   | Start point in seconds within the video timeline |
| 🔢 | `end`         | `INTEGER`   | End point in seconds within the video timeline   |
| 🔤 | `title`       | `TEXT`      | Title or headline of the annotation              |
| 📄 | `notes`       | `BLOB`      | Optional. Additional notes                       |
| 🗓️ | `created_at`  | `NUMERIC`   | Timestamp representing the creation time         |
| 🗓️ | `updated_at`  | `NUMERIC`   | Timestamp representing the last update time      |

#### 👤 User
The records for this entity will represent the users in the system and each record will be stored in the table `users` which has following fields:

| ⏹️ | Name          |     Type    | Description                                   |
|:--:| :---          |    :----:   | :---                                          |
| 🗝️ | `id`          | `INTEGER`   | Auto-numeric identifier for the user          |
| ✳️ | `nickname`    | `TEXT`      | Nickname of the user. Unique along this table |
| 🔢 | `password`    | `TEXT`      | Hash for the password of the user             |
| 🗓️ | `created_at`  | `NUMERIC`   | Timestamp representing the creation time      |
| 🗓️ | `updated_at`  | `NUMERIC`   | Timestamp representing the last update time   |

### 🔀 Workflows
There are three general workflows in this API: user sign up, user login and all the other operations that require authorisation.

#### 🔀 User sign up
Following diagram describes the happy path for a user signup operation:

```mermaid
sequenceDiagram

actor Unknown
participant gin.Engine
participant UsersController
participant User
participant bcrypt
participant GORM
participant Database

Unknown->>+gin.Engine: POST /signup @JSON: credentials
gin.Engine->>+UsersController: Signup(@gin.Context)
UsersController->>+gin.Engine: bind(@JSON credentials)
gin.Engine->>-UsersController: returns @credentials
UsersController->>+bcrypt: Hash(@credentials.Password)
bcrypt-->>-UsersController: returns Hashed Password
UsersController->>+User: create instance (@credentials)
User-->>-UsersController: returns @user
UsersController->>+GORM: Create(@user)
GORM->>GORM: generates SQL Statement
GORM->>+Database: query(INSERT INTO users(...) VALUES(...))
Database-->>-GORM: returns query result
GORM-->>-UsersController: populated @user
UsersController-->>-gin.Engine: HTTP 201 OK and message
gin.Engine->>-Unknown: HTTP 201 Created (JSON with message)
```

#### 🔀 User login
Following diagram describes the happy path for a user login operation:

```mermaid
sequenceDiagram

actor Unknown
participant gin.Engine
participant UsersController
participant User
participant bcrypt
participant JWT
participant GORM
participant Database

Unknown->>+gin.Engine: POST /login @JSON: credentials
gin.Engine->>+UsersController: Login(@gin.Context)
UsersController->>+gin.Engine: bind(@JSON credentials)
gin.Engine->>-UsersController: returns @credentials
UsersController->>+User: create empty instance
User-->>-UsersController: returns empty @user
UsersController->>+GORM: First(@user, @user.Nickname)
GORM->>GORM: generates SQL Statement
GORM->>+Database: query(SELECT * FROM users WHERE nickname = @user.Nickname)
Database-->>-GORM: returns query result
GORM-->>-UsersController: returns populated @user
UsersController->>+bcrypt: CompareHashAndPassword(@credentials.Password)
bcrypt-->>-UsersController: returns comparison result
UsersController->>+JWT: generates token signed with SECRET_TOKEN_KEY
JWT-->>-UsersController: returns signed @jwt.Token(7 days for expiration)
UsersController-->>-gin.Engine: HTTP 200 OK and message
gin.Engine->>-Unknown: HTTP 200 OK JSON with message and Authorisation Cookie with @jwt.Token
```

#### 🔀 Authorised requests
Following diagram describes the happy path for any other operation that needs authorisation, meaning after user has been logged in:

```mermaid
sequenceDiagram

actor KnownUser
participant gin.Engine
participant UsersController
participant User
participant JWT
participant GORM
participant AnotherController
participant AnotherModel
participant Database

KnownUser->>+gin.Engine: GET /another-controller/action @JSON: credentials
gin.Engine->>+UsersController: Authorise(@gin.Context)
UsersController->>+gin.Engine: get Authorisation cookie
gin.Engine->>-UsersController: returns @cookie

UsersController->>+JWT: Parse @cookie.Value
JWT-->>+UsersController: checks algorithm for consistency and ask for key
UsersController-->>-JWT: result for algorithm check and SECRET_TOKEN_KEY
JWT-->>-UsersController: decodes the token and returns Parsed and Decoded @jwt.Token

UsersController->>UsersController: Extracts claims and check expiration
UsersController->>+User: New instance with Nickname only
User-->>-UsersController: returns @user with Nickname only
UsersController->>+GORM: First(@user, @user.Nickname)
GORM->>GORM: generates SQL Statement
GORM->>+Database: query(SELECT * FROM users WHERE nickname = @user.Nickname)
Database-->>-GORM: returns query result
GORM-->>-UsersController: returns populated @user 
UsersController-->>-gin.Engine: @user within the @gin.Context

gin.Engine->>+AnotherController: ActionHandler (@gin.Context with @user)
AnotherController->>+AnotherModel: do something
AnotherModel-->>-AnotherController: returns some result
AnotherController->>+GORM: do something
GORM->>GORM: generates SQL statement
GORM->>+Database: query (generated SQL statement)
Database-->-GORM: returns query result
GORM-->-AnotherController: returns query result
AnotherController->>AnotherController: May do something else (e. g. business logic)
AnotherController-->>-gin.Engine: HTTP 200 OK and message
gin.Engine->>-KnownUser: HTTP 200 OK JSON with message
```

### 🔚 End-points
The input for all the API end-points will be always in JSON format and the Cookie `Authorisation` JWT token in most of the cases and the output will be in the same format. The end-points for the API are described in following table:

| Method   | Address            | Description                             | Success Status | Possible Failure Status                                |
| :---:    | :---               | :----                                   | :---:          | :---                                                   |
| `HEAD`   | `/health`          | Service health check                    | `200 OK`       | `* Any`                                                |
| `POST`   | `/signup`          | User sign up to create users            | `201 Created`  | `400 Bad Request`                                      |
| `POST`   | `/login`           | User login and get authorisation token  | `200 OK`       | `400 Bad Request`, `500 Internal Server Error`         |
| `GET`    | `/videos`          | List of all videos owned by logged user | `200 OK`       | `401 Unauthorised`                                     |
| `POST`   | `/videos`          | Create a video record in the system     | `200 Created`  | `401 Unauthorised`, `400 Bad Request`                  |
| `GET`    | `/videos/:id`      | Get video details and its annotations   | `200 OK`       | `401 Unauthorised`, `404 Not Found`                    |
| `PATCH`  | `/videos/:id`      | Edit details for a given video          | `200 OK`       | `401 Unauthorised`, `400 Bad Request`, `404 Not Found` |
| `DELETE` | `/videos/:id`      | Delete a video from the system          | `200 OK`       | `401 Unauthorised`, `404 Not Found`                    |
| `POST`   | `/annotations`     | Create a annotation record for a video  | `200 Created`  | `401 Unauthorised`, `400 Bad Request`                  |
| `PATCH`  | `/annotations/:id` | Edit details for an annotation          | `200 OK`       | `401 Unauthorised`, `400 Bad Request`, `404 Not Found` |
| `DELETE` | `/annotations/:id` | Delete an annotation                    | `200 OK`       | `401 Unauthorised`, `404 Not Found`                    |

## 🏗️ Implementation details
We are using Golang as programming language for the implementation of the API operations. And the database is a single table in SQLite stored locally.

There is a continuous integration workflow that runs in [GitHub Actions][github-actions] which is responsible to build the API, tun the unit tests if the tests succeed then it generates and pushes the image for the container to DockerHub.

### 📦 Dependencies
We are using following libraries for the implementation:

 * **`gin-gonic`.** A web framework to implement a RESTful API via HTTP.
 * **`gorm`.** A library for Object Relational Model (ORM) in order to represent the records in the database as relational objects.
 * **`gorm/drivers/sqlite`.** Driver that manage SQLite dialect and connect to the database.
 * **`godotenv`.** This CLI tool allows us to load environment configuration via `.env` files and run a command.
 * **`crypto/bcrypt`.** This is part of the standard go library. It's to make use of hashing when sign up and login.
 * **`golang-jwt`.** To generate and use the authorisation tokens.

And also, following ones for the development:
 * **`testify`.** To have more readable assertions on the unit testing.
 * **`mockery`.** To generate mocks used on unit testing.
 * **`monkey`.** To perform monkey patching on the unit testing.

### 🗄️ Storage
A Docker container it's not persistent itself, so the Docker Compose file specify a volume to make the database persistent, that volume can be mapped to a host directory.

## ⏯️ Running
In order to run the application locally you will need to have Docker installed and internet connection. Using the command line with docker you can either go on two modes:

* **Production:** Directly downloading [the latest built of the image][note-vook-image] and run it
* **Development:** Clone [this Git repository][note-vook-repo] and build the image for the container locally

### 🍎 Production Mode
You can `docker` CLI to download the [latest built of the image][note-vook-image] from Docker Hub, and then run create and run a container as follow:
```sh
docker pull zatarain/note-vook:latest
docker run \
  -v $(pwd)/data:/api/data \
  -e ENVIRONMENT=prod \
  --name notevook \
  -p 4000:4000 \
  zatarain/note-vook:latest
```

Those commands will download the latest build image generated by the [Continuous Integration and Deployment Pipeline][ci-cd-pipeline] in the [Github Actions of the repository][notevook-actions].

As you can see you need to specify following things:

* **Volume for Database [`-v $(pwd)/data:/api/data`]:** This will create (if it doesn't already exists) a `data` directory in the current directory where the API will store the database.
* **Environment [`ENVIRONMENT=prod`]:** This is optional, if you omit this, you will run it in development mode. The value of this variable determines the environment file and database tha will be used:

| `ENVIRONMENT` | File       | Database       | Description                                          |
|     :---:     | :---:      | :---:          | :---                                                 |
|               | `.env`     | `data/beta.db` | Development is enabled by leaving the variable empty |
| `test`        | `test.env` | `data/test.db` | This is used when running the Unit Testing           |
| `prod`        | `prod.env` | `data/prod.db` | Production environment                               |

* **Port binding [`-p 4000:4000`]:** The image it's built to run the API on port `4000` withing the container, but you can choose to run it in another host port if you want (e. g. `-p 8080:4000`).

### 🍏 Development Mode
In your terminal, clone repository and build image as follow:
```sh
git clone https://github.com/zatarain/note-vook.git
cd note-vook
docker compose up --build
```

That will follow the configuration specified in the [`compose.yml`][compose-yml] file to build the image and run the unit testing on building time, and then run the API in development mode.

Then you can follow the steps to play manually with the API with the steps in next section.

## ✅ Testing
...

### 🧪 Manual
...

### ♻️ Automated
...

### 💯 Coverage
You can follow the test coverage reports of this project in the [CodeCov website][codecov-notevook]:

![Icicle][codecov-icicle]

## 📚 References
* [SQLite Data Types][sqlite-data-types]
* [GORM Documentation][gorm-docs]
* [Gin Documentation][gin-docs]
* [Testify Documentation][testify-docs]
* [Monkey Patching Documentation][monkey-docs]
* [Crypto/Bcrypt Documentation][bcrypt-docs]
* [GoJWT Documentation][go-jwt-docs]
* [GoDotEnv][go-dotenv-docs]
* [GitHub Actions Documentation][github-actions-docs]

---

[what-is-api]: aws.amazon.com/what-is/api
[what-is-jwt]: https://jwt.io/introduction
[docker]: https://www.docker.com
[docker-hub]: https://hub.docker.com
[note-vook-image]: https://hub.docker.com/repository/docker/zatarain/note-vook/tags
[note-vook-repo]: https://github.com/zatarain/note-vook
[go-lang]: https://go.dev
[sqlite]: https://www.sqlite.org
[sqlite-data-types]: https://www.sqlite.org/datatype3.html
[gorm-docs]: https://gorm.io/docs/
[gin-docs]: https://gin-gonic.com/docs/
[mockery-docs]: https://vektra.github.io/mockery/
[testify-docs]: https://github.com/stretchr/testify#readme
[monkey-docs]: https://github.com/bouk/monkey#readme
[bcrypt-docs]: https://pkg.go.dev/golang.org/x/crypto/bcrypt
[go-jwt-docs]: https://github.com/golang-jwt/jwt#readme
[go-dotenv-docs]: https://github.com/joho/godotenv#readme
[github-actions]: https://github.com/features/actions
[github-actions-docs]: https://docs.github.com/en/actions
[codecov-notevook]: https://app.codecov.io/gh/zatarain/note-vook
[codecov-sunburst]: https://codecov.io/gh/zatarain/note-vook/branch/main/graphs/sunburst.svg?token=bufQuVyLCi
[codecov-grid]: https://codecov.io/gh/zatarain/note-vook/branch/main/graphs/tree.svg?token=bufQuVyLCi
[codecov-icicle]: https://codecov.io/gh/zatarain/note-vook/branch/main/graphs/icicle.svg?token=bufQuVyLCi
[notevook-actions]: https://github.com/zatarain/note-vook/actions
[ci-cd-pipeline]: https://github.com/zatarain/note-vook/blob/main/.github/workflows/pipeline.yml
[compose-yml]: https://github.com/zatarain/note-vook/blob/main/compose.yml

# 🧑🏽‍💻 Project: `NoteVook`
[![License](https://img.shields.io/badge/License-BSD_3--Clause-blue.svg)](https://opensource.org/licenses/BSD-3-Clause) [![CI/CD Pipeline](https://github.com/zatarain/note-vook/actions/workflows/pipeline.yml/badge.svg)](https://github.com/zatarain/note-vook/actions/workflows/pipeline.yml) [![codecov](https://codecov.io/gh/zatarain/note-vook/branch/main/graph/badge.svg?token=bufQuVyLCi)](https://codecov.io/gh/zatarain/note-vook)

This project aims to be an exercise to discuss about software engineering technical topics like software development, pair programming, testing, deployment, etcetera. More specifically, to discuss the development of an [API (Application Programming Interface)][what-is-api] to **manage annotations for videos** implemented written in [go programming language][go-lang].
## 📂 Table of content
* 📹 [Overview](#📹-overview)
	- ☑️ [Requirements](#☑️-requirements)
	- 🤔 [Assumptions](#🤔-assumptions)
* 📐 [Design](#📐-design)
	- 📊 [Data model](#📊-data-model)
		* 🎞️ [Video](#🎞️-video)
		* ✍🏽 [Annotation](#✍🏽-annotation)
		* 👤 [User](#👤-user)
	- 🔀 [Workflows](#🔀-workflows)
* 🏗️ [Implementation details](#🏗️-implementation-details)
  - 📦 [Dependencies](#📦-dependencies)
	- 🗄️ [Storage](#🗄️-storage)
* 📚 [References](#📚-references)

## 📹 Overview
This simple API aims to manage a video annotations database, this means we will have a collection of videos and we will be able to add text notes for some given interval of time of the videos. Additionally we would like to manage some basic security layer based-on [JWT (JSON Web Token)][what-is-jwt]. The application should be ready to deploy as a [Docker][docker] container, so we need to generate an image for it available to download in [Docker Hub][docker-hub].

### ☑️ Requirements
The API should be able to manage a database for the videos and each video may have many annotations related to it. An annotation allow s to capture time related information about the video. For example, we may want to create an annotation that references to a part of the video from `04:00:00` to `00:05:00`.

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
 * The videos can only be annotated by the user creator.
 * A video with the same link can be added multiple times.
 * Users were not part of the original requirements, but I added them as makes simpler the way to explain the authorisation layer.
 * The user names and passwords aren't validated properly, so the client can provide any input except an empty string.
 * Users can anonymously be created in the system.
 * If we would like access to the API end-point programmatically (e.g. via some automation), we would need to create a new user and their correspondent password for that client.
 * Even if we added the security layer with the authorisation process, this is not secure enough, there are several flaws (e. g. non-secure cookie, non-password charset checking, lack of HTTPS certificates, etcetera), but it's implemented in this way just for didactical purposes.

## 📐 Design
The architecture will be a HTTP microservice that will consume some configuration and use ORM to represent the records in the database tables and also a Model-Controller (MC) pattern design, so the controllers will contain the handlers for the API requests, while the models will represent the data. The service will be stateless, so we won't hold any state (e. g. session management) on the server side, instead we will use authorisation tokens.

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
		body string
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

The API manages the persistency of the data with a 🪶 [SQLite][sqlite] database, which is a simple local storage database. So, the records for entities shown in the diagram will be stored as rows in tables. SQLite manages a reduced set of data types so, we will describe the actual data type used in following subsections.

#### 🎞️ Video
This entity will represent the videos in the system and each record will be stored in the table `videos` which has following fields:

|    | Name          |     Type    | Description                                  |
|:--:| :---          |    :----:   | :---                                         |
| 🗝️ | `id`          | `INTEGER`   | Auto-numeric identifier for the video        |
| ✳️ | `user_id`     | `INTEGER`   | Foreign key for the user owner of the video  |
| 🔤 | `title`       | `TEXT`      | Title of the video                           |
| 📄 | `description` | `BLOB`      | Description for the video                    |
| 🔤 | `link`        | `TEXT`      | Link for the video. Unique along user domain |
| 🔢 | `duration`    | `INTEGER`   | Duration of the video in seconds             |
| 🗓️ | `created_at`  | `NUMERIC`   | Timestamp representing the creation time     |
| 🗓️ | `updated_at`  | `NUMERIC`   | Timestamp representing the last update time  |

#### ✍🏽 Annotation
This entity will represent the annotations for the videos in the system and each record will be stored in the table `annotations` which has following fields:

|    | Name          |     Type    | Description                                 |
|:--:| :---          |    :----:   | :---                                        |
| 🗝️ | `id`          | `INTEGER`   | Auto-numeric identifier for the annotation  |
| ✳️ | `video_id`    | `INTEGER`   | Foreign key for the video                   |
| 🔢 | `type`        | `INTEGER`   | Annotation type or category                 |
| 🔢 | `start`       | `INTEGER`   | Start point in the video timeline           |
| 🔢 | `end`         | `INTEGER`   | End point in the video timeline             |
| 🔤 | `title`       | `TEXT`      | Title or headline of the annotation         |
| 📄 | `body`        | `BLOB`      | Optional. Additional notes                  |
| 🗓️ | `created_at`  | `NUMERIC`   | Timestamp representing the creation time    |
| 🗓️ | `updated_at`  | `NUMERIC`   | Timestamp representing the last update time |

#### 👤 User
The records for this entity will represent the users in the system and each record will be stored in the table `users` which has following fields:

|    | Name          |     Type    | Description                                   |
|:--:| :---          |    :----:   | :---                                          |
| 🗝️ | `id`          | `INTEGER`   | Auto-numeric identifier for the user          |
| ✳️ | `nickname`    | `INTEGER`   | Nickname of the user. Unique along this table |
| 🔢 | `password`    | `INTEGER`   | Annotation type or category                   |
| 🗓️ | `created_at`  | `NUMERIC`   | Timestamp representing the creation time      |
| 🗓️ | `updated_at`  | `NUMERIC`   | Timestamp representing the last update time   |

### 🔀 Workflows
...
```mermaid
sequenceDiagram
```

## 🏗️ Implementation details
We are using Golang as programming language for the implementation of the API operations. And the database is a single table in SQLite stored locally.

### 📦 Dependencies
We are using following libraries for the implementation:
 * **`gin-gonic`.** A web framework to implement a RESTful API via HTTP.
 * **`gorm`.** A library for Object Relational Model (ORM) in order to represent the records in the database as relational objects.
 * **`gorm/drivers/sqlite`.** Driver that manage SQLite dialect and connect to the database.
 * **`godotenv`.** This CLI tool allows us to load environment configuration via `.env` files and run a command.
 * **`crypto/bcrypt`.** To make use of `base64` encoding and decoding for the authentication token.
 * **`golang-jwt`.** To generate and use the authorisation tokens.

And also, following ones for the development:
 * **`testify`.** To have more readable assertions on the unit testing.
 * **`mockery`.** To generate mocks used on unit testing.
 * **`monkey`.** To perform monkey patching on the unit testing.

### 🗄️ Storage
A Docker container it's not persistent itself, so the Docker Compose file specify a volume to make the database persistent, that volume can be mapped to a host directory.

## 📚 References

[what-is-api]: aws.amazon.com/what-is/api
[what-is-jwt]: https://jwt.io/introduction
[docker]: https://www.docker.com
[docker-hub]: https://hub.docker.com
[go-lang]: https://go.dev
[sqlite]: https://www.sqlite.org

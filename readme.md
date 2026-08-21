# BandPlan

BandPlan is a full-stack web application built with Go, HTMX, PostgreSQL, and WebSockets that helps bands organize songs, setlists, files, and communication in one place.

It was built from scratch without a full-stack framework to deepen my understanding of backend architecture, authentication, real-time communication, and production deployment.

Actively being developed with new features added regularly.

## Tech Stack

**Backend**
- Go
- PostgreSQL
- Gorilla WebSocket
- Goose

**Frontend**
- HTMX
- HTML
- CSS
- JavaScript

**Infrastructure**
- Docker
- Railway
- Cloudflare


![Go](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-4169E1?style=for-the-badge&logo=postgresql&logoColor=white)
[![WebSockets](https://img.shields.io/badge/WebSockets-Gorilla-222222?style=for-the-badge&logo=go&logoColor=white)](https://github.com/gorilla/websocket)
[![Goose](https://img.shields.io/badge/Goose-v3.27.2-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://github.com/pressly/goose)
[![Docker](https://img.shields.io/badge/Docker-2496ED?style=for-the-badge&logo=docker&logoColor=white)](https://www.docker.com/)
![HTMX](https://img.shields.io/badge/HTMX-3366CC?style=for-the-badge&logo=htmx&logoColor=white)
![HTML5](https://img.shields.io/badge/HTML5-E34F26?style=for-the-badge&logo=html5&logoColor=white)
![CSS3](https://img.shields.io/badge/CSS3-1572B6?style=for-the-badge&logo=css3&logoColor=white)


## Live Demo

app.bandplan.app

## Highlights

* Production deployment on Railway
* Docker containerized
* PostgreSQL database
* Real-time chat using WebSockets
* Session-based authentication
* HTMX server-rendered UI
* Responsive mobile-first design
* iTunes Search API integration
* Cloudflare domain and CDN


## Technical Challenges

* Built real-time messaging using Gorilla WebSocket.
* Implemented authentication middleware with session cookies.
* Solved HTMX browser history and refresh issues.
* Built reusable Go template components.
* Designed a normalized PostgreSQL schema for users, bands, songs, and setlists.
* Implemented image upload and processing.
* Containerized the application with Docker and deployed it to Railway.


## Current Features
* User registration and login
* Session-based authentication
* Band-specific chat
* Song management pages
* Song artwork upload support
* Profile image support
* Setlist pages
* Calendar, files, goals, and event page foundations
* iTunes API client for music metadata lookup


## Planned Features
* Admin privileges for band leaders
* Google Drive integration
* Calendar sync
* Chat-based event creation
* Goal tracking and task management/tracking
* File sharing from R2 cloud storage
* Desktop and mobile applications
* Native iOS application


## What I Learned

Building BandPlan has given me hands-on experience with:

- Designing a multi-table PostgreSQL schema
- Building REST-style HTTP handlers in Go
- Authentication and authorization middleware
- Real-time communication with Gorilla WebSocket
- HTMX server-side rendering
- Docker containerization
- Deploying Go applications to Railway
- Consuming third-party APIs
- Responsive UI design



## Feature Breakdown

The 'Songs' feature will contain a list of all yours songs with all the information you need about your immediately available. Things like recording tempo, live tempo, key, lyrics, and song notes are now just a click away for easy reference when you need them. 

The 'Setlists' feature is a dynamic setlist builder that allows you to create multiple custom setlists for quick reference. Each song will have all the information from the 'Songs' feature a click away and allow for custom segments like intros and breaks. All segments will contain a length feature and a total setlist time along with a built in timer so you can time your rehearsals and save the segment times. This will allow users to know exactly how long their set is down to the second. 

The 'Goals' feature is designed around custom band goals you create. Each goal on the list will open and allow you to create sub goals, assign tasks, and create todo lists. 

The 'Calendar' will be a built in calendar designed to sync with any existing personal calendars and offer a single source of truth for all band events. Users will also be able to create calendar events from the chat. 

The 'Files' feature will sync with Google Drive and/or Dropbox to allow access to all your bands files with the organization your used to. Create new files, and folders and manage your drive without having to change apps and easily share to the chat with one click of a button. 




<p align="center">

<img width="40%" alt="Screenshot 2026-07-16 at 2 29 20 PM" src="https://github.com/user-attachments/assets/c6f3ef98-ced7-4207-a0d9-3ea771d33364" />
<img width="40%" alt="Screenshot 2026-07-16 at 2 29 31 PM" src="https://github.com/user-attachments/assets/b18d267c-57b9-4ec0-9b49-21a6ad57bb80" />
<img width="40%" alt="Screenshot 2026-07-16 at 2 24 14 PM" src="https://github.com/user-attachments/assets/1b3b9449-9b87-4db0-8cca-7ea9fe6fc8b7" />
<img width="40%" alt="Screenshot 2026-07-16 at 2 28 08 PM" src="https://github.com/user-attachments/assets/5239cb31-acd5-457d-ae05-ee2cff053d4c" />

<img width="40%" alt="Screenshot 2026-08-11 at 10 17 05 PM" src="https://github.com/user-attachments/assets/ac067bc6-f595-4ce7-ae0a-aaec74f5fcbf" />

<img width="40%" alt="Screenshot 2026-08-11 at 10 14 04 PM" src="https://github.com/user-attachments/assets/12472af9-44a7-4906-8a1b-f687dc15569a" />

<img width="40%" alt="Screenshot 2026-07-16 at 2 28 13 PM" src="https://github.com/user-attachments/assets/050a997f-2f40-44ee-b795-a7849084b1de" />
<img width="40%" alt="Screenshot 2026-07-16 at 2 23 50 PM" src="https://github.com/user-attachments/assets/81c73190-8724-4b8a-86a3-3e6448c48864" />

</p>

## Architecture
```
Browser
      │
      ▼
HTMX + JavaScript
      │
      ▼
Go HTTP Server
      │
 ┌─────────────┐
 │ Middleware  │
 │ Handlers    │
 │ Services    │
 │ Database    │
 └─────────────┘
      │
      ▼
 PostgreSQL
```


## Local Development

### Requirements

Go

PostgreSQL

Git

Run Locally

Clone the repository:

git clone https://github.com/alecmstockman/bandplan.git
cd bandplan

Install dependencies:

go mod tidy

Start PostgreSQL and create a local database named bandplan.

Run the server:

go run ./cmd/server

Open the app:

http://localhost:8080

## Project Status

BandPlan is an actively developed MVP focused on building a robust backend architecture before expanding into additional collaboration features. Current development is centered around setlists, cloud storage, notifications, and Progressive Web App support.



## Project Structure
```
.
├── cmd
│   ├── server
│   │   └── main.go
│   ├── testitunes
│   │   └── main.go
│   └── testr2
│       └── main.go
├── Dockerfile
├── go.mod
├── go.sum
├── readme.md
├── sql
│   └── schema
├── src
│   ├── clients
│   │   └── client_itunes.go
│   ├── database
│   │   ├── access_codeDB.go
│   │   ├── band_membersDB.go
│   │   ├── bandsDB.go
│   │   ├── breaksDB.go
│   │   ├── connectDB.go
│   │   ├── messagesDB.go
│   │   ├── sessionsDB.go
│   │   ├── setlist_itemsDB.go
│   │   ├── setlistsDB.go
│   │   ├── songsDB.go
│   │   ├── transitionsDB.go
│   │   └── usersDB.go
│   ├── handlers
│   │   ├── handlers_auth.go
│   │   ├── handlers_breaks.go
│   │   ├── handlers_chats.go
│   │   ├── handlers_menu.go
│   │   ├── handlers_models.go
│   │   ├── handlers_profile.go
│   │   ├── handlers_setlist_partials.go
│   │   ├── handlers_setlists.go
│   │   ├── handlers_songs_itunes.go
│   │   ├── handlers_songs.go
│   │   ├── handlers_transitions.go
│   │   ├── handlers_websocket.go
│   │   ├── helpers_images.go
│   │   ├── helpers_templates.go
│   │   └── helpers.go
│   ├── middleware
│   │   └── middleware.go
│   ├── models
│   │   ├── models_events.go
│   │   ├── models_itunes.go
│   │   ├── models-songs.go
│   │   └── models.go
│   ├── realtime
│   │   ├── client.go
│   │   ├── hub.go
│   │   └── websocket.go
│   ├── services
│   │   └── service_song.go
│   └── storage
│       └── r2.go
├── static
│   ├── css
│   │   ├── core
│   │   │   ├── auth.css
│   │   │   ├── body.css
│   │   │   ├── global.css
│   │   │   ├── navigation.css
│   │   │   └── variables.css
│   │   └── pages
│   │       ├── breaks
│   │       │   └── break-add.css
│   │       ├── calendar.css
│   │       ├── chats
│   │       │   ├── chat.css
│   │       │   └── chats.css
│   │       ├── events.css
│   │       ├── files.css
│   │       ├── goals.css
│   │       ├── home
│   │       │   └── home.css
│   │       ├── legal.css
│   │       ├── profile
│   │       │   ├── admin.css
│   │       │   └── profile.css
│   │       ├── promotion.css
│   │       ├── setlists
│   │       │   ├── setlist-add.css
│   │       │   ├── setlist-edit.css
│   │       │   ├── setlist-reorder.css
│   │       │   ├── setlist.css
│   │       │   └── setlists.css
│   │       ├── songs
│   │       │   ├── lyrics.css
│   │       │   ├── song-download.css
│   │       │   ├── song-edit.css
│   │       │   ├── song.css
│   │       │   ├── songs-add.css
│   │       │   ├── songs-itunes-results.css
│   │       │   └── songs.css
│   │       └── transitions
│   │           └── transition.css
│   ├── images
│   │   └── background.jpg
│   ├── js
│   │   ├── chat-websocket.js
│   │   └── main.js
│   └── uploads
│       └── profile-images
└── templates
    ├── auth
    │   ├── access.html
    │   ├── login.html
    │   ├── privacy.html
    │   ├── register.html
    │   ├── terms.html
    │   └── user-agreement.html
    ├── breaks
    │   ├── break-create.html
    │   ├── break-edit.html
    │   └── break.html
    ├── calendar.html
    ├── chats
    │   ├── chat.html
    │   └── chats.html
    ├── events.html
    ├── files.html
    ├── goals.html
    ├── home
    │   └── index.html
    ├── partials
    │   ├── head.html
    │   ├── layout.html
    │   ├── left_sidebar.html
    │   ├── legal.html
    │   ├── right_sidebar.html
    │   └── svg.html
    ├── profile
    │   ├── admin.html
    │   ├── profile.html
    │   └── settings.html
    ├── promotion.html
    ├── setlists
    │   ├── setlist_reorder.html
    │   ├── setlist_update.html
    │   ├── setlist-add.html
    │   ├── setlist-artwork-preview.html
    │   ├── setlist-edit.html
    │   ├── setlist-partials.html
    │   ├── setlist.html
    │   ├── setlists.html
    │   └── transition-create.html
    ├── songs
    │   ├── lyrics.html
    │   ├── song-download.html
    │   ├── song-edit.html
    │   ├── song.html
    │   ├── songs-add-itunes.html
    │   ├── songs-add.html
    │   ├── songs-itunes-results.html
    │   ├── songs-list.html
    │   └── songs.html
    └── transitions
        └── transition.html
```

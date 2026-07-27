# BandPlan

BandPlan is an all-in-one communication and organization app for bands. It is built around a chat-first workflow, giving musicians one place to manage conversations, songs, setlists, files, events, goals, and band-related planning.

Currently in active development.


![Go](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![HTMX](https://img.shields.io/badge/HTMX-3366CC?style=for-the-badge&logo=htmx&logoColor=white)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-4169E1?style=for-the-badge&logo=postgresql&logoColor=white)
[![Goose](https://img.shields.io/badge/Goose-v3.27.2-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://github.com/pressly/goose)
![HTML5](https://img.shields.io/badge/HTML5-E34F26?style=for-the-badge&logo=html5&logoColor=white)
![CSS3](https://img.shields.io/badge/CSS3-1572B6?style=for-the-badge&logo=css3&logoColor=white)



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
* Import song links with Songlink/Odesli API
* Real-time chat with WebSockets
* Admin privileges for band leaders
* Google Drive integration
* Dynamic setlist builder
* Calendar sync
* Chat-based event creation
* Goal tracking and task management/tracking
* File sharing from R2 cloud storage
* Production deployment with Docker
* Desktop and mobile applications
* Native iOS application



## Feature Breakdown

The 'Songs' feature will contain a list of all yours songs with all the information you need about your immediately available. Things like recording tempo, live tempo, key, lyrics, and song notes are now just a click away for easy reference when you need them. 

The 'Setlists' feature is a dynamic setlist builder that allows you to create multiple custom setlists for quick reference. Each song will have all the information from the 'Songs' feature a click away and allow for custom segments like intros and breaks. All segments will contain a length feature and a total setlist time along with a built in timer so you can time your rehearsals and save the segment times. This will allow users to know exactly how long their set is down to the second. 

The 'Goals' feature is designed around custom band goals you create. Each goal on the list will open and allow you to create sub goals, assign tasks, and create todo lists. 

The 'Calendar' will be a built in calendar designed to sync with any existing personal calendars and offer a single source of truth for all band events. Users will also be able to create calendar events from the chat. 

The 'Files' feature will sync with Google Drive and/or Dropbox to allow access to all your bands files with the organization your used to. Create new files, and folders and manage your drive without having to change apps and easily share to the chat with one click of a button. 




<p align="center">

<img width="90%" alt="Screenshot 2026-07-16 at 2 29 20 PM" src="https://github.com/user-attachments/assets/c6f3ef98-ced7-4207-a0d9-3ea771d33364" />
<img width="90%" alt="Screenshot 2026-07-16 at 2 29 31 PM" src="https://github.com/user-attachments/assets/b18d267c-57b9-4ec0-9b49-21a6ad57bb80" />
<img width="50%" alt="Screenshot 2026-07-16 at 2 30 34 PM" src="https://github.com/user-attachments/assets/68a4a867-4be3-48e6-b6a9-f839a1a2e162" />
<img width="90%" alt="Screenshot 2026-07-16 at 2 23 15 PM" src="https://github.com/user-attachments/assets/664d13d2-b558-4f16-a2d9-230d73aff44f" />
<img width="90%" alt="Screenshot 2026-07-16 at 2 24 14 PM" src="https://github.com/user-attachments/assets/1b3b9449-9b87-4db0-8cca-7ea9fe6fc8b7" />
<img width="90%" alt="Screenshot 2026-07-16 at 2 28 08 PM" src="https://github.com/user-attachments/assets/5239cb31-acd5-457d-ae05-ee2cff053d4c" />
<img width="90%" alt="Screenshot 2026-07-16 at 2 28 13 PM" src="https://github.com/user-attachments/assets/050a997f-2f40-44ee-b795-a7849084b1de" />
<img width="90%" alt="Screenshot 2026-07-16 at 2 23 42 PM" src="https://github.com/user-attachments/assets/5da4788f-a337-42f0-8383-c8d7f39d2488" />
<img width="90%" alt="Screenshot 2026-07-16 at 2 23 50 PM" src="https://github.com/user-attachments/assets/81c73190-8724-4b8a-86a3-3e6448c48864" />

</p>




## Architecture
```
.
├── cmd
│   ├── server
│   │   └── main.go
│   └── testitunes
│       └── main.go
│
├── Dockerfile
├── go.mod
├── go.sum
├── readme.md
├── sql
│   └── schema
│       ├── 001_create_users.sql
│       ├── 002_create_bands.sql
│       ├── 003_create_band_members.sql
│       ├── 004_create_messages.sql
│       ├── 005_create_sessions.sql
│       ├── 006_create_songs.sql
│       ├── 007_create_setlists.sql
│       ├── 008_create_setlist_songs.sql
│       ├── 009_create_access_codes.sql
│       └── 20260725170417_make_setlist_position_deferrable.sql
│
├── src
│   ├── clients
│   │   └── client_itunes.go
│   ├── database
│   │   ├── access_codeDB.go
│   │   ├── band_membersDB.go
│   │   ├── bandsDB.go
│   │   ├── connectDB.go
│   │   ├── messagesDB.go
│   │   ├── sessionsDB.go
│   │   ├── setlistsDB.go
│   │   ├── songsDB.go
│   │   └── usersDB.go
│   ├── handlers
│   │   ├── handlers_auth.go
│   │   ├── handlers_chat.go
│   │   ├── handlers_menu.go
│   │   ├── handlers_models.go
│   │   ├── handlers_profile.go
│   │   ├── handlers_setlists.go
│   │   ├── handlers_songs_itunes.go
│   │   ├── handlers_songs.go
│   │   ├── handlers_websocket.go
│   │   ├── helpers_images.go
│   │   ├── helpers_templates.go
│   │   └── helpers.go
│   ├── middleware
│   │   └── middleware.go
│   ├── models
│   │   ├── models_itunes.go
│   │   ├── models-songs.go
│   │   └── models.go
│   ├── realtime
│   │   ├── client.go
│   │   ├── hub.go
│   │   └── websocket.go
│   └── services
│       └── service_song.go
│
├── static
│   ├── css
│   │   ├── core
│   │   │   ├── auth.css
│   │   │   ├── body.css
│   │   │   ├── global.css
│   │   │   ├── navigation.css
│   │   │   └── variables.css
│   │   └── pages
│   │       ├── calendar.css
│   │       ├── chat.css
│   │       ├── events.css
│   │       ├── files.css
│   │       ├── goals.css
│   │       ├── legal.css
│   │       ├── profile
│   │       │   ├── admin.css
│   │       │   └── profile.css
│   │       ├── promotion.css
│   │       ├── setlists
│   │       │   ├── setlist-add.css
│   │       │   ├── setlist.css
│   │       │   └── setlists.css
│   │       └── songs
│   │           ├── lyrics.css
│   │           ├── song-download.css
│   │           ├── song-edit.css
│   │           ├── song.css
│   │           ├── songs-add.css
│   │           ├── songs-itunes-results.css
│   │           └── songs.css
│   │
│   ├── images
│   │   └── background.jpg
│   ├── js
│   │   ├── chat-websocket.js
│   │   └── main.js
│   └── uploads
│       ├── profile-images
│       └── song-images
│
└── templates
    ├── auth
    │   ├── access.html
    │   ├── login.html
    │   ├── privacy.html
    │   ├── register.html
    │   ├── terms.html
    │   └── user-agreement.html
    ├── calendar.html
    ├── events.html
    ├── files.html
    ├── goals.html
    ├── index.html
    ├── partials
    │   ├── head.html
    │   ├── layout.html
    │   ├── left_sidebar.html
    │   ├── legal.html
    │   ├── right_sidebar.html
    │   └── svg.html
    ├── profile
    │   ├── admin.html
    │   ├── profile-pic.html
    │   └── settings.html
    ├── promotion.html
    ├── setlists
    │   ├── setlist-add.html
    │   ├── setlist.html
    │   └── setlists.html
    └── songs
        ├── lyrics.html
        ├── song-download.html
        ├── song-edit.html
        ├── song.html
        ├── songs-add-itunes.html
        ├── songs-add.html
        ├── songs-itunes-results.html
        ├── songs-list.html
        └── songs.html
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

## Status

BandPlan is currently an MVP-stage project. The main focus is building a strong foundation with authentication, chat, songs, setlists, media uploads, and a clean musician-focused interface.



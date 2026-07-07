# BandPlan

BandPlan is an all-in-one communication and organization app for bands. It is built around a chat-first workflow, giving musicians one place to manage conversations, songs, setlists, files, events, goals, and band-related planning.

Currently in active development.


![Go](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![HTMX](https://img.shields.io/badge/HTMX-3366CC?style=for-the-badge&logo=htmx&logoColor=white)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-4169E1?style=for-the-badge&logo=postgresql&logoColor=white)
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

<img width="1512" height="818" alt="Screenshot 2026-06-28 at 10 45 51 PM" src="https://github.com/user-attachments/assets/9e972c9c-1db6-4923-8743-b334588cee95" />

<img width="1497" height="805" alt="Screenshot 2026-06-28 at 10 47 52 PM" src="https://github.com/user-attachments/assets/5fe4d845-5f00-4158-903d-fdeaece3ada8" />

<img width="1481" height="818" alt="Screenshot 2026-06-28 at 10 48 59 PM" src="https://github.com/user-attachments/assets/2b89b571-bbfe-4789-b83c-63893410239c" />

<img width="1165" height="820" alt="Screenshot 2026-06-28 at 10 49 21 PM" src="https://github.com/user-attachments/assets/6156d006-df9f-4992-b363-f2b0e242a82c" />

<img width="703" height="819" alt="Screenshot 2026-06-28 at 10 49 59 PM" src="https://github.com/user-attachments/assets/08983cca-b828-4848-85a3-1db44625ef09" />

<img width="746" height="824" alt="Screenshot 2026-06-28 at 10 50 20 PM" src="https://github.com/user-attachments/assets/238dbb71-0096-4cd3-92d9-e909a5e10636" />




## Architecture
```
.
├── cmd
│   ├── server
│   │   └── main.go
│   └── testitunes
│       └── main.go
│
├── go.mod
├── go.sum
├── readme.md
├── src
│   ├── clients
│   │   └── client_itunes.go
│   ├── database
│   │   ├── band_membersDB.go
│   │   ├── bandsDB.go
│   │   ├── messagesDB.go
│   │   ├── sessionsDB.go
│   │   ├── setlistsDB.go
│   │   ├── songsDB.go
│   │   └── usersDB.go
│   ├── handlers
│   │   ├── handlers_auth.go
│   │   ├── handlers_chat.go
│   │   ├── handlers_menu.go
│   │   ├── handlers_profile.go
│   │   ├── handlers_setlists.go
│   │   ├── handlers_songs.go
│   │   ├── helpers_images.go
│   │   ├── helpers_templates.go
│   │   └── helpers.go
│   ├── models
│   │   ├── models_itunes.go
│   │   └── models.go
│   └── services
│       └── service_song.go
│
├── static
│   ├── auth.css
│   ├── body.css
│   ├── calendar.css
│   ├── chat.css
│   ├── events.css
│   ├── files.css
│   ├── global.css
│   ├── goals.css
│   ├── images
│   │   └── background.jpg
│   ├── js
│   │   └── main.js
│   ├── navigation.css
│   ├── profile
│   │   └── profile-pic.css
│   ├── promotion.css
│   ├── setlists
│   │   ├── setlist.css
│   │   └── setlists.css
│   ├── songs
│   │   ├── song-download.css
│   │   ├── song-edit.css
│   │   ├── song.css
│   │   ├── songs-add.css
│   │   └── songs.css
│   └── uploads
│       ├── profile-images
│       └── song-images
│
└── templates
    ├── auth
    │   ├── login.html
    │   └── register.html
    ├── calendar.html
    ├── events.html
    ├── files.html
    ├── goals.html
    ├── index.html
    ├── partials
    │   ├── head.html
    │   ├── left_sidebar.html
    │   ├── right_sidebar.html
    │   └── svg.html
    ├── profile
    │   └── profile-pic.html
    ├── promotion.html
    ├── setlists
    │   ├── setlist.html
    │   └── setlists.html
    └── songs
        ├── song-download.html
        ├── song-edit.html
        ├── song.html
        ├── songs-add-itunes.html
        ├── songs-add.html
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



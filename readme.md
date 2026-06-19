# BandPlan

BandPlan aims to be an all in one communication and organization app for bands. The app centers around a core chat app and offers a quick and immediate access to the tools you need to move your band forward all in one app. 

Currently in development

![Go](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![HTMX](https://img.shields.io/badge/HTMX-3366CC?style=for-the-badge&logo=htmx&logoColor=white)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-4169E1?style=for-the-badge&logo=postgresql&logoColor=white)
![HTML5](https://img.shields.io/badge/HTML5-E34F26?style=for-the-badge&logo=html5&logoColor=white)
![CSS3](https://img.shields.io/badge/CSS3-1572B6?style=for-the-badge&logo=css3&logoColor=white)

## Planned Features
* Realtime chat
* User accounts
* Admin priveleges for band leaders
* Google drive connection
* Setlist builder
* Song info list
* Calendar

## Feature Breakdown

The 'Songs' feature will contain a list of all yours songs with all the information you need about your immediately available. Things like recording tempo, live tempo, key, lyrics, and song notes are now just a click away for easy reference when you need them. 

The 'Setlists' feature is a dynamic setlist builder that allows you to create multiple custom setlists for quick reference. Each song will have all the information from the 'Songs' feature a click away and allow for custom segments like intros and breaks. All segments will contain a length feature and a total setlist time along with a built in timer so you can time your rehearsals and save the segment times. This will allow users to know exactly how long their set is down to the second. 

The 'Goals' feature is designed around custom band goals you create. Each goal on the list will open and allow you to create sub goals, assign tasks, and create todo lists. 

The 'Calendar' will be a built in calendar designed to sync with any existing personal calendars and offer a single source of truth for all band events. Users will also be able to create calendar events from the chat. 

The 'Files' feature will sync with Google Drive and/or Dropbox to allow access to all your bands files with the organization your used to. Create new files, and folders and manage your drive without having to change apps and easily share to the chat with one click of a button. 


## Architecture
```
.
├── go.mod
├── go.sum
├── main.go
├── readme.md
├── src
│   ├── database
│   │   ├── band_membersDB.go
│   │   ├── bandsDB.go
│   │   ├── messagesDB.go
│   │   ├── sessionsDB.go
│   │   └── usersDB.go
│   ├── handlers
│   │   ├── auth_handlers.go
│   │   ├── chat_handlers.go
│   │   ├── helpers.go
│   │   └── menu_handlers.go
│   └── models
│       └── models.go
├── static
│   ├── auth.css
│   ├── body.css
│   ├── calendar.css
│   ├── chat.css
│   ├── files.css
│   ├── goals.css
│   ├── images
│   │   └── background.jpg
│   ├── navigation.css
│   ├── setlists.css
│   └── songs.css
└── templates
    ├── calendar.html
    ├── files.html
    ├── goals.html
    ├── index.html
    ├── login.html
    ├── register.html
    ├── setlists.html
    └── songs.html
```
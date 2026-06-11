# BandPlan

All in one organization and communication app for bands

Currently in it's infancy

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
* Calendar


## Architecture
```
.
├── go.mod
├── go.sum
├── main.go
├── readme.md
├── src
│   ├── database
│   │   ├── messagesDB.go
│   │   ├── sessionsDB.go
│   │   └── usersDB.go
│   ├── handlers
│   │   ├── handlers.go
│   │   └── helpers.go
│   └── models
│       └── models.go
├── static
│   └── style.css
└── templates
    ├── index.html
    ├── login.html
    └── register.html
```
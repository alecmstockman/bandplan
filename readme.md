# BandPlan

BandPlan is a full-stack web application that helps bands organize communication, songs, setlists, events, and to-do lists in one place.

Built with **Go, HTMX, PostgreSQL, and WebSockets**, it combines a server-rendered interface with real-time messaging and a responsive design for phones and desktops.

As a musician, I built BandPlan to solve a problem I experienced firsthand: band information scattered across group chats, documents, calendars, and file storage. I built the application from scratch without a full-stack framework to deepen my understanding of backend architecture, authentication, relational data modeling, and production deployment.

**Status:** Early beta. Core features are available, with ongoing improvements guided by feedback from musicians and bands.

[Open BandPlan](https://app.bandplan.app/)


<img width="80%" alt="IMG_8709" src="https://github.com/user-attachments/assets/f5818393-2a98-4b01-a2b5-85cb5fa3954c" />


## Features

### Chat and communication

- Real-time messaging using Gorilla WebSocket.
- Primary band chats and additional group chats.
- Message reactions, and pinned messages.

### Song library

- Store song titles, artwork, lyrics, notes, and streaming links.
- Track recording and live tempos, keys, time signatures, tuning, and duration.
- Import music metadata and artwork through the iTunes Search API.

### Setlists

- Create and dynamically reorder custom setlists.
- Add songs, breaks, and transitions.
- View individual item durations and total set length.
- Link to band events.

### Events and calendar

- Organize shows, rehearsals, and other band events.
- Store venue information, event artwork, and schedules such as load-in, sound check, doors, and set time.
- Associate a setlist with an event.
- View events in an interactive FullCalendar interface.

### To-do lists

- Create personal and band to-do lists.
- Keep primary lists alongside lists for specific projects or sessions.
- Organize tasks and track completion.

### Files, accounts, and bands

- Upload and organize band files.
- Register accounts and verify email addresses.
- Customize user profiles and profile images.
- Join multiple bands and switch between them.
- Manage band access through administrator roles and invitation codes.


</details>

## Tech Stack

| Layer | Technologies |
| --- | --- |
| Backend | Go, `net/http`, `html/template` |
| Database | PostgreSQL, Goose migrations |
| Real-time communication | Gorilla WebSocket |
| Frontend | HTMX, HTML, CSS, JavaScript, FullCalendar |
| Authentication | Session cookies, Argon2id password hashing, Gorilla CSRF |
| Storage | Cloudflare R2 |
| Email | Resend | 
| Infrastructure | Docker, Railway, Cloudflare |
| External data | iTunes Search API |

## Architecture

BandPlan runs as a single Go application. HTTP handlers coordinate requests with services and database functions, while Go templates render pages and HTMX updates individual parts of the interface. A WebSocket hub handles real-time chat communication.

PostgreSQL stores application data, including users, band memberships, sessions, messages, songs, setlists, events, and to-do lists. Cloudflare R2 provides object storage, and Resend handles transactional email.

| Location | Responsibility |
| --- | --- |
| `cmd/server` | Application startup and server configuration |
| `src/handlers` | HTTP handlers and response rendering |
| `src/services` | Application logic |
| `src/database` | PostgreSQL queries and persistence |
| `src/models` | Application data structures |
| `src/middleware` | Authentication and request middleware |
| `src/realtime` | WebSocket connections and message distribution |
| `src/clients` | External service clients |
| `src/storage` | Object storage integration |
| `sql/schema` | Database migrations |
| `templates` | Server-rendered HTML templates |
| `static` | CSS, JavaScript, and image assets |

## Engineering Highlights

- **Band-scoped access:** Membership and authorization checks control access to shared band data.
- **Real-time communication:** A WebSocket hub coordinates connected chat clients, with read and write pumps managing connection traffic.
- **Authentication:** Session-based login, Argon2id password hashing, email verification, and CSRF protection support account security.
- **Relational data modeling:** PostgreSQL relationships and constraints connect bands, members, chats, songs, setlists, events, and tasks.
- **File handling:** Image uploads include validation and processing before storage.
- **Request diagnostics:** Structured logging includes request context to help investigate failures.
- **Production deployment:** A Dockerized Go server runs on Railway with Cloudflare services supporting the deployment.

## Local Development

These steps outline the development workflow. Database connection settings and external-service credentials must match the current application configuration; the examples below cover CSRF configuration only.

### Requirements

- Go compatible with the version declared in `go.mod`.
- PostgreSQL.
- Git.
- Goose for applying database migrations.
- Credentials for any external services used by the features you intend to run.

### Mobile Access

BandPlan works in a browser on phones and desktops. Users can also save it to their device’s home screen for convenient access without downloading it from an app store.

### Setup

1. Clone the repository and download dependencies:

   ```sh
   git clone https://github.com/alecmstockman/bandplan.git
   cd bandplan
   go mod download
   ```

2. Start PostgreSQL, create a local database, configure the application's database connection, and apply the migrations in `sql/schema` using Goose.

3. Configure the required environment variables. For CSRF protection, generate a key:

   ```sh
   openssl rand -hex 32
   ```

   Set these values in the server's environment:

   ```env
   CSRF_AUTH_KEY=replace-with-64-hexadecimal-characters
   CSRF_SECURE=false
   ```

   Keep `CSRF_AUTH_KEY` stable and secret. Use `CSRF_SECURE=false` for local HTTP development and `CSRF_SECURE=true` for HTTPS deployments.

   The application does not automatically load `.env`. If you store configuration there, export it before starting the server:

   ```sh
   set -a
   source .env
   set +a
   ```

4. Start the application:

   ```sh
   go run ./cmd/server
   ```

5. Open [localhost:8080](http://localhost:8080).

## Roadmap

- Notifications and background reminders.
- Calendar integration and synchronization.
- Expanded goal tracking and checklists.
- PDF setlist export.
- Additional file-storage integrations.
- Native iOS application.

Roadmap items are planned or in development and are not all available in the beta.

## What I Learned

Building BandPlan has given me hands-on experience taking an application from an idea to a deployed beta:

- Designing relational schemas and maintaining database migrations.
- Building HTTP handlers, services, and reusable server-rendered templates in Go.
- Implementing authentication and authorization across multiple bands.
- Managing persistent WebSocket connections and real-time message delivery.
- Integrating third-party APIs, transactional email, and object storage.
- Diagnosing issues through structured logs and feedback from early users.
- Balancing feature development, usability improvements, and backend maintenance.

## Feedback

BandPlan is actively being developed. Bug reports and feature suggestions are welcome through [GitHub Issues](https://github.com/alecmstockman/bandplan/issues).



## Project Structure

```text
cmd/server/       Application startup and server configuration
src/
  handlers/       HTTP handlers and response rendering
  services/       Application logic
  database/       PostgreSQL queries and persistence
  middleware/     Authentication and request middleware
  models/         Application data structures
  realtime/       WebSocket connections and message distribution
  clients/        External service clients
  storage/        Object storage integration
  helpers/        Shared utilities
  logging/        Request logging
sql/schema/       Database migrations
templates/        Server-rendered HTML templates
static/           CSS, JavaScript, images, and PWA assets
```

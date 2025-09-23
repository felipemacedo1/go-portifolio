# Portfolio Backend API

A modern Golang backend API for a dynamic portfolio application, built with Gin framework, MongoDB Atlas, and GitHub API integration.

## 🚀 Features

- **RESTful API** built with Gin framework
- **MongoDB Atlas** integration for data persistence
- **GitHub API** integration for dynamic repository and profile data
- **CORS** enabled for cross-origin requests
- **Docker** support for easy deployment
- **Environment-based** configuration

## 📋 API Endpoints

### Profile
- `GET /api/v1/profile` - Get user profile
- `POST /api/v1/profile/sync` - Sync profile from GitHub

### Repositories
- `GET /api/v1/repositories` - Get all repositories
- `POST /api/v1/repositories/sync` - Sync repositories from GitHub

### Skills
- `GET /api/v1/skills` - Get all skills
- `POST /api/v1/skills` - Create new skill
- `PUT /api/v1/skills/:id` - Update skill
- `DELETE /api/v1/skills/:id` - Delete skill

### Projects
- `GET /api/v1/projects` - Get all projects
- `GET /api/v1/projects/featured` - Get featured projects
- `POST /api/v1/projects` - Create new project
- `PUT /api/v1/projects/:id` - Update project
- `DELETE /api/v1/projects/:id` - Delete project

### Experience
- `GET /api/v1/experience` - Get work experience
- `POST /api/v1/experience` - Create new experience
- `PUT /api/v1/experience/:id` - Update experience
- `DELETE /api/v1/experience/:id` - Delete experience

### Health
- `GET /health` - Health check endpoint
- `GET /` - API information

## 🛠️ Setup & Installation

### Prerequisites
- Go 1.21+
- MongoDB Atlas account (free tier)
- GitHub Personal Access Token (optional, for private repos)

### Local Development

1. **Clone the repository**
   ```bash
   git clone https://github.com/felipemacedo1/b.git
   cd b
   ```

2. **Install dependencies**
   ```bash
   go mod tidy
   ```

3. **Set up environment variables**
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

4. **Run the application**
   ```bash
   go run main.go
   ```

The API will be available at `http://localhost:8080`

### Docker Deployment

1. **Build and run with Docker Compose**
   ```bash
   docker-compose up -d
   ```

2. **Or build and run with Docker**
   ```bash
   docker build -t portfolio-backend .
   docker run -p 8080:8080 --env-file .env portfolio-backend
   ```

## ⚙️ Configuration

Create a `.env` file based on `.env.example`:

```env
# MongoDB Configuration
MONGO_URI=mongodb+srv://username:password@cluster.mongodb.net/
DATABASE_NAME=portfolio

# GitHub Configuration
GITHUB_TOKEN=your_github_token_here
GITHUB_USER=felipemacedo1

# Server Configuration
ENVIRONMENT=development
PORT=8080
```

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `MONGO_URI` | MongoDB connection string | `mongodb://localhost:27017` |
| `DATABASE_NAME` | Database name | `portfolio` |
| `GITHUB_TOKEN` | GitHub Personal Access Token | `` |
| `GITHUB_USER` | GitHub username | `felipemacedo1` |
| `ENVIRONMENT` | Application environment | `development` |
| `PORT` | Server port | `8080` |

## 🏗️ Project Structure

```
.
├── main.go                 # Application entry point
├── internal/
│   ├── api/               # API handlers and routes
│   │   ├── routes.go      # Route definitions
│   │   ├── profile.go     # Profile endpoints
│   │   ├── repository.go  # Repository endpoints
│   │   ├── skills.go      # Skills endpoints
│   │   ├── projects.go    # Projects endpoints
│   │   └── experience.go  # Experience endpoints
│   ├── config/            # Configuration management
│   │   └── config.go
│   ├── database/          # Database connection
│   │   └── database.go
│   ├── github/            # GitHub API client
│   │   └── client.go
│   └── models/            # Data models
│       └── models.go
├── Dockerfile             # Docker configuration
├── docker-compose.yml     # Docker Compose configuration
├── .env.example          # Environment variables example
├── .gitignore            # Git ignore rules
└── README.md             # Project documentation
```

## 🔧 Development

### Running Tests
```bash
go test ./...
```

### Code Formatting
```bash
go fmt ./...
```

### Building for Production
```bash
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .
```

## 📝 Usage Examples

### Sync GitHub Data
```bash
# Sync profile
curl -X POST http://localhost:8080/api/v1/profile/sync

# Sync repositories
curl -X POST http://localhost:8080/api/v1/repositories/sync
```

### Create a New Skill
```bash
curl -X POST http://localhost:8080/api/v1/skills \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Go",
    "category": "Programming Languages",
    "level": "advanced"
  }'
```

### Create a New Project
```bash
curl -X POST http://localhost:8080/api/v1/projects \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Portfolio Backend",
    "description": "A modern Golang backend API",
    "technologies": ["Go", "Gin", "MongoDB"],
    "github_url": "https://github.com/felipemacedo1/b",
    "featured": true
  }'
```

## 🚀 Deployment

This backend is designed to work seamlessly with MongoDB Atlas (free tier) and can be deployed to various platforms:

- **Heroku**: Use the included Dockerfile
- **Railway**: Direct git deployment
- **Vercel**: Serverless deployment
- **DigitalOcean**: Docker-based deployment
- **AWS/GCP/Azure**: Container services

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Gin](https://gin-gonic.com/) - HTTP web framework
- [MongoDB Go Driver](https://go.mongodb.org/mongo-driver/) - MongoDB driver for Go
- [GitHub API](https://docs.github.com/en/rest) - GitHub REST API
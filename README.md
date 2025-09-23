# 🚀 Portfolio Backend API

Uma API REST em Golang para servir dados dinâmicos para seu portfólio, integrada com a GitHub API e MongoDB Atlas.

## 📋 Índice

- [Características](#-características)
- [Tecnologias](#-tecnologias)
- [Início Rápido](#-início-rápido)
- [Configuração](#-configuração)
- [Endpoints da API](#-endpoints-da-api)
- [Autenticação](#-autenticação)
- [Deploy](#-deploy)
- [Desenvolvimento](#-desenvolvimento)
- [Exemplos](#-exemplos)
- [Contribuição](#-contribuição)

## ✨ Características

- **🔄 Integração GitHub**: Sincronização automática com perfil, repositórios e contribuições
- **📝 Gestão de Conteúdo**: CRUD completo para skills, experiência, projetos e educação
- **⚡ Cache Inteligente**: Sistema de cache com TTL configurável e cleanup automático
- **🔐 Autenticação**: JWT e API tokens para operações protegidas
- **🛡️ Rate Limiting**: Proteção contra abuse com limites por IP
- **📊 Analytics**: Métricas detalhadas de performance e uso
- **🌐 CORS**: Configurado para integração com GitHub Pages
- **🐳 Docker**: Containerização para deploy simplificado
- **📖 Logs Estruturados**: Sistema de logging completo com request ID

## 🛠 Tecnologias

- **Linguagem**: Go 1.21+
- **Framework**: Gin (HTTP framework)
- **Banco de Dados**: MongoDB Atlas
- **Cache**: Sistema interno com MongoDB
- **Autenticação**: JWT + Bearer tokens
- **Deploy**: Docker + Railway/Render/Vercel

## 🚀 Início Rápido

### Pré-requisitos

- Go 1.21+
- MongoDB Atlas (plano gratuito)
- Token do GitHub (opcional, mas recomendado)

### Instalação Local

1. **Clone o repositório**
```bash
git clone https://github.com/felipemacedo1/portfolio-backend.git
cd portfolio-backend
```

2. **Configure as variáveis de ambiente**
```bash
cp .env.example .env
# Edite o arquivo .env com suas configurações
```

3. **Instale as dependências**
```bash
go mod tidy
```

4. **Execute a aplicação**
```bash
go run main.go
```

A API estará disponível em `http://localhost:8080`

### Docker

```bash
# Build da imagem
docker build -t portfolio-backend .

# Execute o container
docker run -p 8080:8080 --env-file .env portfolio-backend
```

## ⚙️ Configuração

### Variáveis de Ambiente

Crie um arquivo `.env` baseado no `.env.example`:

```bash
# Database
MONGODB_URI=mongodb+srv://user:pass@cluster.mongodb.net/portfolio
DATABASE_NAME=portfolio

# GitHub API
GITHUB_TOKEN=ghp_your_personal_access_token
GITHUB_USERNAME=felipemacedo1

# Server Config
PORT=8080
GIN_MODE=release
CORS_ORIGINS=https://felipemacedo1.github.io,http://localhost:3000

# Auth
JWT_SECRET=your_super_secret_key
API_TOKEN=bearer_token_for_write_operations

# Cache & Performance
GITHUB_CACHE_TTL=6h
CONTENT_CACHE_TTL=24h
RATE_LIMIT_REQUESTS=100
RATE_LIMIT_WINDOW=3600s

# Monitoring
LOG_LEVEL=info
ENABLE_METRICS=true
```

### MongoDB Atlas Setup

1. Crie uma conta gratuita no [MongoDB Atlas](https://cloud.mongodb.com/)
2. Crie um cluster (M0 - gratuito)
3. Configure o acesso de rede (0.0.0.0/0 para desenvolvimento)
4. Crie um usuário de banco de dados
5. Obtenha a string de conexão

### GitHub Token

Para evitar rate limiting da GitHub API:

1. Acesse [GitHub Settings > Developer settings > Personal access tokens](https://github.com/settings/tokens)
2. Gere um novo token (classic)
3. Selecione os escopos: `public_repo`, `read:user`
4. Adicione o token ao arquivo `.env`

## 📡 Endpoints da API

### Health & Info

```http
GET /health                    # Health check
GET /readiness                 # Readiness probe (Kubernetes)
GET /liveness                  # Liveness probe (Kubernetes)
GET /api/v1/info              # Informações da API
```

### Content Management

```http
GET /api/v1/content           # Todo conteúdo do portfólio
GET /api/v1/content/skills    # Skills técnicas
GET /api/v1/content/experience # Experiência profissional
GET /api/v1/content/projects  # Projetos desenvolvidos
GET /api/v1/content/education # Formação acadêmica
GET /api/v1/content/meta      # Informações pessoais
GET /api/v1/content/search?q=query # Busca no conteúdo

# Endpoints protegidos (requer autenticação)
PUT /api/v1/content           # Atualizar conteúdo
GET /api/v1/content/history/:type # Histórico de versões
```

### GitHub Integration

```http
GET /api/v1/github/profile/:username      # Perfil GitHub
GET /api/v1/github/repos/:username        # Repositórios públicos
GET /api/v1/github/contributions/:username # Gráfico de contribuições
GET /api/v1/github/stats/:username        # Estatísticas agregadas
GET /api/v1/github/rate-limit             # Status do rate limit

# Endpoints protegidos
POST /api/v1/github/sync/:username        # Sincronizar dados
```

### Analytics

```http
GET /api/v1/analytics/summary             # Resumo geral
GET /api/v1/analytics/contributions/:period # Contribuições por período
GET /api/v1/analytics/cache-stats         # Estatísticas do cache
GET /api/v1/analytics/performance         # Métricas de performance
```

### Admin (Requer API Key)

```http
POST /api/v1/admin/cache/clear            # Limpar cache
GET /api/v1/admin/system/stats            # Estatísticas do sistema
POST /api/v1/admin/content/import         # Importar conteúdo
```

## 🔐 Autenticação

### Bearer Token (Operações de Escrita)

```bash
curl -H "Authorization: Bearer YOUR_API_TOKEN" \
     -X PUT https://api.example.com/api/v1/content
```

### API Key (Admin)

```bash
curl -H "X-API-Key: YOUR_API_KEY" \
     -X POST https://api.example.com/api/v1/admin/cache/clear
```

## 🚢 Deploy

### Railway.app

1. Conecte seu repositório GitHub
2. Configure as variáveis de ambiente
3. Deploy automático!

```bash
# CLI do Railway
npm install -g @railway/cli
railway login
railway init
railway add
railway deploy
```

### Render.com

1. Conecte seu repositório GitHub
2. Configure como Web Service
3. Adicione as variáveis de ambiente
4. Deploy automático!

### Dockerfile Multi-stage

O projeto inclui um Dockerfile otimizado:

```dockerfile
# Build da imagem
docker build \
  --build-arg VERSION=1.0.0 \
  --build-arg BUILD_TIME=$(date -u +"%Y-%m-%dT%H:%M:%SZ") \
  --build-arg GIT_COMMIT=$(git rev-parse HEAD) \
  -t portfolio-backend .
```

## 💻 Desenvolvimento

### Estrutura do Projeto

```
portfolio-backend/
├── main.go                    # Entry point
├── .env.example              # Template de variáveis
├── Dockerfile               # Container para deploy
├── README.md               # Documentação
├── go.mod & go.sum         # Dependências Go
├── config/
│   └── config.go           # Configurações
├── controllers/
│   ├── health.go           # Health check
│   ├── content.go          # Gestão de conteúdo
│   ├── github.go           # Integração GitHub
│   └── analytics.go        # Analytics
├── services/
│   ├── github_service.go   # Client GitHub API
│   ├── content_service.go  # Business logic
│   └── cache_service.go    # Sistema de cache
├── models/
│   ├── content.go          # Modelos de conteúdo
│   ├── github.go           # Modelos GitHub
│   └── response.go         # Responses padronizados
├── database/
│   └── mongodb.go          # Conexão MongoDB
├── middleware/
│   ├── cors.go             # CORS
│   ├── auth.go             # Autenticação
│   ├── rate_limit.go       # Rate limiting
│   └── logger.go           # Logging
├── routes/
│   └── routes.go           # Definição de rotas
└── utils/
    ├── helpers.go          # Funções utilitárias
    └── validator.go        # Validações
```

### Scripts de Desenvolvimento

```bash
# Executar em modo desenvolvimento
go run main.go

# Build para produção
go build -o portfolio-backend .

# Executar testes
go test ./...

# Verificar dependências
go mod tidy

# Formatar código
go fmt ./...

# Análise estática
go vet ./...
```

### Variáveis de Build

```bash
# Build com informações de versão
go build -ldflags="-X main.version=1.0.0 -X main.buildTime=$(date -u +%Y-%m-%dT%H:%M:%SZ) -X main.gitCommit=$(git rev-parse HEAD)" -o portfolio-backend .
```

## 📝 Exemplos

### Buscar Perfil GitHub

```bash
curl https://api.example.com/api/v1/github/profile/felipemacedo1
```

Resposta:
```json
{
  "success": true,
  "data": {
    "login": "felipemacedo1",
    "name": "Felipe Macedo",
    "avatar_url": "https://avatars.githubusercontent.com/u/...",
    "bio": "Desenvolvedor Full Stack",
    "public_repos": 25,
    "followers": 50,
    "following": 30
  },
  "timestamp": "2025-01-01T12:00:00Z"
}
```

### Atualizar Conteúdo

```bash
curl -X PUT https://api.example.com/api/v1/content \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "type": "skills",
    "data": {
      "backend": [
        {"name": "Go", "level": 85},
        {"name": "Java", "level": 90}
      ]
    }
  }'
```

### Buscar Analytics

```bash
curl https://api.example.com/api/v1/analytics/summary
```

Resposta:
```json
{
  "success": true,
  "data": {
    "summary": {
      "total_repositories": 25,
      "total_stars": 150,
      "total_forks": 45,
      "contribution_streak": 30
    },
    "github": {
      "top_languages": [
        {"name": "Go", "percentage": 45.5},
        {"name": "JavaScript", "percentage": 30.2}
      ]
    }
  }
}
```

## 🔄 Integração com Frontend

### JavaScript/React

```javascript
const API_BASE = 'https://your-api.railway.app/api/v1';

// Buscar dados do portfólio
async function fetchPortfolio() {
  const response = await fetch(`${API_BASE}/content`);
  const data = await response.json();
  return data.data;
}

// Buscar dados do GitHub
async function fetchGitHubStats(username) {
  const response = await fetch(`${API_BASE}/github/stats/${username}`);
  const data = await response.json();
  return data.data;
}
```

### Tratamento de Erros

```javascript
async function apiCall(endpoint) {
  try {
    const response = await fetch(`${API_BASE}${endpoint}`);
    const data = await response.json();
    
    if (!data.success) {
      throw new Error(data.error);
    }
    
    return data.data;
  } catch (error) {
    console.error('API Error:', error.message);
    // Fallback para dados locais
    return getLocalData();
  }
}
```

## 📊 Monitoramento

### Health Checks

A API fornece endpoints para monitoramento:

- `/health` - Status geral da aplicação
- `/readiness` - Pronto para receber tráfego
- `/liveness` - Aplicação está viva

### Métricas

- Response time médio
- Taxa de erro
- Hit rate do cache
- Uso de recursos

### Logs

Logs estruturados em JSON incluem:

- Request ID único
- Informações do usuário
- Tempo de resposta
- Status code
- Detalhes de erro

## 🤝 Contribuição

1. Fork o projeto
2. Crie uma branch para sua feature (`git checkout -b feature/AmazingFeature`)
3. Commit suas mudanças (`git commit -m 'Add some AmazingFeature'`)
4. Push para a branch (`git push origin feature/AmazingFeature`)
5. Abra um Pull Request

## 📄 Licença

Este projeto está sob a licença MIT. Veja o arquivo [LICENSE](LICENSE) para mais detalhes.

## 🆘 Suporte

Se você encontrar algum problema ou tiver dúvidas:

1. Verifique a [documentação](#-endpoints-da-api)
2. Consulte os [exemplos](#-exemplos)
3. Abra uma [issue](https://github.com/felipemacedo1/portfolio-backend/issues)

## 🙏 Agradecimentos

- [Gin Framework](https://gin-gonic.com/)
- [MongoDB](https://www.mongodb.com/)
- [GitHub API](https://docs.github.com/en/rest)
- [Railway](https://railway.app/) para hosting gratuito

---

**Portfolio Backend API** - Transforme seu portfólio estático em uma experiência dinâmica! 🚀
// MongoDB initialization script
// This script runs when the MongoDB container starts for the first time

// Switch to the portfolio database
db = db.getSiblingDB('portfolio');

// Create collections with validation
db.createCollection('contents', {
    validator: {
        $jsonSchema: {
            bsonType: 'object',
            required: ['id', 'type', 'title', 'created_at'],
            properties: {
                id: {
                    bsonType: 'string',
                    description: 'Unique identifier for the content'
                },
                type: {
                    bsonType: 'string',
                    enum: ['about', 'project', 'skill', 'experience', 'education', 'contact'],
                    description: 'Type of content'
                },
                title: {
                    bsonType: 'string',
                    description: 'Title of the content'
                },
                description: {
                    bsonType: 'string',
                    description: 'Description of the content'
                },
                created_at: {
                    bsonType: 'date',
                    description: 'Creation timestamp'
                },
                updated_at: {
                    bsonType: 'date',
                    description: 'Last update timestamp'
                }
            }
        }
    }
});

db.createCollection('github_cache', {
    validator: {
        $jsonSchema: {
            bsonType: 'object',
            required: ['key', 'data', 'expires_at'],
            properties: {
                key: {
                    bsonType: 'string',
                    description: 'Cache key'
                },
                data: {
                    bsonType: 'object',
                    description: 'Cached data'
                },
                expires_at: {
                    bsonType: 'date',
                    description: 'Cache expiration time'
                }
            }
        }
    }
});

db.createCollection('analytics', {
    validator: {
        $jsonSchema: {
            bsonType: 'object',
            required: ['endpoint', 'timestamp'],
            properties: {
                endpoint: {
                    bsonType: 'string',
                    description: 'API endpoint accessed'
                },
                method: {
                    bsonType: 'string',
                    description: 'HTTP method used'
                },
                ip: {
                    bsonType: 'string',
                    description: 'Client IP address'
                },
                user_agent: {
                    bsonType: 'string',
                    description: 'Client user agent'
                },
                response_time: {
                    bsonType: 'number',
                    description: 'Response time in milliseconds'
                },
                status_code: {
                    bsonType: 'number',
                    description: 'HTTP status code'
                },
                timestamp: {
                    bsonType: 'date',
                    description: 'Request timestamp'
                }
            }
        }
    }
});

// Create indexes for better performance
db.contents.createIndex({ 'id': 1 }, { unique: true });
db.contents.createIndex({ 'type': 1 });
db.contents.createIndex({ 'created_at': -1 });
db.contents.createIndex({ 'updated_at': -1 });

db.github_cache.createIndex({ 'key': 1 }, { unique: true });
db.github_cache.createIndex({ 'expires_at': 1 }, { expireAfterSeconds: 0 });

db.analytics.createIndex({ 'endpoint': 1 });
db.analytics.createIndex({ 'timestamp': -1 });
db.analytics.createIndex({ 'ip': 1 });

// Insert sample data
db.contents.insertMany([
    {
        id: 'about-me',
        type: 'about',
        title: 'Sobre Mim',
        description: 'Desenvolvedor Full Stack apaixonado por tecnologia e inovação.',
        content: {
            bio: 'Sou um desenvolvedor experiente com foco em Go, JavaScript e tecnologias modernas.',
            skills: ['Go', 'JavaScript', 'MongoDB', 'Docker', 'Kubernetes'],
            location: 'Brasil',
            email: 'contato@portfolio.com'
        },
        metadata: {
            featured: true,
            order: 1
        },
        created_at: new Date(),
        updated_at: new Date()
    },
    {
        id: 'portfolio-backend',
        type: 'project',
        title: 'Portfolio Backend API',
        description: 'API REST completa para portfolio interativo com integração GitHub.',
        content: {
            technologies: ['Go', 'Gin', 'MongoDB', 'Docker', 'GitHub API'],
            github_url: 'https://github.com/username/portfolio-backend',
            live_url: 'https://api.portfolio.com',
            features: [
                'Integração com GitHub API',
                'Cache inteligente',
                'Rate limiting',
                'Autenticação JWT',
                'Documentação completa'
            ]
        },
        metadata: {
            featured: true,
            order: 1,
            status: 'completed'
        },
        created_at: new Date(),
        updated_at: new Date()
    },
    {
        id: 'go-lang',
        type: 'skill',
        title: 'Go (Golang)',
        description: 'Linguagem de programação moderna e eficiente.',
        content: {
            level: 'advanced',
            years_experience: 3,
            projects_count: 15,
            certifications: []
        },
        metadata: {
            category: 'backend',
            order: 1
        },
        created_at: new Date(),
        updated_at: new Date()
    },
    {
        id: 'senior-developer',
        type: 'experience',
        title: 'Desenvolvedor Senior',
        description: 'Desenvolvimento de aplicações enterprise com Go e microserviços.',
        content: {
            company: 'Tech Company',
            position: 'Senior Go Developer',
            start_date: '2022-01-01',
            end_date: null,
            current: true,
            responsibilities: [
                'Desenvolvimento de APIs REST',
                'Arquitetura de microserviços',
                'Mentoria de desenvolvedores junior',
                'Code review e padrões de qualidade'
            ],
            technologies: ['Go', 'MongoDB', 'Docker', 'Kubernetes', 'gRPC']
        },
        metadata: {
            order: 1
        },
        created_at: new Date(),
        updated_at: new Date()
    }
]);

print('✅ MongoDB initialization completed successfully!');
print('📊 Collections created: contents, github_cache, analytics');
print('📝 Indexes created for optimal performance');
print('🎯 Sample data inserted');
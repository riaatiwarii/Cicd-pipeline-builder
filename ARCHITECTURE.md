# Architecture Overview

## High-Level Architecture

```
┌─────────────────────────────────────────────────────────────────────┐
│                           Users                                     │
└────────────┬────────────────────────────────────────────────────┬───┘
             │                                                    │
             ▼                                                    ▼
    ┌────────────────┐                                  ┌─────────────────┐
    │   React Web    │                                  │  Mobile/CLI     │
    │   Dashboard    │                                  │   Clients       │
    └────────┬────────┘                                  └────────┬────────┘
             │                                                   │
             └───────────────────────┬─────────────────────────┘
                                     │
                                     ▼
                        ┌────────────────────────┐
                        │  API Gateway/LB        │
                        │  (Nginx/ALB)           │
                        └───────────┬────────────┘
                                    │
            ┌───────────────────────┼───────────────────────┐
            │                       │                       │
            ▼                       ▼                       ▼
    ┌──────────────┐      ┌──────────────┐      ┌──────────────┐
    │ API Server   │      │ API Server   │      │ API Server   │
    │ Instance 1   │      │ Instance 2   │      │ Instance N   │
    └──────┬───────┘      └──────┬───────┘      └──────┬───────┘
           │                     │                     │
           └─────────────────────┼─────────────────────┘
                                 │
           ┌─────────────────────┼─────────────────────┐
           │                     │                     │
           ▼                     ▼                     ▼
    ┌──────────────┐    ┌──────────────┐    ┌──────────────┐
    │ PostgreSQL   │    │ Redis Queue  │    │ Message Bus  │
    │ (Primary)    │    │              │    │ (Event Bus)  │
    └──────────────┘    └──────────────┘    └──────────────┘
           ▲
           │
    ┌──────┴───────┐
    │              │
    ▼              ▼
 Replica      Backup
    
    ┌──────────────────────────────────────────────────┐
    │        Job Processing Layer                      │
    └──────────────────────────────────────────────────┘
    │                    │                    │
    ▼                    ▼                    ▼
┌─────────────┐  ┌─────────────┐  ┌─────────────┐
│  Processor  │  │  Processor  │  │  Processor  │
│  Worker 1   │  │  Worker 2   │  │  Worker N   │
└─────────────┘  └─────────────┘  └─────────────┘
    │                │                │
    └────────────────┼────────────────┘
                     │
                     ▼
            ┌────────────────────┐
            │  Docker Daemon     │
            │  Container Runtime │
            └────────────────────┘
                     │
                     ▼
            ┌────────────────────┐
            │   Build Jobs       │
            │  (Ephemeral        │
            │   Containers)      │
            └────────────────────┘
```

## Component Details

### 1. API Layer (Go)

**Purpose**: RESTful API for pipeline management, builds, webhooks

**Key Responsibilities**:
- User authentication & authorization (JWT)
- Pipeline CRUD operations
- Build triggering and status tracking
- Webhook reception and routing
- Real-time log streaming (WebSocket)
- Configuration management

**Technology**:
- Framework: Gin or Echo
- Database: GORM with PostgreSQL
- Queue: Redis Streams or AMQP
- Authentication: JWT

**Scalability**:
- Stateless design enables horizontal scaling
- Load balanced via Nginx/ALB
- Session data stored in database

### 2. Database Layer (PostgreSQL)

**Purpose**: Persistent storage for all system state

**Key Tables**:
- `users` - User accounts & roles
- `pipelines` - Pipeline definitions
- `stages` - Pipeline stages/steps
- `builds` - Pipeline executions
- `jobs` - Individual stage runs
- `logs` - Streaming logs
- `webhooks` - Webhook configurations
- `triggers` - Scheduled/manual triggers
- `artifacts` - Build outputs
- `credentials` - Encrypted secrets

**ACID Compliance**: Ensures reliable pipeline state

**Replication**: Primary-replica for HA

### 3. Message Queue (Redis/RabbitMQ)

**Purpose**: Decoupling API from job processors

**Functions**:
- Job queue for processors
- Pub/Sub for real-time updates
- Cache for frequently accessed data
- Session storage (optional)

**Why Separate**:
- Enables independent scaling of API and processors
- Fault tolerance (jobs survive API restarts)
- Load distribution

### 4. Job Processor (Go Workers)

**Purpose**: Execute pipeline stages in isolated containers

**Key Responsibilities**:
- Poll job queue
- Clone git repositories
- Pull Docker images
- Execute scripts in containers
- Capture and stream logs
- Update job status
- Handle timeouts and failures

**Execution Flow**:
```
Get Job from Queue
    ↓
Clone Repository
    ↓
Parse Pipeline Config
    ↓
For Each Stage:
    - Create Docker container
    - Execute script
    - Stream logs to DB/WebSocket
    - Capture exit code
    ↓
Update Build Status
    ↓
Publish Completion Event
    ↓
Clean up Resources
```

**Scalability**:
- Multiple worker instances
- Auto-scaling based on queue depth
- Parallel job execution

### 5. Frontend (React)

**Purpose**: Web dashboard for pipeline management

**Key Features**:
- Pipeline list and creation
- Build history and logs
- Real-time log streaming
- Trigger management
- Webhook configuration
- Build artifacts download

**Technology**:
- Framework: React 18
- State Management: Zustand
- UI Components: Material-UI
- Build Tool: Vite
- Real-time: WebSocket/SSE

## Data Flow

### Pipeline Creation Flow

```
User Input (UI)
    ↓
POST /api/v1/pipelines
    ↓
API Server (Validate)
    ↓
Database (Store Pipeline)
    ↓
Response to Client
    ↓
UI Updates
```

### Build Execution Flow

```
Webhook or Manual Trigger
    ↓
API Validates Request
    ↓
Create Build Record (DB)
    ↓
Enqueue Job to Redis
    ↓
Respond to Trigger (202 Accepted)
    ↓
Processor Picks Up Job
    ↓
Clone Repo
    ↓
Parse Config
    ↓
For Each Stage:
    - Create Container
    - Execute
    - Stream Logs
    - Capture Result
    ↓
Update Build Status (Success/Failed)
    ↓
Publish Event
    ↓
Frontend Updates (WebSocket)
```

## Security Architecture

### Authentication & Authorization

```
User Login
    ↓
Validate Credentials
    ↓
Generate JWT Token
    ↓
Token Includes: user_id, role, exp
    ↓
All API Requests: Bearer <Token>
    ↓
JWT Verification Middleware
    ↓
Extract user context
    ↓
Check Role-Based Access
```

### Secrets Management

```
User Input (GitHub Token, etc.)
    ↓
Encryption (AES-256)
    ↓
Store in Database
    ↓
On Use: Decrypt
    ↓
Pass to Container
    ↓
Container Variables (env)
```

### Container Isolation

```
Each Job
    ↓
Separate Docker Container
    ↓
Resource Limits (CPU, Memory)
    ↓
Network Isolation
    ↓
Filesystem Isolation
    ↓
Clean-up on Exit
```

## Deployment Models

### Development

```
docker-compose up
- Single machine
- All services in one network
- Suitable for testing
```

### Production (Kubernetes)

```
API Replicas: 2-10 (HPA)
Processors: 3-20 (HPA based on queue)
Database: Managed PostgreSQL
Redis: Managed Redis
Load Balancer: Ingress/ALB
```

### Hybrid

```
On-Premise:
- API Servers
- Job Processors
- Private Git Repos

Cloud:
- Managed Database
- Managed Redis
- CDN for Frontend
- Logging/Monitoring
```

## Performance Optimization

### Caching

```
Pipeline Config: Cache in Redis (invalidate on update)
Recent Logs: Temporary cache in memory
User Permissions: Cache for 5 minutes
Build Status: Real-time via WebSocket
```

### Database

```
Indexes:
- builds.pipeline_id (for queries)
- jobs.build_id (for filtering)
- logs.job_id, timestamp (for retrieval)

Replication:
- Read queries to replica
- Write queries to primary
- Automatic failover
```

### API

```
Rate Limiting: 1000 req/hour per user
Connection Pooling: 20-50 connections
Request Timeout: 30 seconds
Response Compression: Gzip
```

### Jobs

```
Parallel Execution: Multiple processors
Container Reuse: Image layer caching
Artifact Storage: CDN/S3
Log Rotation: Archive old logs
```

## Monitoring & Observability

### Metrics

```
Prometheus Endpoint: /metrics

Key Metrics:
- API request rate, latency, errors
- Job success rate, execution time
- Queue depth, processing time
- Database connections, queries
- Memory, CPU usage
```

### Logging

```
Structured Logs:
- Timestamp
- Service (api, processor)
- Level (info, warn, error)
- Message
- Context (user_id, build_id)

Aggregation: ELK Stack, Loki
Retention: 30 days
```

### Alerting

```
Conditions:
- API error rate > 5%
- Build processor failure rate > 10%
- Queue depth > 1000
- Database replication lag > 1 second
- Disk usage > 80%
```

## High Availability

### API Server

```
Load Balancer (Nginx/ALB)
    ↓
Multiple API Instances
    ↓
Health Checks (every 10s)
    ↓
Auto-replacement of failed instances
```

### Database

```
Primary PostgreSQL
    ↓
Synchronous Replica
    ↓
Read replicas for queries
    ↓
Automated failover on primary failure
```

### Message Queue

```
Redis Cluster or RabbitMQ cluster
    ↓
Persistent queue storage
    ↓
Automatic failover
```

## Disaster Recovery

### Backup Strategy

```
Database:
- Full backups: Daily
- Incremental: Hourly
- Point-in-time recovery: 30 days

Artifacts:
- Store in S3/GCS
- Lifecycle: Archive after 90 days

Configuration:
- Git repository
- IaC (Terraform/CloudFormation)
```

### Recovery Plan

```
RTO (Recovery Time Objective): 15 minutes
RPO (Recovery Point Objective): 1 hour

Steps:
1. Detect failure
2. Promote read replica to primary
3. Restore from latest backup if needed
4. Verify data integrity
5. Resume operations
```

## Future Enhancements

```
Pipeline Orchestration:
- DAG (Directed Acyclic Graph) support
- Conditional stages
- Parallel execution
- Matrix builds

Advanced Features:
- Pipeline templates
- Shared libraries
- Artifact caching
- Build cache layers
- Distributed tracing
```

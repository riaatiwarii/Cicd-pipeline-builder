# CI/CD Pipeline Builder - Kubernetes Deployment

This guide covers deploying the CI/CD Pipeline Builder to Kubernetes for production use.

## Prerequisites

- Kubernetes cluster (1.20+)
- kubectl configured
- Helm 3.0+ (optional but recommended)
- PostgreSQL 15+ instance
- Redis 7+ instance

## Deployment Steps

### 1. Create Namespaces

```bash
kubectl create namespace cicd
kubectl create namespace cicd-runners
```

### 2. Configure Secrets

```bash
# Database credentials
kubectl create secret generic cicd-db-secret \
  --from-literal=username=cicd \
  --from-literal=password=$(openssl rand -base64 32) \
  -n cicd

# JWT Secret
kubectl create secret generic cicd-jwt-secret \
  --from-literal=jwt-secret=$(openssl rand -base64 32) \
  -n cicd

# Docker Registry credentials
kubectl create secret docker-registry cicd-docker-secret \
  --docker-server=registry.io \
  --docker-username=your_username \
  --docker-password=your_password \
  -n cicd
```

### 3. Deploy PostgreSQL (if not using managed service)

```bash
kubectl apply -f - <<EOF
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: postgres-pvc
  namespace: cicd
spec:
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 50Gi
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: postgres
  namespace: cicd
spec:
  replicas: 1
  selector:
    matchLabels:
      app: postgres
  template:
    metadata:
      labels:
        app: postgres
    spec:
      containers:
      - name: postgres
        image: postgres:15-alpine
        ports:
        - containerPort: 5432
        env:
        - name: POSTGRES_USER
          valueFrom:
            secretKeyRef:
              name: cicd-db-secret
              key: username
        - name: POSTGRES_PASSWORD
          valueFrom:
            secretKeyRef:
              name: cicd-db-secret
              key: password
        - name: POSTGRES_DB
          value: cicd_db
        volumeMounts:
        - name: postgres-storage
          mountPath: /var/lib/postgresql/data
      volumes:
      - name: postgres-storage
        persistentVolumeClaim:
          claimName: postgres-pvc
---
apiVersion: v1
kind: Service
metadata:
  name: postgres
  namespace: cicd
spec:
  selector:
    app: postgres
  ports:
  - port: 5432
    targetPort: 5432
  type: ClusterIP
EOF
```

### 4. Deploy Redis (if not using managed service)

```bash
kubectl apply -f - <<EOF
apiVersion: apps/v1
kind: Deployment
metadata:
  name: redis
  namespace: cicd
spec:
  replicas: 1
  selector:
    matchLabels:
      app: redis
  template:
    metadata:
      labels:
        app: redis
    spec:
      containers:
      - name: redis
        image: redis:7-alpine
        ports:
        - containerPort: 6379
        command:
        - redis-server
        - "--appendonly"
        - "yes"
---
apiVersion: v1
kind: Service
metadata:
  name: redis
  namespace: cicd
spec:
  selector:
    app: redis
  ports:
  - port: 6379
    targetPort: 6379
  type: ClusterIP
EOF
```

### 5. Deploy API Server

```bash
kubectl apply -f - <<EOF
apiVersion: apps/v1
kind: Deployment
metadata:
  name: api
  namespace: cicd
spec:
  replicas: 3
  selector:
    matchLabels:
      app: api
  template:
    metadata:
      labels:
        app: api
    spec:
      containers:
      - name: api
        image: cicd-api:latest
        imagePullPolicy: Always
        ports:
        - containerPort: 8080
        env:
        - name: DATABASE_URL
          value: postgres://postgres:postgres@postgres:5432/cicd_db
        - name: REDIS_URL
          value: redis://redis:6379
        - name: JWT_SECRET
          valueFrom:
            secretKeyRef:
              name: cicd-jwt-secret
              key: jwt-secret
        - name: PORT
          value: "8080"
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 10
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /ready
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
      imagePullSecrets:
      - name: cicd-docker-secret
---
apiVersion: v1
kind: Service
metadata:
  name: api
  namespace: cicd
spec:
  selector:
    app: api
  ports:
  - port: 8080
    targetPort: 8080
  type: LoadBalancer
---
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: api-hpa
  namespace: cicd
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: api
  minReplicas: 2
  maxReplicas: 10
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 70
  - type: Resource
    resource:
      name: memory
      target:
        type: Utilization
        averageUtilization: 80
EOF
```

### 6. Deploy Job Processors

```bash
kubectl apply -f - <<EOF
apiVersion: apps/v1
kind: Deployment
metadata:
  name: processor
  namespace: cicd-runners
spec:
  replicas: 5
  selector:
    matchLabels:
      app: processor
  template:
    metadata:
      labels:
        app: processor
    spec:
      containers:
      - name: processor
        image: cicd-processor:latest
        imagePullPolicy: Always
        env:
        - name: DATABASE_URL
          value: postgres://postgres:postgres@postgres.cicd:5432/cicd_db
        - name: REDIS_URL
          value: redis://redis.cicd:6379
        - name: WORKSPACE_PATH
          value: /workspace
        volumeMounts:
        - name: docker-socket
          mountPath: /var/run/docker.sock
        - name: workspace
          mountPath: /workspace
      volumes:
      - name: docker-socket
        hostPath:
          path: /var/run/docker.sock
      - name: workspace
        emptyDir: {}
      imagePullSecrets:
      - name: cicd-docker-secret
---
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: processor-hpa
  namespace: cicd-runners
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: processor
  minReplicas: 3
  maxReplicas: 20
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 60
EOF
```

### 7. Deploy Frontend

```bash
kubectl apply -f - <<EOF
apiVersion: apps/v1
kind: Deployment
metadata:
  name: frontend
  namespace: cicd
spec:
  replicas: 2
  selector:
    matchLabels:
      app: frontend
  template:
    metadata:
      labels:
        app: frontend
    spec:
      containers:
      - name: frontend
        image: cicd-frontend:latest
        imagePullPolicy: Always
        ports:
        - containerPort: 3000
        env:
        - name: REACT_APP_API_URL
          value: "https://api.yourdomain.com"
      imagePullSecrets:
      - name: cicd-docker-secret
---
apiVersion: v1
kind: Service
metadata:
  name: frontend
  namespace: cicd
spec:
  selector:
    app: frontend
  ports:
  - port: 3000
    targetPort: 3000
  type: LoadBalancer
EOF
```

### 8. Configure Ingress

```bash
kubectl apply -f - <<EOF
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: cicd-ingress
  namespace: cicd
  annotations:
    cert-manager.io/cluster-issuer: letsencrypt-prod
spec:
  ingressClassName: nginx
  tls:
  - hosts:
    - yourdomain.com
    secretName: cicd-tls
  rules:
  - host: yourdomain.com
    http:
      paths:
      - path: /api
        pathType: Prefix
        backend:
          service:
            name: api
            port:
              number: 8080
      - path: /
        pathType: Prefix
        backend:
          service:
            name: frontend
            port:
              number: 3000
EOF
```

## Monitoring

### Prometheus Metrics

The API exposes Prometheus metrics on `/metrics`. Configure your Prometheus scrape job:

```yaml
scrape_configs:
  - job_name: 'cicd-api'
    kubernetes_sd_configs:
      - role: pod
        namespaces:
          names:
          - cicd
    relabel_configs:
      - source_labels: [__meta_kubernetes_pod_label_app]
        action: keep
        regex: api
```

### Logging with ELK Stack

Use Fluentd to ship logs to Elasticsearch:

```bash
kubectl apply -f - <<EOF
apiVersion: v1
kind: ConfigMap
metadata:
  name: fluentd-config
  namespace: cicd
data:
  fluent.conf: |
    <source>
      @type tail
      path /var/log/containers/cicd_*.log
      pos_file /var/log/fluentd-containers.log.pos
      tag kubernetes.*
      read_from_head true
      <parse>
        @type json
        time_format %Y-%m-%dT%H:%M:%S.%NZ
      </parse>
    </source>
    <match kubernetes.**>
      @type elasticsearch
      @id output_elasticsearch
      @log_level info
      include_tag_key true
      host elasticsearch.default.svc
      port 9200
      path _bulk
      logstash_format true
    </match>
EOF
```

## Scaling Recommendations

1. **API Server**: 2-10 replicas depending on load
2. **Job Processors**: 3-20 replicas, use HPA based on queue depth
3. **Database**: Use managed PostgreSQL service for production (e.g., RDS, CloudSQL)
4. **Redis**: Use managed Redis service (e.g., ElastiCache, MemoryStore)
5. **Storage**: Use cloud object storage (S3, GCS) for artifacts

## Security Best Practices

1. Enable RBAC and restrict API access
2. Use Network Policies to restrict traffic
3. Enable Pod Security Policies
4. Regular security scanning with Aqua/Snyk
5. Secrets management with Vault or external secret management
6. Enable audit logging
7. Use resource quotas and limits

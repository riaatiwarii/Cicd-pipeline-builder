-- PostgreSQL Schema for CI/CD Pipeline Builder

-- Extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "citext";

-- Enum Types
CREATE TYPE pipeline_status AS ENUM ('active', 'inactive', 'archived');
CREATE TYPE build_status AS ENUM ('pending', 'running', 'success', 'failed', 'cancelled');
CREATE TYPE trigger_type AS ENUM ('webhook', 'schedule', 'manual');
CREATE TYPE provider_type AS ENUM ('github', 'gitlab', 'gitea', 'custom');

-- Users Table
CREATE TABLE users (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  username CITEXT UNIQUE NOT NULL,
  email CITEXT UNIQUE NOT NULL,
  password_hash VARCHAR(255) NOT NULL,
  full_name VARCHAR(255),
  role VARCHAR(50) DEFAULT 'user',
  is_active BOOLEAN DEFAULT true,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  last_login TIMESTAMP
);

-- Pipelines Table
CREATE TABLE pipelines (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  name VARCHAR(255) NOT NULL,
  description TEXT,
  repo_url VARCHAR(500) NOT NULL,
  repo_branch VARCHAR(255) DEFAULT 'main',
  config_path VARCHAR(500) DEFAULT '.cicd.yml',
  status pipeline_status DEFAULT 'active',
  is_public BOOLEAN DEFAULT false,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  UNIQUE(user_id, name)
);

-- Stages Table
CREATE TABLE stages (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  pipeline_id UUID NOT NULL REFERENCES pipelines(id) ON DELETE CASCADE,
  name VARCHAR(255) NOT NULL,
  stage_order INT NOT NULL,
  description TEXT,
  docker_image VARCHAR(500) NOT NULL DEFAULT 'ubuntu:22.04',
  script TEXT NOT NULL,
  timeout_seconds INT DEFAULT 3600,
  allow_failure BOOLEAN DEFAULT false,
  needs_docker BOOLEAN DEFAULT false,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  UNIQUE(pipeline_id, name)
);

-- Builds Table (represents a pipeline run)
CREATE TABLE builds (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  pipeline_id UUID NOT NULL REFERENCES pipelines(id) ON DELETE CASCADE,
  status build_status DEFAULT 'pending',
  branch VARCHAR(255),
  commit_hash VARCHAR(40),
  commit_message TEXT,
  triggered_by UUID REFERENCES users(id),
  trigger_type trigger_type NOT NULL,
  started_at TIMESTAMP,
  completed_at TIMESTAMP,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Jobs Table (individual stage execution within a build)
CREATE TABLE jobs (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  build_id UUID NOT NULL REFERENCES builds(id) ON DELETE CASCADE,
  stage_id UUID NOT NULL REFERENCES stages(id) ON DELETE CASCADE,
  status build_status DEFAULT 'pending',
  exit_code INT,
  started_at TIMESTAMP,
  completed_at TIMESTAMP,
  container_id VARCHAR(255),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Logs Table (streaming logs from jobs)
CREATE TABLE logs (
  id BIGSERIAL PRIMARY KEY,
  job_id UUID NOT NULL REFERENCES jobs(id) ON DELETE CASCADE,
  line_number INT,
  log_level VARCHAR(50) DEFAULT 'info',
  message TEXT NOT NULL,
  timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create index on logs for faster queries
CREATE INDEX idx_logs_job_id ON logs(job_id);
CREATE INDEX idx_logs_timestamp ON logs(timestamp DESC);

-- Webhooks Table
CREATE TABLE webhooks (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  pipeline_id UUID NOT NULL REFERENCES pipelines(id) ON DELETE CASCADE,
  provider provider_type NOT NULL,
  events TEXT[] DEFAULT ARRAY['push', 'pull_request'],
  secret_key VARCHAR(500),
  webhook_url VARCHAR(500),
  is_active BOOLEAN DEFAULT true,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Triggers Table (schedules and manual triggers)
CREATE TABLE triggers (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  pipeline_id UUID NOT NULL REFERENCES pipelines(id) ON DELETE CASCADE,
  name VARCHAR(255) NOT NULL,
  trigger_type trigger_type NOT NULL,
  cron_expression VARCHAR(255), -- For scheduled triggers
  is_active BOOLEAN DEFAULT true,
  last_triggered_at TIMESTAMP,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Artifacts Table
CREATE TABLE artifacts (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  build_id UUID NOT NULL REFERENCES builds(id) ON DELETE CASCADE,
  job_id UUID REFERENCES jobs(id) ON DELETE SET NULL,
  name VARCHAR(255) NOT NULL,
  file_path VARCHAR(500) NOT NULL,
  file_size BIGINT,
  mime_type VARCHAR(100),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  expires_at TIMESTAMP
);

-- Credentials Table (encrypted credentials)
CREATE TABLE credentials (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  credential_type VARCHAR(50) NOT NULL, -- 'github_token', 'docker_registry', 'ssh_key', etc
  name VARCHAR(255) NOT NULL,
  encrypted_value TEXT NOT NULL,
  is_active BOOLEAN DEFAULT true,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  UNIQUE(user_id, name)
);

-- Environment Variables Table
CREATE TABLE env_variables (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  pipeline_id UUID NOT NULL REFERENCES pipelines(id) ON DELETE CASCADE,
  name VARCHAR(255) NOT NULL,
  value TEXT,
  is_secret BOOLEAN DEFAULT false,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  UNIQUE(pipeline_id, name)
);

-- Webhook Events Log (for audit trail)
CREATE TABLE webhook_events (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  webhook_id UUID NOT NULL REFERENCES webhooks(id) ON DELETE CASCADE,
  build_id UUID REFERENCES builds(id),
  event_type VARCHAR(100),
  payload JSONB,
  status VARCHAR(50),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for performance
CREATE INDEX idx_pipelines_user_id ON pipelines(user_id);
CREATE INDEX idx_pipelines_status ON pipelines(status);
CREATE INDEX idx_builds_pipeline_id ON builds(pipeline_id);
CREATE INDEX idx_builds_status ON builds(status);
CREATE INDEX idx_builds_created_at ON builds(created_at DESC);
CREATE INDEX idx_jobs_build_id ON jobs(build_id);
CREATE INDEX idx_jobs_stage_id ON jobs(stage_id);
CREATE INDEX idx_jobs_status ON jobs(status);
CREATE INDEX idx_stages_pipeline_id ON stages(pipeline_id);
CREATE INDEX idx_webhooks_pipeline_id ON webhooks(pipeline_id);
CREATE INDEX idx_triggers_pipeline_id ON triggers(pipeline_id);
CREATE INDEX idx_artifacts_build_id ON artifacts(build_id);
CREATE INDEX idx_credentials_user_id ON credentials(user_id);
CREATE INDEX idx_env_variables_pipeline_id ON env_variables(pipeline_id);
CREATE INDEX idx_webhook_events_webhook_id ON webhook_events(webhook_id);

-- Create default admin user (password: admin123, use a hash in production)
INSERT INTO users (username, email, password_hash, full_name, role, is_active)
VALUES ('admin', 'admin@cicd.local', '$2a$10$Lx5Z5r5Z5Z5Z5Z5Z5Z5Z5', 'Administrator', 'admin', true)
ON CONFLICT DO NOTHING;

-- Create updated_at trigger function
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create triggers for updated_at
CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_pipelines_updated_at BEFORE UPDATE ON pipelines
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_stages_updated_at BEFORE UPDATE ON stages
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_builds_updated_at BEFORE UPDATE ON builds
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_jobs_updated_at BEFORE UPDATE ON jobs
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_webhooks_updated_at BEFORE UPDATE ON webhooks
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_triggers_updated_at BEFORE UPDATE ON triggers
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_credentials_updated_at BEFORE UPDATE ON credentials
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_env_variables_updated_at BEFORE UPDATE ON env_variables
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

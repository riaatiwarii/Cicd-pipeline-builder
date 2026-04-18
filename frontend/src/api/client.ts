import axios from 'axios';

const API_URL = import.meta.env.REACT_APP_API_URL || 'http://localhost:8080';

const api = axios.create({
  baseURL: `${API_URL}/api/v1`,
});

// Add token to requests
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// Auth API
export const authAPI = {
  login: (username: string, password: string) =>
    api.post('/auth/login', { username, password }),
  register: (username: string, email: string, password: string, fullName: string) =>
    api.post('/auth/register', { username, email, password, full_name: fullName }),
  refreshToken: () => api.post('/auth/refresh'),
};

// Pipeline API
export const pipelineAPI = {
  list: () => api.get('/pipelines'),
  get: (id: string) => api.get(`/pipelines/${id}`),
  create: (pipeline: any) => api.post('/pipelines', pipeline),
  update: (id: string, pipeline: any) => api.put(`/pipelines/${id}`, pipeline),
  delete: (id: string) => api.delete(`/pipelines/${id}`),
  trigger: (id: string) => api.post(`/pipelines/${id}/trigger`),
};

// Build API
export const buildAPI = {
  list: (pipelineId: string, limit = 10, offset = 0) =>
    api.get(`/pipelines/${pipelineId}/builds`, { params: { limit, offset } }),
  get: (id: string) => api.get(`/builds/${id}`),
  getLogs: (id: string) => api.get(`/builds/${id}/logs`),
  getArtifacts: (id: string) => api.get(`/builds/${id}/artifacts`),
  cancel: (id: string) => api.post(`/builds/${id}/cancel`),
};

// Webhook API
export const webhookAPI = {
  list: () => api.get('/webhooks'),
  get: (id: string) => api.get(`/webhooks/${id}`),
  create: (webhook: any) => api.post('/webhooks', webhook),
  update: (id: string, webhook: any) => api.put(`/webhooks/${id}`, webhook),
  delete: (id: string) => api.delete(`/webhooks/${id}`),
};

// Trigger API
export const triggerAPI = {
  list: () => api.get('/triggers'),
  get: (id: string) => api.get(`/triggers/${id}`),
  create: (trigger: any) => api.post('/triggers', trigger),
  update: (id: string, trigger: any) => api.put(`/triggers/${id}`, trigger),
  delete: (id: string) => api.delete(`/triggers/${id}`),
};

export default api;

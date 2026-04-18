import { create } from 'zustand';

interface AuthState {
  token: string | null;
  user: any | null;
  isAuthenticated: boolean;
  login: (token: string, user: any) => void;
  logout: () => void;
  setUser: (user: any) => void;
}

export const useAuthStore = create<AuthState>((set) => ({
  token: localStorage.getItem('token'),
  user: JSON.parse(localStorage.getItem('user') || 'null'),
  isAuthenticated: !!localStorage.getItem('token'),
  login: (token, user) => {
    localStorage.setItem('token', token);
    localStorage.setItem('user', JSON.stringify(user));
    set({ token, user, isAuthenticated: true });
  },
  logout: () => {
    localStorage.removeItem('token');
    localStorage.removeItem('user');
    set({ token: null, user: null, isAuthenticated: false });
  },
  setUser: (user) => {
    localStorage.setItem('user', JSON.stringify(user));
    set({ user });
  },
}));

interface PipelineState {
  pipelines: any[];
  currentPipeline: any | null;
  loading: boolean;
  error: string | null;
  setPipelines: (pipelines: any[]) => void;
  setCurrentPipeline: (pipeline: any) => void;
  setLoading: (loading: boolean) => void;
  setError: (error: string | null) => void;
}

export const usePipelineStore = create<PipelineState>((set) => ({
  pipelines: [],
  currentPipeline: null,
  loading: false,
  error: null,
  setPipelines: (pipelines) => set({ pipelines }),
  setCurrentPipeline: (pipeline) => set({ currentPipeline: pipeline }),
  setLoading: (loading) => set({ loading }),
  setError: (error) => set({ error }),
}));

interface BuildState {
  builds: any[];
  currentBuild: any | null;
  logs: string[];
  loading: boolean;
  error: string | null;
  setBuilds: (builds: any[]) => void;
  setCurrentBuild: (build: any) => void;
  setLogs: (logs: string[]) => void;
  addLog: (log: string) => void;
  setLoading: (loading: boolean) => void;
  setError: (error: string | null) => void;
}

export const useBuildStore = create<BuildState>((set) => ({
  builds: [],
  currentBuild: null,
  logs: [],
  loading: false,
  error: null,
  setBuilds: (builds) => set({ builds }),
  setCurrentBuild: (build) => set({ currentBuild: build }),
  setLogs: (logs) => set({ logs }),
  addLog: (log) => set((state) => ({ logs: [...state.logs, log] })),
  setLoading: (loading) => set({ loading }),
  setError: (error) => set({ error }),
}));

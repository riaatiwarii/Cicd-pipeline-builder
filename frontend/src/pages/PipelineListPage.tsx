import React, { useEffect, useState } from 'react';
import {
  Box,
  Button,
  Container,
  Grid,
  Card,
  CardContent,
  CardActions,
  Typography,
  CircularProgress,
  Alert,
  Dialog,
  TextField,
} from '@mui/material';
import { useNavigate } from 'react-router-dom';
import { pipelineAPI } from '../api/client';
import { usePipelineStore } from '../store/store';
import AddIcon from '@mui/icons-material/Add';
import PlayArrowIcon from '@mui/icons-material/PlayArrow';

export const PipelineListPage: React.FC = () => {
  const navigate = useNavigate();
  const { pipelines, loading, error, setPipelines, setLoading, setError } = usePipelineStore();
  const [openCreateDialog, setOpenCreateDialog] = useState(false);
  const [newPipeline, setNewPipeline] = useState({ name: '', repo_url: '', config_path: '.cicd.yml' });

  useEffect(() => {
    fetchPipelines();
  }, []);

  const fetchPipelines = async () => {
    setLoading(true);
    try {
      const { data } = await pipelineAPI.list();
      setPipelines(data || []);
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to fetch pipelines');
    } finally {
      setLoading(false);
    }
  };

  const handleCreatePipeline = async () => {
    try {
      await pipelineAPI.create(newPipeline);
      setNewPipeline({ name: '', repo_url: '', config_path: '.cicd.yml' });
      setOpenCreateDialog(false);
      fetchPipelines();
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to create pipeline');
    }
  };

  const handleTriggerPipeline = async (id: string) => {
    try {
      await pipelineAPI.trigger(id);
      fetchPipelines();
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to trigger pipeline');
    }
  };

  return (
    <Container maxWidth="lg" sx={{ py: 4 }}>
      <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 4 }}>
        <Typography variant="h4">Pipelines</Typography>
        <Button
          variant="contained"
          startIcon={<AddIcon />}
          onClick={() => setOpenCreateDialog(true)}
        >
          New Pipeline
        </Button>
      </Box>

      {error && <Alert severity="error" sx={{ mb: 2 }}>{error}</Alert>}

      {loading ? (
        <Box sx={{ display: 'flex', justifyContent: 'center' }}>
          <CircularProgress />
        </Box>
      ) : (
        <Grid container spacing={2}>
          {pipelines.map((pipeline: any) => (
            <Grid item xs={12} sm={6} md={4} key={pipeline.id}>
              <Card>
                <CardContent>
                  <Typography variant="h6" gutterBottom>
                    {pipeline.name}
                  </Typography>
                  <Typography color="textSecondary" variant="body2">
                    {pipeline.repo_url}
                  </Typography>
                  <Typography variant="caption" color="textSecondary">
                    Branch: {pipeline.repo_branch}
                  </Typography>
                </CardContent>
                <CardActions>
                  <Button size="small" onClick={() => navigate(`/pipelines/${pipeline.id}`)}>
                    View
                  </Button>
                  <Button
                    size="small"
                    startIcon={<PlayArrowIcon />}
                    onClick={() => handleTriggerPipeline(pipeline.id)}
                  >
                    Run
                  </Button>
                </CardActions>
              </Card>
            </Grid>
          ))}
        </Grid>
      )}

      <Dialog open={openCreateDialog} onClose={() => setOpenCreateDialog(false)}>
        <Box sx={{ p: 3, minWidth: 400 }}>
          <Typography variant="h6" gutterBottom>
            Create New Pipeline
          </Typography>
          <TextField
            fullWidth
            label="Pipeline Name"
            value={newPipeline.name}
            onChange={(e) => setNewPipeline({ ...newPipeline, name: e.target.value })}
            margin="normal"
          />
          <TextField
            fullWidth
            label="Repository URL"
            value={newPipeline.repo_url}
            onChange={(e) => setNewPipeline({ ...newPipeline, repo_url: e.target.value })}
            margin="normal"
          />
          <TextField
            fullWidth
            label="Config Path"
            value={newPipeline.config_path}
            onChange={(e) => setNewPipeline({ ...newPipeline, config_path: e.target.value })}
            margin="normal"
          />
          <Box sx={{ mt: 3, display: 'flex', gap: 1 }}>
            <Button variant="contained" onClick={handleCreatePipeline}>
              Create
            </Button>
            <Button variant="outlined" onClick={() => setOpenCreateDialog(false)}>
              Cancel
            </Button>
          </Box>
        </Box>
      </Dialog>
    </Container>
  );
};

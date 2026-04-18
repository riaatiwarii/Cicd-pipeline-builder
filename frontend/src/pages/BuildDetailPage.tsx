import React, { useEffect, useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import {
  Box,
  Button,
  Container,
  Typography,
  CircularProgress,
  Paper,
  Stepper,
  Step,
  StepLabel,
  Chip,
  Card,
  CardContent,
} from '@mui/material';
import { buildAPI } from '../api/client';
import BackIcon from '@mui/icons-material/ArrowBack';

export const BuildDetailPage: React.FC = () => {
  const { buildId } = useParams<{ buildId: string }>();
  const navigate = useNavigate();
  const [build, setBuild] = useState<any>(null);
  const [logs, setLogs] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetchData();
    // Poll for updates every 2 seconds
    const interval = setInterval(fetchData, 2000);
    return () => clearInterval(interval);
  }, [buildId]);

  const fetchData = async () => {
    if (!buildId) return;
    try {
      const [buildRes, logsRes] = await Promise.all([
        buildAPI.get(buildId),
        buildAPI.getLogs(buildId),
      ]);
      setBuild(buildRes.data);
      setLogs(logsRes.data || []);
    } catch (err) {
      console.error('Failed to fetch build details', err);
    } finally {
      setLoading(false);
    }
  };

  if (loading) {
    return (
      <Container sx={{ display: 'flex', justifyContent: 'center', py: 10 }}>
        <CircularProgress />
      </Container>
    );
  }

  const getStepStatus = (jobStatus: string) => {
    if (jobStatus === 'success') return 'completed';
    if (jobStatus === 'failed') return 'error';
    if (jobStatus === 'running') return 'in_progress';
    return 'not_started';
  };

  return (
    <Container maxWidth="lg" sx={{ py: 4 }}>
      <Button
        startIcon={<BackIcon />}
        onClick={() => navigate(-1)}
        sx={{ mb: 2 }}
      >
        Back
      </Button>

      <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 4 }}>
        <Box>
          <Typography variant="h4" gutterBottom>
            Build Details
          </Typography>
          <Typography color="textSecondary">
            {build?.commit_hash?.substring(0, 7)} on {build?.branch}
          </Typography>
        </Box>
        <Chip
          label={build?.status}
          color={build?.status === 'success' ? 'success' : build?.status === 'running' ? 'default' : 'error'}
          variant="outlined"
        />
      </Box>

      <Paper sx={{ p: 3, mb: 4 }}>
        <Typography variant="h6" gutterBottom>
          Pipeline Progress
        </Typography>
        {build?.jobs && build.jobs.length > 0 && (
          <Stepper activeStep={build.jobs.length - 1} alternativeLabel>
            {build.jobs.map((job: any, index: number) => (
              <Step key={index}>
                <StepLabel>{job.stage_id}</StepLabel>
              </Step>
            ))}
          </Stepper>
        )}
      </Paper>

      <Paper sx={{ p: 3 }}>
        <Typography variant="h6" gutterBottom>
          Logs
        </Typography>
        <Box
          component="pre"
          sx={{
            bg: '#1e1e1e',
            color: '#d4d4d4',
            p: 2,
            borderRadius: 1,
            overflow: 'auto',
            maxHeight: '500px',
            fontSize: '0.875rem',
            fontFamily: 'monospace',
          }}
        >
          {logs.length === 0 ? (
            'No logs available'
          ) : (
            logs.map((log: any, index: number) => (
              <div key={index}>
                {log.message}
              </div>
            ))
          )}
        </Box>
      </Paper>
    </Container>
  );
};

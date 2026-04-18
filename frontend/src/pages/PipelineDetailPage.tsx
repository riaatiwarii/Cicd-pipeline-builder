import React, { useEffect, useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import {
  Box,
  Button,
  Container,
  Typography,
  Card,
  CardContent,
  Chip,
  CircularProgress,
  Paper,
  Tab,
  Tabs,
} from '@mui/material';
import { pipelineAPI, buildAPI } from '../api/client';
import PlayArrowIcon from '@mui/icons-material/PlayArrow';
import BackIcon from '@mui/icons-material/ArrowBack';

interface TabPanelProps {
  children?: React.ReactNode;
  index: number;
  value: number;
}

function TabPanel(props: TabPanelProps) {
  const { children, value, index, ...other } = props;

  return (
    <div
      role="tabpanel"
      hidden={value !== index}
      id={`tabpanel-${index}`}
      aria-labelledby={`tab-${index}`}
      {...other}
    >
      {value === index && <Box sx={{ p: 3 }}>{children}</Box>}
    </div>
  );
}

export const PipelineDetailPage: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const [pipeline, setPipeline] = useState<any>(null);
  const [builds, setBuilds] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);
  const [tabValue, setTabValue] = useState(0);

  useEffect(() => {
    fetchData();
  }, [id]);

  const fetchData = async () => {
    if (!id) return;
    setLoading(true);
    try {
      const [pipelineRes, buildsRes] = await Promise.all([
        pipelineAPI.get(id),
        buildAPI.list(id),
      ]);
      setPipeline(pipelineRes.data);
      setBuilds(buildsRes.data || []);
    } catch (err) {
      console.error('Failed to fetch pipeline details', err);
    } finally {
      setLoading(false);
    }
  };

  const handleTrigger = async () => {
    if (!id) return;
    try {
      await pipelineAPI.trigger(id);
      fetchData();
    } catch (err) {
      console.error('Failed to trigger pipeline', err);
    }
  };

  if (loading) {
    return (
      <Container sx={{ display: 'flex', justifyContent: 'center', py: 10 }}>
        <CircularProgress />
      </Container>
    );
  }

  return (
    <Container maxWidth="lg" sx={{ py: 4 }}>
      <Button
        startIcon={<BackIcon />}
        onClick={() => navigate('/pipelines')}
        sx={{ mb: 2 }}
      >
        Back
      </Button>

      <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 4 }}>
        <Box>
          <Typography variant="h4" gutterBottom>
            {pipeline?.name}
          </Typography>
          <Typography color="textSecondary">
            Repository: {pipeline?.repo_url}
          </Typography>
          <Typography variant="caption" color="textSecondary">
            Branch: {pipeline?.repo_branch}
          </Typography>
        </Box>
        <Button
          variant="contained"
          startIcon={<PlayArrowIcon />}
          size="large"
          onClick={handleTrigger}
        >
          Trigger Build
        </Button>
      </Box>

      <Paper sx={{ width: '100%' }}>
        <Tabs value={tabValue} onChange={(e, newValue) => setTabValue(newValue)}>
          <Tab label="Builds" id="tab-0" aria-controls="tabpanel-0" />
          <Tab label="Configuration" id="tab-1" aria-controls="tabpanel-1" />
        </Tabs>

        <TabPanel value={tabValue} index={0}>
          <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2 }}>
            {builds.length === 0 ? (
              <Typography color="textSecondary">No builds yet</Typography>
            ) : (
              builds.map((build: any) => (
                <Card key={build.id}>
                  <CardContent>
                    <Box sx={{ display: 'flex', justifyContent: 'space-between' }}>
                      <Box>
                        <Typography variant="subtitle1">Build #{build.id.substring(0, 8)}</Typography>
                        <Typography color="textSecondary" variant="body2">
                          {build.branch} - {build.commit_hash?.substring(0, 7)}
                        </Typography>
                      </Box>
                      <Chip
                        label={build.status}
                        color={build.status === 'success' ? 'success' : 'error'}
                        variant="outlined"
                      />
                    </Box>
                  </CardContent>
                </Card>
              ))
            )}
          </Box>
        </TabPanel>

        <TabPanel value={tabValue} index={1}>
          <Box
            component="pre"
            sx={{
              bg: '#f5f5f5',
              p: 2,
              borderRadius: 1,
              overflow: 'auto',
              fontSize: '0.875rem',
            }}
          >
            Configuration will be displayed here
          </Box>
        </TabPanel>
      </Paper>
    </Container>
  );
};

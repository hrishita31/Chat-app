import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import Card from '@mui/material/Card';
import CardContent from '@mui/material/CardContent';
import CardActions from '@mui/material/CardActions';
import Typography from '@mui/material/Typography';
import TextField from '@mui/material/TextField';
import Button from '@mui/material/Button';
import Box from '@mui/material/Box';
import Link from '@mui/material/Link';
import PersonIcon from '@mui/icons-material/Person';
import HttpsIcon from '@mui/icons-material/Https';

const SignUp = () => {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [name, setName] = useState('');
  const [message, setMessage] = useState('');
  const navigate = useNavigate();


  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      const response = await fetch(`${import.meta.env.VITE_HOST}/user/?name=${name}&username=${username}&password=${password}`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
      });
      const data = await response.json();
      if (response.ok) {
        setMessage(data.message || 'Sign up successful');
        navigate('/login'); // Redirect to Login page
      } else {
        setMessage(data.message || 'Sign up failed');
      }
    } catch (error) {
      setMessage('An error occurred. Please try again later.');
    }
  };

  return (
    <Box className="flex justify-center items-center h-screen">
      <Card className="w-full max-w-md shadow-md rounded-lg p-8">
        <CardContent className="flex flex-col items-center">
        <img
            src="src/assets/chat.jpg"
            alt="Chat Icon"
            className="w-24 h-24 mb-4"
          />
          <Typography variant="h5" component="div" className="text-center mb-4">
            Sign Up for Chat App
          </Typography>
          {message && (
            <Typography variant="body2" color="error" className="mb-4">
              {message}
            </Typography>
          )}
          <form onSubmit={handleSubmit} className="w-full">
            <TextField
              label="Name"
              variant="outlined"
              fullWidth
              margin="normal"
              value={name}
              onChange={(e) => setName(e.target.value)}
              
              sx={{
                '& .MuiOutlinedInput-root': {
                  borderRadius: '50px',
                },
              }}
            />
            <TextField
              label="Username"
              variant="outlined"
              fullWidth
              margin="normal"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              InputProps={{
                startAdornment: (
                  <Box mr={1}>
                    <PersonIcon color="action" />
                  </Box>
                ),
               
              }}
              sx={{
                '& .MuiOutlinedInput-root': {
                  borderRadius: '50px',
                },
              }}
            />
            <TextField
              label="Password"
              variant="outlined"
              type="password"
              fullWidth
              margin="normal"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              InputProps={{
                startAdornment: (
                  <Box mr={1}>
                    <HttpsIcon color="action" />
                  </Box>
                ),
               
              }}
              sx={{
                '& .MuiOutlinedInput-root': {
                  borderRadius: '50px',
                },
              }}
            />
            <CardActions className="flex flex-col w-full gap-4">
              <Button type="submit" variant="contained" sx={{ width: '100%', borderRadius: '50px', paddingY: 2, marginTop: 2, backgroundColor: '#9E97F4', ':hover': { backgroundColor: '#7a70e3' } }}>
                Sign Up
              </Button>
              <Typography variant="body2" className="text-center mt-4">
                Already have an account?{' '}
                <Link href="/login" sx={{ color: '#E91A78' }}>
                  Login
                </Link>
              </Typography>
            </CardActions>
          </form>
        </CardContent>
      </Card>
    </Box>
  );
};

export default SignUp;

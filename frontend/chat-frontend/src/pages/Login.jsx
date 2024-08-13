import { useState } from "react";
import { useNavigate } from "react-router-dom";
import Card from "@mui/material/Card";
import CardContent from "@mui/material/CardContent";
import CardActions from "@mui/material/CardActions";
import Typography from "@mui/material/Typography";
import TextField from "@mui/material/TextField";
import Button from "@mui/material/Button";
import Box from "@mui/material/Box";
import Link from "@mui/material/Link";
import PersonIcon from '@mui/icons-material/Person';
import HttpsIcon from '@mui/icons-material/Https';

const LoginPage = () => {
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [message, setMessage] = useState("");
  const navigate = useNavigate();

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      const response = await fetch(
        `${import.meta.env.VITE_HOST}/user/loginUser?username=${username}&password=${password}`,
        {
          method: "GET",
        }
      );

      const data = await response.json();
      if (response.ok) {
        setMessage(data.message || "Login successful");
        localStorage.setItem("token", data.data);
        localStorage.setItem("username", username);
        navigate("/home"); // Redirect to Home page
      } else {
        setMessage(data.message || "Login failed");
      }
    } catch (error) {
      setMessage("An error occurred. Please try again later.");
    }
  };

  return (
    <Box className="flex justify-center items-center h-screen">
      <Card className="w-full max-w-md shadow-md rounded-lg p-8">
        <CardContent className="flex flex-col items-center">
          <img
            src="src/assets/register.jpg"
            
            alt="Chat Icon"
            className="w-30 h-30 mb-4"
          />
          <Typography variant="h5" component="div" className="text-center mb-4">
            Login to Chat App
          </Typography>
          <form onSubmit={handleSubmit} className="w-full">
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
                "& .MuiOutlinedInput-root": {
                  borderRadius: "50px",
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
                "& .MuiOutlinedInput-root": {
                  borderRadius: "50px",
                },
              }}
            />
            <CardActions className="flex flex-col w-full gap-4">
              <Button
                type="submit"
                variant="contained"
                sx={{
                  width: "100%",
                  borderRadius: "50px",
                  paddingY: 2,
                  marginTop: 2,
                  backgroundColor: "#9E97F4",
                  ":hover": { backgroundColor: "#7a70e3" },
                }}
              >
                Login
              </Button>
              {message && (
                <Typography variant="body2" color="error" className="mb-4">
                  {message}
                </Typography>
              )}
              <Typography variant="body2" className="text-center mt-4">
                Don&apos;t have an account?{" "}
                <Link href="/sign-up" sx={{ color: "#E91A78" }}>
                  Sign Up
                </Link>
              </Typography>
            </CardActions>
          </form>
        </CardContent>
      </Card>
    </Box>
  );
};

export default LoginPage;

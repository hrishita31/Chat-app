import {  useNavigate } from "react-router-dom";
import CardActions from "@mui/material/CardActions";
import Typography from "@mui/material/Typography";
import Button from "@mui/material/Button";
// import "../layout/Home.css";



const OnBoarding = () => {
  const navigate = useNavigate();
  return (
    <div className="h-full mt-32">
      <img
        src="src/assets/welcomePage.jpg"
        alt="Chat Icon"
        className="w-50 h-50 mb-4 mt-4 rounded-md mx-auto"
      />
      <Typography variant="h5" component="div" className="text-center mb-4">
        Welcome to Chat App
      </Typography>

      <CardActions className="flex flex-col w-full gap-4">
        <Typography variant="body2" className="text-center mt-4">
          Chat with your friends and family Have fun!
        </Typography>
        <Button
          type="submit"
          variant="contained"
          onClick={()=>navigate("/login")}
          sx={{
            width: "30%",
            borderRadius: "50px",
            paddingY: 2,
            marginTop: 2,
            backgroundColor: "#9E97F4",
            ":hover": { backgroundColor: "#7a70e3" },
          }}
        >
          Get Started!
        </Button>
      </CardActions>
    </div>
  );
};

export default OnBoarding;

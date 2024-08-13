import { useState, useEffect } from "react";
import TextField from "@mui/material/TextField";
import PersonSearchIcon from "@mui/icons-material/PersonSearch";
import Box from "@mui/material/Box";
import Button from "@mui/material/Button";
import LogoutIcon from "@mui/icons-material/Logout";
import { useNavigate } from "react-router-dom";

const Sidebar = ({ userProfile, setSelectedChat }) => {
  const [friends, setFriends] = useState([]);
  const [searchResults, setSearchResults] = useState([]);
  const [searchTerm, setSearchTerm] = useState("");
  const navigate = useNavigate();

  useEffect(() => {
    

    fetchFriends();
  }, []);
  const fetchFriends = async () => {
    try {
      const response = await fetch(
        `${import.meta.env.VITE_HOST}/friend/findAllFriends/${localStorage.getItem(
          "username"
        )}`,
        {
          method: "GET",
          headers: {
            Authorization: `Bearer ${localStorage.getItem("token")}`,
          },
        }
      );

      const data = await response.json();
      if (response.ok) {
        setFriends(data.data || []);
      } else {
        console.error("Failed to fetch friends");
      }
    } catch (error) {
      console.error("An error occurred:", error);
    }
  };

  const handleSearchChange = (e) => {
    const searchTerm = e.target.value;
    setSearchTerm(searchTerm);
    if (searchTerm.length > 0) {
      handleSearch(searchTerm);
    } else {
      setSearchResults([]);
    }
  };

  const handleSearch = async (searchTerm) => {
    try {
      const response = await fetch(
        `${import.meta.env.VITE_HOST}/user/search/${searchTerm}`,
        {
          method: "GET",
          headers: {
            Authorization: `Bearer ${localStorage.getItem("token")}`,
          },
        }
      );

      const data = await response.json();
      if (response.ok) {
        setSearchResults(data.data || ["No results found"]);
      } else {
        console.error("Failed to fetch search results");
      }
    } catch (error) {
      console.error("An error occurred:", error);
    }
  };

  const handleAddFriend = async (friendUsername) => {
    try {
      const response = await fetch(
        `${import.meta.env.VITE_HOST}/friend/addFriend?username1=${localStorage.getItem(
          "username"
        )}&username2=${friendUsername}`,
        {
          method: "GET",
          headers: {
            Authorization: `Bearer ${localStorage.getItem("token")}`,
          },
        }
      );

      const data = await response.json();
      if (response.ok) {
        // Refetch friends after adding a new friend
        fetchFriends();
        // Clear search results and selected chat
        setSearchResults([]);
        setSelectedChat(null);
        setSearchTerm("");
      } else {
        console.error("Failed to add friend");
      }
    } catch (error) {
      console.error("An error occurred:", error);
    }
    fetchFriends();
  };

  const handleLogout = () => {
    localStorage.removeItem("token");
    localStorage.removeItem("username");
    navigate("/login");
  };

  return (
    <div className="w-1/4 bg-gray-300 p-4 shadow-lg flex flex-col justify-between h-auto">
      <div>
        <div className="flex items-center mb-4">
          <img
            src="https://images.rawpixel.com/image_png_800/cHJpdmF0ZS9sci9pbWFnZXMvd2Vic2l0ZS8yMDIzLTAxL3JtNjA5LXNvbGlkaWNvbi13LTAwMi1wLnBuZw.png"
            alt="Profile"
            className="w-12 h-12 rounded-full mr-4"
          />
          <div className="grid grid-rows-2">
            <span className="font-bold text-lg">{userProfile.Name}</span>
            <span className="text-gray-600">{userProfile.Username}</span>
          </div>
        </div>
        <div className="w-full mb-4 relative">
          <TextField
            label="Search"
            variant="outlined"
            fullWidth
            margin="normal"
            value={searchTerm}
            onChange={handleSearchChange}
            InputProps={{
              startAdornment: (
                <Box mr={1}>
                  <PersonSearchIcon color="action" />
                </Box>
              ),
            }}
            sx={{
              "& .MuiOutlinedInput-root": {
                borderRadius: "20px",
              },
            }}
          />
          {searchResults.length > 0 && (
            <div className="absolute top-16 w-[100%] p-3 left-[0.25rem] right-5 bg-white z-10 rounded-b-md">
              {searchResults.map((result) => (
                <div key={result} className="flex flex-row gap-1">
                  <p
                    className="flex items-center p-2 rounded-sm m-1 h-9 cursor-pointer"
                    onClick={() => {
                      setSelectedChat(result); // Assuming result is the username
                      setSearchResults([]);
                      setSearchTerm("");
                    }}
                  >
                    {result}
                    <button className="bg-slate-500 rounded-md text-white px-2 py-1 ml-2"
                      onClick={() => {
                        handleAddFriend(result); // Pass searchTerm or selected result here
                      }}
                    >
                      Add Friend
                    </button>
                  </p>
                </div>
              ))}
            </div>
          )}
        </div>
        <div className="friends-list">
          <h3 className="text-gray-700 font-semibold mb-2">Friends</h3>
          {friends.length === 0 ? (
            <p className="text-gray-500">You have no friends.</p>
          ) : (
            friends.map((friend) => (
              <div
                key={friend}
                className="friend-item py-2 cursor-pointer hover:bg-gray-100 rounded-md mb-2 p-4"
                onClick={() => setSelectedChat(friend)}
              >
                <p className="text-gray-800">{friend}</p>
              </div>
            ))
          )}
        </div>
      </div>
      <div className="flex items-center">
        <Button
          variant="contained"
          sx={{
            backgroundColor: "#9E97F4",
            ":hover": { backgroundColor: "#7a70e3" },
            color: "black",
          }}
          endIcon={<LogoutIcon />}
          onClick={handleLogout}
        >
          LogOut
        </Button>
      </div>
    </div>
  );
};

export default Sidebar;

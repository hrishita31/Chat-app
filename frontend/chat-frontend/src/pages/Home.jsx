import { useState, useEffect } from "react";
import Sidebar from "../components/Sidebar";
import ChatWindow from "../components/ChatWindow";

const HomePage = () => {
  const [selectedChat, setSelectedChat] = useState(null);
  // const [searchResults, setSearchResults] = useState([]);
  const [userProfile, setUserProfile] = useState({});

  useEffect(() => {
    // Fetch user profile details on mount
    const fetchUserProfile = async () => {
      try {
        const response = await fetch(
          `${import.meta.env.VITE_HOST}/user/?username=${localStorage.getItem(
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
          setUserProfile(data.data);
        } else {
          console.error("Failed to fetch user profile");
        }
      } catch (error) {
        console.error("An error occurred:", error);
      }
    };

    fetchUserProfile();
  }, []);

  return (
    <div className=" flex min-h-screen max-h-screen bg-[#ebf4f8]">
      <Sidebar userProfile={userProfile} setSelectedChat={setSelectedChat} />
      {selectedChat ? (
        <ChatWindow selectedChat={selectedChat} />
      ) : (
        <div className="flex-grow flex items-center justify-center bg-white">
          <p className="text-white text-lg">Select a chat to start messaging</p>
        </div>
      )}
    </div>
  );
};

export default HomePage;

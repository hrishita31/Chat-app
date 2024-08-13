import { useEffect, useRef, useState } from "react";
import {
  Button,
  TextField,
  Typography,
  Box,
  Paper,
  List,
  ListItem,
  ListItemText,
} from "@mui/material";
import "../layout/ChatWindow.css";
import useWindowDimensions from './WindowDimensions';


// eslint-disable-next-line react/prop-types
const ChatWindow = ({ selectedChat }) => {
  const [messages, setMessages] = useState([]);
  const [inputMessage, setInputMessage] = useState("");
  const socketRef = useRef(null);
  const fetchChatHistory = async (selectedChat) => {
    try {
      const response = await fetch(
        `${import.meta.env.VITE_HOST}/ws/histories?username=${localStorage.getItem(
          "username"
        )}&receiver=${selectedChat}`,
        {
          method: "GET",
          headers: {
            Authorization: `Bearer ${localStorage.getItem("token")}`,
          },
        }
      );

      const data = await response.json();
      if (response.ok) {
        setMessages(data.data || []);
      } else {
        console.error("Failed to fetch chat history");
      }
    } catch (error) {
      console.error("An error occurred:", error);
    }
  };

  const messageEl = useRef(null);

  useEffect(() => {
    if (messageEl) {
      messageEl.current.addEventListener('DOMNodeInserted', event => {
        const { currentTarget: target } = event;
        target.scroll({ top: target.scrollHeight, behavior: 'smooth' });
      });
    }
  }, [])

  useEffect(() => {

    fetchChatHistory(selectedChat);

    const wsUrl = `${import.meta.env.VITE_WSHOST}/ws/?username=${localStorage.getItem(
      "username"
    )}`;
    if(socketRef.current) {
      socketRef.current.close()
    }
    socketRef.current = new WebSocket(wsUrl);

    socketRef.current.onopen = () => {
      console.log("WebSocket connection established");
    };

    socketRef.current.onmessage = (event) => {
      const message = JSON.parse(event.data);
      console.log("Message received:", message);
      setMessages((prevMsgs) => [message,...prevMsgs]);
    };

    socketRef.current.onerror = (error) => {
      console.error("WebSocket error:", error);
    };

    socketRef.current.onclose = () => {
      console.log('WebSocket connection closed');
    };

    // return () => {
    //   if (socketRef.current) {
    //     socketRef.current.close();
    //   }
    // };


  }, [selectedChat]);

  const handleSendMessage = () => {
    if (socketRef.current && inputMessage) {
      if (socketRef.current.readyState === WebSocket.OPEN) {
        const messageData = {
          sender: localStorage.getItem("username"),
          message: inputMessage,
          receiver: selectedChat,
          time: new Date().toISOString(),
        };
        socketRef.current.send(JSON.stringify(messageData));
        setMessages([...messages, messageData])
        setInputMessage("");
      } else {
        console.error(
          "WebSocket is not open. Current readyState:",
          socketRef.current.readyState
        );
      }
    }
  };

  
  const { height, width } = useWindowDimensions();

  return (
    <Box className="chat-window flex flex-col w-3/4 h-full-screen bg-white shadow-lg">
      <Paper
        sx={{
          padding: "12px",
          marginTop: "10px",
          marginLeft: "6px",
          marginRight: "6px",
          marginBottom: "4px",
          borderRadius: "6px",
          display: "flex",
          alignItems: "center",
          flexDirection: "row",
        }}
      >
        <img
          src="https://images.rawpixel.com/image_png_800/cHJpdmF0ZS9sci9pbWFnZXMvd2Vic2l0ZS8yMDIzLTAxL3JtNjA5LXNvbGlkaWNvbi13LTAwMi1wLnBuZw.png"
          alt="Profile"
          className="w-12 h-12 rounded-full mr-4"
        />
        <Typography variant="h6">{selectedChat}</Typography>
      </Paper>
      <Box className={"chat-messages flex-1 p-4 overflow-auto max-height="+height} ref={messageEl}>
        <List>
          {messages.sort((a,b) => a.time > b.time ? 1 : -1).map((message, index) => (
            <ListItem
              sx={{
                padding: "6px",
                marginLeft: "6px",
                marginRight: "6px",
                marginTop: "4px",
                marginBottom: "4px",
                borderRadius: "16px",
                borderBottomLeftRadius: "0px",
              }}
              key={index}
              className={
                message.receiver === localStorage.getItem("username")
                  ? "bg-blue-100"
                  : "bg-gray-100"
              }
            >
              <ListItemText
                primary={message.message}
                secondary={new Date(message.time).toLocaleTimeString()}
              />
            </ListItem>
          ))}
        </List>
      </Box>
      <Box className="chat-input p-4 bg-gray-200 flex">
        <TextField
          fullWidth
          variant="outlined"
          placeholder="Type a message"
          value={inputMessage}
          onKeyDown={(e) => { 
                        if (e.key === "Enter") { 
                            handleSendMessage()
                        } 
                    }
          }
          onChange={(e) => setInputMessage(e.target.value)}
        />
        <Button
          variant="contained"
          color="primary"
 
          onClick={handleSendMessage}
          className="ml-4"
        >
          Send
        </Button>
      </Box>
    </Box>
  );
};

export default ChatWindow;

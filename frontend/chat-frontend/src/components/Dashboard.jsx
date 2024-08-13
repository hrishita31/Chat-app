import { useState, useEffect } from 'react';
import '../layout/Dashboard.css';

const Dashboard = () => {
  const [userDetails, setUserDetails] = useState({});
  const [message, setMessage] = useState('');

  useEffect(() => {
    const fetchUserDetails = async () => {
      try {
        const response = await fetch(`${import.meta.env.VITE_HOST}/getUser?username=${encodeURIComponent(localStorage.getItem('username'))}`, {
          method: 'GET',
          headers: {
            'Authorization': `Bearer ${localStorage.getItem('token')}`,
          },
        });

        const data = await response.json();
        if (response.ok) {
          setUserDetails(data);
        } else {
          setMessage(data.message || 'Failed to fetch user details');
        }
      } catch (error) {
        setMessage('An error occurred. Please try again later.');
      }
    };

    fetchUserDetails();
  }, []);

  return (
    <div className="dashboard-page">
      <h2>User Dashboard</h2>
      {message && <p className="error-message">{message}</p>}
      {userDetails ? (
        <div className="user-details">
          <p><strong>Name:</strong> {userDetails.name}</p>
          <p><strong>Username:</strong> {userDetails.username}</p>
        </div>
      ) : (
        <p>Loading user details...</p>
      )}
    </div>
  );
};

export default Dashboard;

import './App.css';
import SignUp from './pages/SignUp';
import Login from './pages/Login.jsx';
import Home from './pages/Home';
import NotFound from './pages/NotFound';
import OnBoarding from './pages/onBoarding.jsx';

import { BrowserRouter, Routes, Route } from 'react-router-dom';

function App() {
  return (
    <BrowserRouter>
      <Routes>
      <Route path="/" element={<OnBoarding />} />
        <Route path="/home" element={<Home />} />
        <Route path="/login" element={<Login />} />
        <Route path="/sign-up" element={<SignUp />} />
        <Route path="*" element={<NotFound />} />
      </Routes>
    </BrowserRouter>
  );
}

export default App;

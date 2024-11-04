// App.js
import React from 'react';
import { Routes, Route } from 'react-router-dom'
import Main from './components/pages/Main.js';
import Login from './components/pages/Login.js';
import Profile from './components/pages/Profile.js';
import Nav from "./components/pages/Nav.js"
import Callback from './components/pages/Callback.js'
import { AuthProvider } from './components/auth/AuthContext.tsx'
import './App.css';

function App() {
  return (
    <div className="app">
      <AuthProvider>
        <Nav />
        <Routes>
          <Route path="/" element={<Main />} />
          <Route path="/account/login" element={<Login />} />
          <Route path="/account/profile" element={<Profile />} />
          <Route path="/callback" element={<Callback />} />
        </Routes>
      </AuthProvider>
    </div>
  );
}

export default App;
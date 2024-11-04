// App.js
import React from 'react';
import { Routes, Route } from 'react-router-dom'
import Main from './components/Main';
import Login from './components/Login';
import Profile from './components/Profile';
import Nav from "./components/Nav"
import Callback from './components/Callback'
import { AuthProvider } from './components/authContext.tsx'
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
// App.js
import React from 'react';
import { Routes, Route } from 'react-router-dom'
import TweetForm from './components/TweetForm';
import Login from './components/Login';
import Nav from "./components/Nav"
import Callback from './components/Callback'
import './App.css';

function App() {
  return (
    <div className="app">
      <Nav />
      <Routes>        
        <Route path="/" element={<TweetForm />} />
        <Route path="/account/login" element={<Login />} />
        <Route path="/callback" element={<Callback />} />
      </Routes>
    </div>
  );
}

export default App;
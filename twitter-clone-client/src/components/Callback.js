import React, { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import CallbackStyles from "./Callback.module.css"

function Callback() {
  const navigate = useNavigate();
  const [idToken, setIdToken] = useState('');

  useEffect(() => {
    console.log('Callback component mounted');
    const search = window.location.search;
    const query = new URLSearchParams(search);
    const idToken = query.get('id_token');
    const refreshToken = query.get('refresh_token');

    if (idToken) {
      setIdToken(idToken);
      // Handle successful login response
      console.log(`Organization login successful. idToken: ${idToken}\nrefreshToken: ${refreshToken}`);
    } else {
      // Handle error or invalid token
      console.error('ID token not found in the callback URL');
    }
  }, [navigate]);

  return (        
      <div className={CallbackStyles.container}>
          Frontend redirect.
      </div>
  );
};

export default Callback;
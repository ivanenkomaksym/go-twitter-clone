import React, { useState, useEffect } from "react";
import LoginStyles from "../../styles/pages/Login.module.css"
import { loginAuthorizeUrl } from "../../common.js"
import { GoogleLogin } from '@react-oauth/google';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../auth/AuthContext.tsx';

function Login() {
    const [clicked, setClicked] = useState(false);
    useEffect(() => {
        if (clicked) {
            window.location.assign(loginAuthorizeUrl);
        }
    });

    const { checkAuth } = useAuth();
    const navigate = useNavigate();

    return (
        <div className={LoginStyles.loginContainer}>
            <div className={LoginStyles.loginContainerv2}>
                <h1>Welcome back</h1>
                <button onClick={() => setClicked(true)} className={LoginStyles.loginBTN}>
                    Login
                </button>
                <GoogleLogin
                    onSuccess={credentialResponse => {
                        console.log(credentialResponse);
                        localStorage.setItem('authToken', credentialResponse.credential);
                        checkAuth();
                        navigate("/"); // Redirect to root
                    }}

                    onError={() => {
                        console.log('Login Failed');
                    }}
                    useOneTap
                    auto_select
                />
            </div>
        </div>
    )
}

export default Login;
import React, { useState, useEffect } from "react";
import LoginStyles from "../../styles/pages/Login.module.css"
import {loginAuthorizeUrl} from "../../config.js"

function Login() {    
    const [clicked, setClicked] = useState(false);
    useEffect(() => {
        if (clicked) {
            window.location.assign(loginAuthorizeUrl);
        }
    });

    return (
        <div className={LoginStyles.loginContainer}>
            <div className={LoginStyles.loginContainerv2}>
                <h1>Welcome back</h1>
                <button onClick={() => setClicked(true)} className={LoginStyles.loginBTN}>
                    Login
                </button>
            </div>
        </div>
    )
}

export default Login;
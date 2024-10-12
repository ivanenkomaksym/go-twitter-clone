import React, { useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import LoginStyles from "./Login.module.css"

function Login() {
    const handleLoginClick = () => {
        console.log(`Send request to login.`);
    };

    return (
        <div className={LoginStyles.loginContainer}>
            <div className={LoginStyles.loginContainerv2}>
                <h1>Welcome back</h1>
                <button className={LoginStyles.loginBTN} onClick={handleLoginClick}>
                    Login
                </button>
            </div>
        </div>
    )
}

export default Login;
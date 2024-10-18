import React, { useEffect } from 'react';
import CallbackStyles from "./Callback.module.css"
import { signIn } from "../authhandlers.js"
import { useNavigate } from 'react-router-dom';

function Callback() {

    useEffect(() => {
        const signInAsync = async () => {
            await signIn();
            navigate("/"); // Redirect to root
        };

        signInAsync();
    }, []);

    const navigate = useNavigate();

    return (
        <div className={CallbackStyles.container}>
            Frontend redirect.
        </div>
    );
};

export default Callback;
import React, { useEffect } from 'react';
import CallbackStyles from "./Callback.module.css"
import { useNavigate } from 'react-router-dom';
import { useAuth } from './authContext.tsx';

function Callback() {

    const { checkAuth } = useAuth();

    useEffect(() => {
        checkAuth();
        navigate("/"); // Redirect to root
    }, []);

    const navigate = useNavigate();

    return (
        <div className={CallbackStyles.container}>
            Frontend redirect.
        </div>
    );
};

export default Callback;
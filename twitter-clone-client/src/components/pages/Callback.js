import React, { useEffect } from 'react';
import CallbackStyles from "../../styles/pages/Callback.module.css"
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../auth/AuthContext.tsx';

function Callback() {
    const { checkAuth } = useAuth();
    const navigate = useNavigate();

    useEffect(() => {
        checkAuth();
        navigate("/"); // Redirect to root
    }, [checkAuth, navigate]);

    return (
        <div className={CallbackStyles.container}>
            Frontend redirect.
        </div>
    );
};

export default Callback;
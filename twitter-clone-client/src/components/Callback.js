import React, { useEffect } from 'react';
import CallbackStyles from "./Callback.module.css"
import { signIn } from "../authhandlers.js"

function Callback() {

    useEffect(() => {
        const signInAsync = async () => {
            await signIn();
        };

        signInAsync();
    }, []);

    return (
        <div className={CallbackStyles.container}>
            Frontend redirect.
        </div>
    );
};

export default Callback;
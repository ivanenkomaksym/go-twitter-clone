import React from "react";
import {Link} from "react-router-dom"
import {useEffect, useState} from "react";
import { useLocation } from 'react-router-dom'
import NavStyles from "./Nav.module.css"
import {loadUserFromLocalStorage} from "../authhandlers.js"
import { logOutAuthorizeUrl } from "../config.js"
import IsAuthnEnabled from "../useFeatureFlags.js";

function Nav() {
    const [userData, setUserData] = useState(null);
    const location = useLocation();

    useEffect(() => {
        const localUser = loadUserFromLocalStorage();
        if (localUser) {
            setUserData(localUser);
        } else {
            setUserData(null);
        }
    }, [location]);

    function handleLogOut(e) {
        e.preventDefault();
        window.location.assign(logOutAuthorizeUrl);
    }

    return (
        <nav className={NavStyles.mainNav}>
            <div className={NavStyles.navContainer}>
                <Link className={NavStyles.linkBTN} to="/">Home</Link>
            </div>
            <div>
                {IsAuthnEnabled() ?
                (
                    <>
                    {userData ? (
                        <div className={NavStyles.rightSideNav}>
                            <i className="fa-solid fa-user"></i>
                            <div>
                                <span className={NavStyles.accountText}>Account</span>
                                <div className={NavStyles.navContainer}>
                                    <Link className={NavStyles.linkBTN} to="/account/profile">Profile</Link>
                                    <span className={NavStyles.orText}>or</span>
                                    <Link onClick={handleLogOut} className={NavStyles.linkBTN} to="/">Logout</Link>
                                </div>
                            </div>
                        </div>
                    ) : (
                        <div className={NavStyles.rightSideNav}>
                            <i className="fa-solid fa-user"></i>
                            <div>
                                <span className={NavStyles.accountText}>Account</span>
                                <div className={NavStyles.navContainer}>
                                    <Link className={NavStyles.linkBTN} to="/account/login">Login</Link>
                                    <span className={NavStyles.orText}>or</span>
                                    <Link className={NavStyles.linkBTN} to="/account/signup">Signup</Link>
                                </div>
                            </div>
                        </div>
                    )}
                    </>
                )
                    : <p />
                }
            </div>
        </nav>
    );
}

export default Nav;
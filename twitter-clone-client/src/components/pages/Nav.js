import React from "react";
import { Link } from "react-router-dom"
import NavStyles from "../../styles/pages/Nav.module.css"
import { logOutAuthorizeUrl } from "../../config.js"
import { useAuth } from '../auth/AuthContext.tsx';

function Nav() {
    const { isAuthenticated } = useAuth();

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
                {isAuthenticated ? (
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
            </div>
        </nav>
    );
}

export default Nav;
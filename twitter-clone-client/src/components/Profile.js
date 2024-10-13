import React, { useEffect, useState } from 'react';
import ProfileStyles from "./Profile.module.css"
import {loadUserFromLocalStorage} from "../authhandlers.js"

function Profile() {
  const [userInfo, setUserInfo] = useState(null);

  // Retrieve user information from local storage
  useEffect(() => {
    const userInfoFromLocalStorage = loadUserFromLocalStorage()
    setUserInfo(userInfoFromLocalStorage);
  }, []);

  return (
    <div className={ProfileStyles.profileContainer}>
      {userInfo ? (
        <div>
          <h2>User Profile</h2>
          <table className={ProfileStyles.profileTable}>
            <tbody>
              <tr>
                <td>First Name:</td>
                <td>{userInfo.firstName}</td>
              </tr>
              <tr>
                <td>Last Name:</td>
                <td>{userInfo.lastName}</td>
              </tr>
              <tr>
                <td>Email:</td>
                <td>{userInfo.email}</td>
              </tr>
              <tr>
                <td>Profile Picture:</td>
                <td><img src={userInfo.picture} alt="User Profile" className={ProfileStyles.profilePicture} /></td>
              </tr>
            </tbody>
          </table>
        </div>
      ) : (
        <div>
          <p>No user information available.</p>
        </div>
      )}
    </div>
  );
}

export default Profile;
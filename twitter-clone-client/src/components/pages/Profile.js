import ProfileStyles from "../../styles/pages/Profile.module.css"
import { useAuth } from '../auth/AuthContext.tsx';

function Profile() {
  const { user } = useAuth();

  return (
    <div className={ProfileStyles.profileContainer}>
      {user ? (
        <div>
          <h2>User Profile</h2>
          <table className={ProfileStyles.profileTable}>
            <tbody>
              <tr>
                <td>First Name:</td>
                <td>{user.firstName}</td>
              </tr>
              <tr>
                <td>Last Name:</td>
                <td>{user.lastName}</td>
              </tr>
              <tr>
                <td>Email:</td>
                <td>{user.email}</td>
              </tr>
              <tr>
                <td>Profile Picture:</td>
                <td><img src={user.picture} alt="User Profile" className={ProfileStyles.profilePicture} /></td>
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
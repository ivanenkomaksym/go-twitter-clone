import axios from "axios"
import User from "./models/user"
import { userInfoUrl } from "./config.js"

export async function fetchtUserInfo() {
    const instance = axios.create({
        withCredentials: true,
    });

    const response = await instance.get(userInfoUrl)
    if (response.status != 200) {
        clearUserFromLocalStorage(); // Clear user data on 401 Unauthorized
        return null;
    }

    const data = response.data

    const user = new User(
        data.firstName,         // First name from the response
        data.lastName,          // Last name
        data.email,             // Email
        data.picture,           // Profile picture
        document.cookie         // Use document.cookie if you need to extract id_token
    );

    return user;
}

export async function signIn() {
    console.log("auth.signIn");

    try {
        var user = await fetchtUserInfo();

        // Save user to local storage
        saveUserToLocalStorage(user);
    }
    catch (err) {
        clearUserFromLocalStorage(); // Clear user data on 401 Unauthorized
        console.log(err);
    }
}

export function saveUserToLocalStorage(user) {
    localStorage.setItem('user_info', JSON.stringify({ ...user }));
}

export function loadUserFromLocalStorage() {
    return JSON.parse(localStorage.getItem("user_info"));
}

export function clearUserFromLocalStorage() {
    localStorage.removeItem("user_info");
}
import axios from "axios"
import User from "./models/user"
import { userInfoUrl } from "./config.js"

export async function signIn() {
    console.log("auth.signIn");

    try {
        const instance = axios.create({
            withCredentials: true,
        });

        const response = await instance.get(userInfoUrl)
        const data = response.data

        const user = new User(
            data.firstName,         // First name from the response
            data.lastName,          // Last name
            data.email,             // Email
            data.picture,           // Profile picture
            document.cookie         // Use document.cookie if you need to extract id_token
        );

        // Save user to local storage
        saveUserToLocalStorage(user);
    }
    catch (err) {
        console.log(err);
    }
}

export function saveUserToLocalStorage(user) {
    localStorage.setItem('user_info', JSON.stringify({ ...user }));
}

export function loadUserFromLocalStorage() {
    return JSON.parse(localStorage.getItem("user_info"));
}
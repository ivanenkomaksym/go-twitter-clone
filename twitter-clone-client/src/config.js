const featureFlags = {
    authn: false
};

export default featureFlags;

export const loginAuthorizeUrl = "http://localhost:8016/auth/google/login";
export const logOutAuthorizeUrl = "http://localhost:8016/auth/google/logout";
export const userInfoUrl = "http://localhost:8016/auth/google/userinfo";
const config = {
    applicationUrl: "http://localhost:3000",
  };
  
  export default config;

  export const loginAuthorizeUrl = config.applicationUrl + "/auth/google/login";
  export const logOutAuthorizeUrl = config.applicationUrl + "/auth/google/logout";
  export const userInfoUrl = config.applicationUrl + "/auth/google/userinfo";
  export const feedsUrl = config.applicationUrl + '/api/feeds';
  export const tweetsUrl = config.applicationUrl + '/api/tweets';
  
  export const getTaggedFeedsUrl = (tag) => `${feedsUrl}/${tag}`
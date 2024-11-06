console.log('NODE_ENV:', process.env.NODE_ENV);
const environment = process.env.NODE_ENV || 'development'; // Default to development if NODE_ENV is not set
const config = require('./env.json')[environment];

console.log('Current Environment:', environment);
console.log('Configuration:', config);

export default config;

export const loginAuthorizeUrl = config.applicationUrl + "/auth/google/login";
export const logOutAuthorizeUrl = config.applicationUrl + "/auth/google/logout";
export const userInfoUrl = config.applicationUrl + "/auth/google/userinfo";
export const feedsUrl = config.applicationUrl + '/api/feeds';
export const tweetsUrl = config.applicationUrl + '/api/tweets';

export const getTaggedFeedsUrl = (tag) => `${feedsUrl}/${tag}`
console.log('NODE_ENV:', process.env.NODE_ENV);
console.log('REACT_APP_CLIENT_APPLICATIONURL:', process.env.REACT_APP_CLIENT_APPLICATIONURL);
const environment = process.env.NODE_ENV || 'development'; // Default to development if NODE_ENV is not set
const config = {
    applicationUri: process.env.REACT_APP_CLIENT_APPLICATIONURL  || require('./env.json')[environment].applicationUri
};

console.log('Current Environment:', environment);
console.log('Configuration:', config);

export default config;

export const loginAuthorizeUrl = config.applicationUri + "/auth/google/login";
export const logOutAuthorizeUrl = config.applicationUri + "/auth/google/logout";
export const userInfoUrl = config.applicationUri + "/auth/google/userinfo";
export const feedsUrl = config.applicationUri + '/api/feeds';
export const tweetsUrl = config.applicationUri + '/api/tweets';

export const getTaggedFeedsUrl = (tag) => `${feedsUrl}/${tag}`
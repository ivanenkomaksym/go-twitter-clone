console.log('NODE_ENV:', process.env.NODE_ENV);
const environment = process.env.NODE_ENV || 'development'; // Default to development if NODE_ENV is not set
const config = require('./env.json')[environment];

console.log('Current Environment:', environment);
console.log('Configuration:', config);

export default config;
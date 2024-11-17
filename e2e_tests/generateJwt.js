import jwt from 'jsonwebtoken';

// Load the service account key from environment variables
const serviceAccountKey = JSON.parse(process.env.GOOGLE_SERVICE_ACCOUNT_KEY);
console.log("serviceAccountKey: ", serviceAccountKey);

// Extract the necessary values from the service account key JSON
const clientEmail = serviceAccountKey.client_email;
const privateKey = serviceAccountKey.private_key;

/**
 * Generate a JWT signed with RS256 for Google OAuth2.0 service account authentication.
 * @returns {string} - Signed JWT token.
 */
export function createJWT() {
    const now = Math.floor(Date.now() / 1000);
    const payload = {
        iss: clientEmail,                       // Issuer - service account email
        scope: "https://www.googleapis.com/auth/cloud-platform", // Scope(s) you need
        aud: "https://oauth2.googleapis.com/token",              // Audience - Google's OAuth2 token URL
        exp: now + 3600,                        // Expiration time (1 hour)
        iat: now                                // Issued at time
    };

    // Sign the token with RS256 algorithm
    const token = jwt.sign(payload, privateKey, { algorithm: 'RS256' });
    return token;
}
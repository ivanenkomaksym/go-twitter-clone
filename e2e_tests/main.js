import axios from 'axios'; // Use import for ES Modules
import { createJWT } from './generateJwt.js'; // Ensure the correct path
import { expect } from 'chai'; // Use import for ES Modules

describe('E2E Test', () => {
    let accessToken;

    before(async () => {
        // Step 1: Create JWT
        const jwtToken = createJWT();
        console.log("JWT Token: ", jwtToken);

        // Step 2: Request Access Token
        try {
            const response = await axios.post('https://oauth2.googleapis.com/token', null, {
                params: {
                    grant_type: 'urn:ietf:params:oauth:grant-type:jwt-bearer',
                    assertion: jwtToken
                },
                headers: {
                    'Content-Type': 'application/x-www-form-urlencoded'
                }
            });

            accessToken = response.data.access_token;
            console.log("Access Token: ", accessToken);
        } catch (error) {
            console.error("Error fetching access token: ", error.response.data);
            throw error; // Rethrow to fail the test
        }
    });

    it('test', async () => {
        try {
            const url = 'http://localhost:8016/api/tweets';

            // Act
            const response = await axios.get(url, {
                headers: {
                    Authorization: `${accessToken}`
                }
            });
            
            // Assert
            expect(response.status).to.equal(200);
        } catch (error) {
            // Print the error details
            console.error("Error making authenticated request: ", error.response ? error.response.data : error.message);
            throw error; // Rethrow to fail the test
        }
    });
});

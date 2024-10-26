const { createJWT } = require('./generateJwt');

describe('E2E Test', () => {
    before(async () => {
        const token = createJWT();
        console.log("token: ", token);
    });

    it('test', async () => {
        // Use jwtClient to make authorized requests
        
        // Assertions here
    });
});

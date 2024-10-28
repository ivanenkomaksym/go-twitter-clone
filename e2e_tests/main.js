import axios from 'axios';
import { createJWT } from './generateJwt.js';
import { expect } from 'chai';

describe('E2E Tweet API Tests', () => {
    let accessToken;
    const golangtag = 'golang';
    const e2etestingtag = 'e2etesting';
    let createdTweetId = "ca17e472-4d75-4bd5-9f75-4f9e61892591";

    before(async () => {
        // Generate JWT and fetch access token
        const jwtToken = createJWT();

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
        } catch (error) {
            console.error("Error fetching access token: ", error.response.data);
            throw error;
        }
    });

    it('should create a new tweet', async () => {
        const url = 'http://localhost:8016/api/tweets';
        const newTweet = {
            id: createdTweetId,
            title: 'e2e tests on GO app in CI',
            content: 'How to run end-to-end tests on full scale GO environment during CI using github workflow actions',
            author: 'testUser',
            tags: [golangtag, e2etestingtag],
            created_at: new Date().toISOString()
        };

        try {
            const response = await axios.post(url, newTweet, {
                headers: {
                    Authorization: `${accessToken}`
                }
            });

            expect(response.status).to.equal(201);
        } catch (error) {
            console.error("Error creating tweet: ", error.response ? error.response.data : error.message);
            throw error;
        }
    });

    it('should retrieve all tweets', async () => {
        const url = 'http://localhost:8016/api/tweets';

        try {
            const response = await axios.get(url, {
                headers: {
                    Authorization: `${accessToken}`
                }
            });

            expect(response.status).to.equal(200);
            expect(response.data).to.be.an('array');
            expect(response.data.some(tweet => tweet.id === createdTweetId)).to.be.true;
        } catch (error) {
            console.error("Error retrieving all tweets: ", error.response ? error.response.data : error.message);
            throw error;
        }
    });

    it('should retrieve all feeds and verify the tagged feed exists', async () => {
        const url = 'http://localhost:8016/api/feeds';

        try {
            const response = await axios.get(url, {
                headers: {
                    Authorization: `${accessToken}`
                }
            });

            expect(response.status).to.equal(200);
            expect(response.data.feeds).to.be.an('array');

            const golangfeedExists = response.data.feeds.some(feed => feed.name === golangtag);
            expect(golangfeedExists).to.be.true;
            
            const e2etestingfeedExists = response.data.feeds.some(feed => feed.name === e2etestingtag);
            expect(e2etestingfeedExists).to.be.true;
        } catch (error) {
            console.error("Error retrieving all feeds: ", error.response ? error.response.data : error.message);
            throw error;
        }
    });

    it('should retrieve a specific feed by tag name and include the created tweet', async () => {
        const url = `http://localhost:8016/api/feeds/${golangtag}`;

        try {
            const response = await axios.get(url, {
                headers: {
                    Authorization: `${accessToken}`
                }
            });

            expect(response.status).to.equal(200);
            expect(response.data).to.have.property('name', golangtag);
            expect(response.data).to.have.property('tweets').that.is.an('array');

            // Confirm that the specific tweet exists within the feed's tweets
            const tweetExistsInFeed = response.data.tweets.some(tweet => tweet.id === createdTweetId);
            expect(tweetExistsInFeed).to.be.true;
        } catch (error) {
            console.error("Error retrieving feed by name: ", error.response ? error.response.data : error.message);
            throw error;
        }
    });

    it('should retrieve a specific tweet by ID', async () => {
        const url = `http://localhost:8016/api/tweets/${createdTweetId}`;

        try {
            const response = await axios.get(url, {
                headers: {
                    Authorization: `${accessToken}`
                }
            });

            expect(response.status).to.equal(200);
            expect(response.data).to.have.property('id', createdTweetId);
            expect(response.data).to.have.property('title', 'My First Tweet');
        } catch (error) {
            console.error("Error retrieving tweet by ID: ", error.response ? error.response.data : error.message);
            throw error;
        }
    });

    it('should delete a tweet by ID', async () => {
        const url = `http://localhost:8016/api/tweets/${createdTweetId}`;

        try {
            const response = await axios.delete(url, {
                headers: {
                    Authorization: `${accessToken}`
                }
            });

            expect(response.status).to.equal(204);
        } catch (error) {
            console.error("Error deleting tweet: ", error.response ? error.response.data : error.message);
            throw error;
        }
    });

    it('should verify tweet deletion', async () => {
        const url = `http://localhost:8016/api/tweets/${createdTweetId}`;

        try {
            await axios.get(url, {
                headers: {
                    Authorization: `${accessToken}`
                }
            });
            throw new Error("Tweet was not deleted"); // Fail if tweet is found
        } catch (error) {
            expect(error.response.status).to.equal(404);
        }
    });
});

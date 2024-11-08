import axios from 'axios';
import { createJWT } from './generateJwt.js';
import { expect } from 'chai';

describe('E2E Tweet API Tests', () => {
    let accessToken;
    const baseurl = 'http://localhost:8016/api';
    const golangtag = 'golang';
    const e2etestingtag = 'e2etesting';
    const tweettitle = 'e2e tests on GO app in CI';
    let createdTweetId = "";

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

    it('should return unauthorized given no authorization', async () => {
        const url = `${baseurl}/tweets`;

        try {
            await axios.get(url);
        } catch (error) {
            expect(error.response.status).to.equal(401);
        }
    });

    it('should create a new tweet', async () => {
        const url = `${baseurl}/tweets`;
        const newTweetRequest = {
            title: tweettitle,
            content: 'How to run end-to-end tests on full scale GO environment during CI using github workflow actions',
            tags: [golangtag, e2etestingtag]
        };

        try {
            const response = await axios.post(url, newTweetRequest, {
                headers: {
                    Authorization: `${accessToken}`
                }
            });

            expect(response.status).to.equal(201);

            // Verify the response data
            expect(response.data).to.be.an('object');
            expect(response.data).to.have.property('id').that.is.a('string');
            createdTweetId = response.data.id;
            
            expect(response.data).to.have.property('title').that.equals(newTweetRequest.title);
            expect(response.data).to.have.property('content').that.equals(newTweetRequest.content);
            expect(response.data).to.have.property('tags').that.is.an('array').that.deep.equals(newTweetRequest.tags);
            expect(response.data).to.have.property('created_at').that.is.a('string'); // Check format if needed
            expect(response.data).to.have.property('user').that.is.an('object');
        } catch (error) {
            console.error("Error creating tweet: ", error.response ? error.response.data : error.message);
            throw error;
        }
    });

    it('should retrieve all tweets', async () => {
        const url = `${baseurl}/tweets`;

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
        const url = `${baseurl}/feeds`;

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
        const url = `${baseurl}/feeds/${golangtag}`;

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
        const url = `${baseurl}/tweets/${createdTweetId}`;

        try {
            const response = await axios.get(url, {
                headers: {
                    Authorization: `${accessToken}`
                }
            });

            expect(response.status).to.equal(200);
            expect(response.data).to.have.property('id', createdTweetId);
            expect(response.data).to.have.property('title', tweettitle);
        } catch (error) {
            console.error("Error retrieving tweet by ID: ", error.response ? error.response.data : error.message);
            throw error;
        }
    });

    it('should delete a tweet by ID', async () => {
        const url = `${baseurl}/tweets/${createdTweetId}`;

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
        const url = `${baseurl}/tweets/${createdTweetId}`;

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

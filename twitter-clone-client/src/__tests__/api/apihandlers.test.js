jest.mock('../../common', () => require('../../__mocks__/configMock'));

import axios from 'axios';
import { fetchUserInfo, fetchTagsFromServer, addTweetToServer, fetchTaggedTweets } from '../../api/apihandlers.js'
import config, { feedsUrl, tweetsUrl, userInfoUrl } from '../../common';

// Mock axios
jest.mock('axios');

describe('apiHandlers', () => {
  afterEach(() => {
    jest.clearAllMocks();
  });

  describe('fetchUserInfo', () => {
    it('should return user data when the response status is 200', async () => {
      const mockData = { name: 'John Doe', email: 'john@example.com' };
      axios.create.mockReturnValue({
        get: jest.fn().mockResolvedValue({ status: 200, data: mockData }),
      });

      const result = await fetchUserInfo();
      expect(result).toEqual(mockData);
      expect(axios.create().get).toHaveBeenCalledWith(`${config.applicationUri}${userInfoUrl}`);
    });

    it('should return null when the response status is not 200', async () => {
      axios.create.mockReturnValue({
        get: jest.fn().mockResolvedValue({ status: 404 }),
      });

      const result = await fetchUserInfo();
      expect(result).toBeNull();
    });

    it('should return null and log an error when the request fails', async () => {
      console.error = jest.fn();
      axios.create.mockReturnValue({
        get: jest.fn().mockRejectedValue(new Error('Request failed')),
      });

      const result = await fetchUserInfo();
      expect(result).toBeNull();
      expect(console.error).toHaveBeenCalledWith('Error:', expect.any(Error));
    });
  });

  describe('fetchTagsFromServer', () => {
    it('should return formatted tags data when response status is 200', async () => {
      const mockData = { feeds: [{ name: 'tag1', tweets: 10 }, { name: 'tag2', tweets: 5 }] };
      axios.create.mockReturnValue({
        get: jest.fn().mockResolvedValue({ status: 200, data: mockData }),
      });

      const result = await fetchTagsFromServer();
      expect(result).toEqual([
        { name: 'tag1', nofTweets: 10 },
        { name: 'tag2', nofTweets: 5 },
      ]);
    });

    it('should return an empty array when response status is not 200', async () => {
      axios.create.mockReturnValue({
        get: jest.fn().mockResolvedValue({ status: 404 }),
      });

      const result = await fetchTagsFromServer();
      expect(result).toEqual([]);
    });

    it('should return an empty array and log an error when the request fails', async () => {
      console.error = jest.fn();
      axios.create.mockReturnValue({
        get: jest.fn().mockRejectedValue(new Error('Request failed')),
      });

      const result = await fetchTagsFromServer();
      expect(result).toEqual([]);
      expect(console.error).toHaveBeenCalledWith('Error:', expect.any(Error));
    });
  });

  describe('addTweetToServer', () => {
    it('should return true when the tweet is added successfully', async () => {
      axios.post.mockResolvedValue({ status: 201 });

      const result = await addTweetToServer({ title: 'Title', content: 'Content', author: 'Author' }, ['tag1']);
      expect(result).toBe(true);
      expect(axios.post).toHaveBeenCalledWith(
        `${config.applicationUri}${tweetsUrl}`,
        expect.objectContaining({
          title: 'Title',
          content: 'Content',
          author: 'Author',
          tags: ['tag1'],
        }),
        expect.objectContaining({
          withCredentials: true,
          headers: { 'Content-Type': 'application/json' },
        })
      );
    });

    it('should return false when the response status is not 201', async () => {
      axios.post.mockResolvedValue({ status: 400 });

      const result = await addTweetToServer({ title: 'Title', content: 'Content', author: 'Author' }, ['tag1']);
      expect(result).toBe(false);
    });

    it('should return false and log an error when the request fails', async () => {
      console.error = jest.fn();
      axios.post.mockRejectedValue(new Error('Request failed'));

      const result = await addTweetToServer({ title: 'Title', content: 'Content', author: 'Author' }, ['tag1']);
      expect(result).toBe(false);
      expect(console.error).toHaveBeenCalledWith('Error:', expect.any(Error));
    });
  });

  describe('fetchTaggedTweets', () => {
    it('should set tagged tweets and return true when response status is 200', async () => {
      const mockTweets = [{ id: 1, content: 'Tweet 1' }];
      const setTaggedTweets = jest.fn();
      axios.create.mockReturnValue({
        get: jest.fn().mockResolvedValue({ status: 200, data: { tweets: mockTweets } }),
      });

      const result = await fetchTaggedTweets('tag1', setTaggedTweets);
      expect(result).toBe(true);
      expect(setTaggedTweets).toHaveBeenCalledWith(mockTweets);
    });

    it('should return false and set empty array when response status is not 200', async () => {
      const setTaggedTweets = jest.fn();
      axios.create.mockReturnValue({
        get: jest.fn().mockResolvedValue({ status: 404 }),
      });

      const result = await fetchTaggedTweets('tag1', setTaggedTweets);
      expect(result).toBe(false);
      expect(setTaggedTweets).toHaveBeenCalledWith([]);
    });

    it('should return false and log an error when the request fails', async () => {
      console.error = jest.fn();
      const setTaggedTweets = jest.fn();
      axios.create.mockReturnValue({
        get: jest.fn().mockRejectedValue(new Error('Request failed')),
      });

      const result = await fetchTaggedTweets('tag1', setTaggedTweets);
      expect(result).toBe(false);
      expect(console.error).toHaveBeenCalledWith('Error:', expect.any(Error));
    });
  });
});

import axios from "axios"
import { v4 as uuidv4 } from 'uuid';
import config from '../common';

// Function to fetch tags from the server
export const fetchTagsFromServer = async () => {
  try {
    const instance = axios.create({
      withCredentials: true,
    });
    const response = await instance.get(config.applicationUri + '/api/feeds');

    if (response.status == 200) {
      const data = await response.data;
      const fetchedTags = data.feeds.map(feed => ({
        name: feed.name,
        nofTweets: feed.tweets
      }));
      return fetchedTags;
    } else {
      console.error('Failed to fetch tags.');
      return [];
    }
  } catch (error) {
    console.error('Error:', error);
    return [];
  }
};

export const addTweetToServer = async (formData, tweetTags) => {
  const currentDate = new Date().toISOString();

  try {
    const response = await axios.post(`${config.applicationUri}/api/tweets`, {
      id: uuidv4(),
      title: formData.title,
      content: formData.content,
      author: formData.author,
      tags: tweetTags,
      created_at: currentDate,
      likes: null
    }, {
      withCredentials: true,
      headers: {
        'Content-Type': 'application/json'
      }
    });

    if (response.status == 201) {
      console.log('Tweet added successfully!');
      return true;
    } else {
      console.error('Failed to add tweet.');
      return false;
    }
  } catch (error) {
    console.error('Error:', error);
    return false;
  }
};

export const fetchTaggedTweets = async (tag, setTaggedTweets) => {
  try {
    const instance = axios.create({
      withCredentials: true,
    });
    const response = await instance.get(config.applicationUri + `/api/feeds/${tag}`);

    if (response.status == 200) {
      const feedData = await response.data;
      setTaggedTweets(feedData.tweets);
      return true;
    } else {
      console.error('Failed to fetch feed.');
      setTaggedTweets([]);
      return false;
    }
  } catch (error) {
    console.error('Error:', error);
    return false;
  }
};
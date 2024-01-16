import { v4 as uuidv4 } from 'uuid';
import config from './common';

// Function to fetch tags from the server
export const fetchTagsFromServer = async () => {
    try {
      const response = await fetch(config.applicationUri + '/api/feeds');
  
      if (response.ok) {
        const data = await response.json();
        const fetchedTags = data.feeds.map(feed => feed.name);
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
      const response = await fetch(config.applicationUri + '/api/tweets', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({
          id: uuidv4(),
          title: formData.title,
          content: formData.content,
          author: formData.author,
          tags: tweetTags,
          created_at: currentDate,
          likes: null
        })
      });
  
      if (response.ok) {
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
      const response = await fetch(config.applicationUri + `/api/feeds/${tag}`);
  
      if (response.ok) {
        const feedData = await response.json();
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
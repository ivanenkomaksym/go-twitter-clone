import React, { useState, useEffect } from 'react';
import { v4 as uuidv4 } from 'uuid';
import './TweetForm.css';

const TweetForm = () => {
  // Load hashtags from local storage on component mount
  const [formData, setFormData] = useState({
    title: '',
    content: '',
    author: ''
  });

  const [tags, setTags] = useState([]);
  const [tweetTags, setTweetTags] = useState([]);
  const [selectedTag, setSelectedTag] = useState(null);
  const [taggedTweets, setTaggedTweets] = useState([]);
  const [eventSource, setEventSource] = useState(null);

  // Effect to fetch tags from server on component mount
  useEffect(() => {
    const fetchTags = async () => {
      try {
        const response = await fetch('http://localhost:8016/api/feeds');
        
        if (response.ok) {
          const data = await response.json();
          const fetchedTags = data.feeds.map(feed => feed.name);
          setTags(fetchedTags);
        } else {
          console.error('Failed to fetch tags.');
        }
      } catch (error) {
        console.error('Error:', error);
      }
    };

    fetchTags();
  }, []);

  // Effect to set up EventSource for data changes
  useEffect(() => {
    const eventSource = new EventSource('http://localhost:8016/api/feeds');

    eventSource.addEventListener('data', (event) => {
      try {
        const data = JSON.parse(event.data);
        const updatedTags = data.feeds.map(feed => feed.name);
        setTags(updatedTags);
      } catch (error) {
        console.error('Error parsing event data:', error);
      }
    });

    return () => {
      // Clean up EventSource when component is unmounted
      eventSource.close();
    };
  }, []);

  const handleAddTweet = async () => {
    const currentDate = new Date().toISOString();

    const contentTags = formData.content.match(/#[a-zA-Z0-9_]+/g) || [];
    const resultContentTags = contentTags.map(tag => tag.substring(1));

    // Update the tags state before calling setFormData
    setTags((prevTags) => {
      const uniqueTags = Array.from(new Set([...prevTags, ...resultContentTags]));
      return uniqueTags;
    });

    const tweetData = {
      id: uuidv4(),
      title: formData.title,
      content: formData.content,
      author: formData.author,
      tags: tweetTags,
      created_at: currentDate,
      likes: null
    };

    try {
      const response = await fetch('http://localhost:8016/api/tweets', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(tweetData)
      });

      if (response.ok) {
        console.log('Tweet added successfully!');
        // Clear the form after successful submission
        setFormData({
          title: '',
          content: ''
        });
      } else {
        console.error('Failed to add tweet.');
      }
    } catch (error) {
      console.error('Error:', error);
    }
  };

  const handleTagClick = async (tag) => {
    try {
      setSelectedTag(tag);
      
      const response = await fetch(`http://localhost:8016/api/feeds/${tag}`);
      if (response.ok) {
        const feedData = await response.json();
        setTaggedTweets(feedData.tweets);

        // Close existing EventSource connection if any
        if (eventSource) {
          eventSource.close();
        }

        // Set up a new EventSource connection
        const newEventSource = new EventSource(`http://localhost:8016/api/feeds/${tag}`);

        // Add event listener for the 'data' event
        newEventSource.addEventListener('data', (event) => {
          try {
            const data = JSON.parse(event.data);
            console.info('data EventSource', data)
            setTaggedTweets(data.tweets || []);
          } catch (error) {
            console.error('Error parsing event data:', error);
          }
        });

        // Save the new EventSource instance to close later
        setEventSource(newEventSource);
      } else {
        setTaggedTweets([]);
        console.error('Failed to fetch feed.');
      }
    } catch (error) {
      console.error('Error:', error);
    }
  };

  useEffect(() => {
    const contentTags = formData.content.match(/#[a-zA-Z0-9_]+/g) || [];
    setTweetTags(contentTags.map(tag => tag.substring(1)));
  }, [formData.content]);

  useEffect(() => {
    // Clean up EventSource when component is unmounted or tag changes
    return () => {
      if (eventSource) {
        eventSource.close();
      }
    };
  }, [selectedTag]);

  const handleInputChange = (e) => {
    const { name, value } = e.target;
    setFormData({
      ...formData,
      [name]: value
    });
  };

  return (
    <div className="tweet-form-container">
      <h3>Entered Hashtags:</h3>
          {tags.map((tag, index) => (
              <a href="#" onClick={() => handleTagClick(tag)}>{`#${tag} | `}</a>
          ))}

      <h2>Add new tweet</h2>
      <form className="tweet-form">
        <div className="form-group">
          <label>Title:</label>
          <input type="text" name="title" value={formData.title} onChange={handleInputChange} />
        </div>

        <div className="form-group">
          <label>Content:</label>
          <textarea name="content" value={formData.content} onChange={handleInputChange}></textarea>
        </div>

        <div className="form-group">
          <label>Author:</label>
          <input type="text" name="author" value={formData.author} onChange={handleInputChange} />
        </div>

        <button type="button" onClick={handleAddTweet}>Add tweet</button>
      </form>

      {selectedTag && (
        <div className="tagged-tweets-container">
          <h3>Tweets with #{selectedTag}</h3>
          <div className="tweet-boxes">
            {taggedTweets.map((tweet) => (
              <div key={tweet.id} className="tweet-box">
                <div className="tweet-title">
                  <strong>{tweet.title}</strong>
                </div>
                <div className="tweet-content">
                  {tweet.content}
                </div>
                <div className="tweet-author">
                  <span className="author-text">by {tweet.author}</span>
                </div>
              </div>
            ))}
          </div>
        </div>

      )}
    </div>
  );
};

export default TweetForm;

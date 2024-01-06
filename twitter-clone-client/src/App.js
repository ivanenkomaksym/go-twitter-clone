import React, { useState, useEffect } from 'react';
import { v4 as uuidv4 } from 'uuid';
import './TweetForm.css';

const TweetForm = () => {
  // Load hashtags from local storage on component mount
  const savedHashtags = JSON.parse(localStorage.getItem('hashtags')) || [];
  const [formData, setFormData] = useState({
    title: '',
    content: '',
    author: ''
  });

  const [tags, setTags] = useState(savedHashtags);
  const [selectedTag, setSelectedTag] = useState(null);
  const [taggedTweets, setTaggedTweets] = useState([]);

  const handleAddTweet = async () => {
    const currentDate = new Date().toISOString();

    const tweetData = {
      id: uuidv4(),
      title: formData.title,
      content: formData.content,
      author: formData.author,
      tags: tags,
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
          content: '',
          author: ''
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
      const response = await fetch(`http://localhost:8016/api/feeds/${tag}`);
      if (response.ok) {
        const feedData = await response.json();
        setSelectedTag(tag);
        setTaggedTweets(feedData.tweets);
      } else {
        console.error('Failed to fetch feed.');
      }
    } catch (error) {
      console.error('Error:', error);
    }
  };

  useEffect(() => {
    const contentTags = formData.content.match(/#[a-zA-Z0-9_]+/g) || [];
    setTags(contentTags.map(tag => tag.substring(1)));
  }, [formData.content]);

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
      <ul className="hashtags-list">
        {tags.map((tag, index) => (
          <li key={index}>
            <a href="#" onClick={() => handleTagClick(tag)}>{`#${tag}`}</a>
          </li>
        ))}
      </ul>

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
          <ul>
            {taggedTweets.map((tweet) => (
              <li key={tweet.id}>
                <strong>{tweet.title}</strong> - {tweet.content} by {tweet.author}
              </li>
            ))}
          </ul>
        </div>
      )}
    </div>
  );
};

export default TweetForm;

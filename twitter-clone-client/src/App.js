// Import necessary dependencies
import React, { useState, useEffect } from 'react';
import './TweetForm.css'; // Import the CSS file

// Define your functional component
const TweetForm = () => {
  // State for form data
  const [formData, setFormData] = useState({
    title: '',
    content: '',
    author: ''
  });

  // State for tags
  const [tags, setTags] = useState([]);

  // Function to handle form submission
  const handleAddTweet = async () => {
    // Create the JSON object to send in the POST request
    const tweetData = {
      id: '123',
      title: formData.title,
      content: formData.content,
      author: formData.author,
      tags: tags,
      created_at: '2024-01-01T00:00:00Z',
      likes: null
    };

    try {
      // Make the POST request
      const response = await fetch('http://localhost:8016/api/tweets', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(tweetData)
      });

      // Handle the response as needed
      if (response.ok) {
        console.log('Tweet added successfully!');
      } else {
        console.error('Failed to add tweet.');
      }
    } catch (error) {
      console.error('Error:', error);
    }
  };

  // Function to update tags based on content
  useEffect(() => {
    // Extract hashtags from content
    const contentTags = formData.content.match(/#[a-zA-Z0-9_]+/g) || [];
    setTags(contentTags.map(tag => tag.substring(1))); // Remove '#' from tags
  }, [formData.content]);

  // Function to handle form input changes
  const handleInputChange = (e) => {
    const { name, value } = e.target;
    setFormData({
      ...formData,
      [name]: value
    });
  };

  // JSX for the component
  return (
    <div className="tweet-form-container">
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
    </div>
  );
};

export default TweetForm;

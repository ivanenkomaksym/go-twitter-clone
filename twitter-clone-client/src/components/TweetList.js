// TweetList.js
import React from 'react';

const TweetList = ({ taggedTweets, selectedTag }) => {
  return (
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
  );
};

export default TweetList;

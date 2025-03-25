// TweetList.js
import React from 'react';
import TweetCard from './TweetCard';

const TweetList = ({ taggedTweets, handleTagClick }) => {
  return (
    <div className="tweet-list">
      {taggedTweets.map(tweet => (
        <TweetCard key={tweet.id} tweet={tweet} handleTagClick={handleTagClick} />
      ))}
    </div>
  );
};

export default TweetList;

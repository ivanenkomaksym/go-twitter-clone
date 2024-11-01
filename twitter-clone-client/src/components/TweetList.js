// TweetList.js
import React from 'react';
import TweetCard from './TweetCard';

const TweetList = ({ taggedTweets, selectedTag }) => {
  return (
    <div className="tweet-list">
      {taggedTweets.map(tweet => (
        <TweetCard key={tweet.id} tweet={tweet} />
      ))}
    </div>
  );
};

export default TweetList;

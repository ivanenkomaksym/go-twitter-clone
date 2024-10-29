// TagList.js
import React from 'react';
import './TagList.css';

const TagList = ({ tags, handleTagClick }) => {
  return (
    <div className="tag-list-sidebar">
      <h3>Trending Hashtags</h3>
      <ul>
        {tags.map((tag, index) => (
          <div  key={index} onClick={() => handleTagClick(tag)}>
            <div className="tag-item">
              <span className="tag-name">{`#${tag.name}`}</span>
              <span className="tag-count">{`${tag.nofTweets} tweets`}</span>
            </div>
          </div>
        ))}
      </ul>
    </div>
  );
};

export default TagList;

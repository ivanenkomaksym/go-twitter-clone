// TagList.js
import React from 'react';

const TagList = ({ tags, handleTagClick }) => {
  return (
    <div>
      <h3>Entered Hashtags:</h3>
      {tags.map((tag, index) => (
        <a href="#" key={index} onClick={() => handleTagClick(tag)}>{`#${tag} | `}</a>
      ))}
    </div>
  );
};

export default TagList;

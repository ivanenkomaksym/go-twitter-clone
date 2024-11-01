import React from 'react';
import './TweetCard.css';

const TweetCard = ({ tweet }) => {
    return (
        <div className="tweet-card">
            <div className="tweet-icon">
                <img src="https://pbs.twimg.com/profile_images/1527675441892544513/vXL1OQVh_400x400.jpg" alt={`${tweet.author}'s icon`} />
            </div>
            <div className="tweet-content">
                <div className="tweet-header">
                    <span className="tweet-author">{tweet.author}</span>
                    <span className="tweet-date">6h</span>
                </div>
                <div className="tweet-text">
                    {tweet.content.split(' ').map((word, index) => 
                        word.startsWith('#') ? <strong key={index} style={{ color: 'blue' }}>{word}</strong> : word + ' '
                    )}
                </div>
                <div className="tweet-actions">
                    <button className="tweet-action">
                        <i className="fas fa-comment"></i> {tweet.comments}
                    </button>
                    <button className="tweet-action">
                        <i className="fas fa-retweet"></i> {tweet.retweets}
                    </button>
                    <button className="tweet-action">
                        <i className="fas fa-heart"></i> {tweet.likes}
                    </button>
                    <button className="tweet-action">
                        <i className="fas fa-chart-bar"></i>
                    </button>
                    <button className="tweet-action">
                        <i className="fas fa-bookmark"></i>
                    </button>
                </div>
            </div>
        </div>
    );
};

export default TweetCard;

import React from 'react';
import './TweetCard.css';
import { formatDistanceToNow, parseISO } from 'date-fns';

const hashtagRegex = /(^|\B)#(?![0-9_]+\b)([a-zA-Z0-9_]{1,30})(\b|\r)/g; // Match hashtags
const renderedText = (text) => text.replace(hashtagRegex, (match, prefix, hash, tag) => {
    return `<a href="#" onclick="event.preventDefault();">${match}</a>`;
});

const createMarkup = (text) => {
    return { __html: renderedText(text) };
};

const TweetCard = ({ tweet, handleTagClick }) => {
    const createdAt = formatDistanceToNow(parseISO(tweet.created_at), { addSuffix: true });
    // const renderedText = replaceHashtagsWithLinks(tweet.content);

    return (
        <div className="tweet-card">
            <div className="tweet-icon">
                <img src={tweet.user.picture} alt={`${tweet.user.email}'s icon`} />
            </div>
            <div className="tweet-content">
                <div className="tweet-header">
                    <span className="tweet-user">{tweet.user.email}</span>
                    <span className="tweet-date">{createdAt}</span>
                </div>
                <div className="tweet-text" style={{ whiteSpace: 'pre-wrap' }}>
                    <div dangerouslySetInnerHTML={createMarkup(tweet.content)} onClick={(e) => {
                        if (e.target.tagName === 'A') {
                            const tagName = e.target.innerText.slice(1); // Remove the '#' character
                            console.log("clicked tagName:", tagName);
                            handleTagClick({name: tagName});
                        }
                    }} />
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

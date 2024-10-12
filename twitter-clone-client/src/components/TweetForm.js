// TweetForm.js
import React, { useState, useEffect } from 'react';
import { Link  } from 'react-router-dom';
import TagList from './TagList';
import TweetList from './TweetList';
import * as apiHandlers from '../apihandlers';
import * as eventSourceHandlers from '../eventSourceHandlers';
import './TweetForm.css';
import IsAuthnEnabled from '../useFeatureFlags.js';

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

    useEffect(() => {
        const fetchTags = async () => {
            const fetchedTags = await apiHandlers.fetchTagsFromServer();
            setTags(fetchedTags);
        };

        fetchTags();
    }, []);

    useEffect(() => {
        const eventSource = eventSourceHandlers.setUpFeedsEventSource(setTags);

        setEventSource(eventSource);

        return () => {
            if (eventSource) {
                eventSource.close();
            }
        };
    }, []);

    const handleAddTweet = async () => {
        const success = await apiHandlers.addTweetToServer(formData, tweetTags);

        if (success) {
            // Clear the form after successful submission
            setFormData({
                title: '',
                content: '',
                author: formData.author
            });
        }
    };

    const handleTagClick = async (tag) => {
        setSelectedTag(tag);
        const success = await apiHandlers.fetchTaggedTweets(tag, setTaggedTweets, setEventSource);

        if (success) {
            const eventSource = eventSourceHandlers.setUpFeedsTagEventSource(tag, setTaggedTweets);
            setEventSource(eventSource);

            return () => {
                if (eventSource) {
                    eventSource.close();
                }
            };
        }
    };

    useEffect(() => {
        const contentTags = formData.content.match(/#[a-zA-Z0-9_]+/g) || [];
        setTweetTags(contentTags.map(tag => tag.substring(1)));
    }, [formData.content]);

    useEffect(() => {
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

    const isAuthnEnabled = IsAuthnEnabled();

    const isAuthenticated = () => {
        if (!isAuthnEnabled)
            return true;

        const userInfo = JSON.parse(localStorage.getItem('user_info'));
        if (userInfo) {
            const { expirationDate } = userInfo;
            if (new Date(expirationDate) < new Date()) {
                console.log("Clearing local storage due to token expiration");
                localStorage.clear();
                return false;
            }
            return true;
        }
        return false;
    };

    return (
        <div className="container">
            {isAuthenticated() ? (
                <>
                    <div className="tweet-form-container">
                        <TagList tags={tags} handleTagClick={handleTagClick} />

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

                        {selectedTag && <TweetList taggedTweets={taggedTweets} selectedTag={selectedTag} />}
                    </div>
                </>
            ) : (
                <div className="welcomeMessage">
                    Welcome. <Link to="/account/login">Login</Link> to continue.
                </div>
            )};
        </div>
    );
}

export default TweetForm;

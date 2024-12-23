import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import TagList from '../tweet/TagList.js';
import TweetList from '../tweet/TweetList.js';
import * as apiHandlers from '../../api/apihandlers.js';
import * as eventSourceHandlers from '../../api/eventSourceHandlers.js';
import '../../styles/pages/Main.css';
import InputForm from '../tweet/InputForm.js';
import { useAuth } from '../auth/AuthContext.tsx';

function Main() {
    const [tags, setTags] = useState([]);
    const [tweetTags, setTweetTags] = useState([]);
    const [selectedTag, setSelectedTag] = useState(null);
    const [taggedTweets, setTaggedTweets] = useState([]);
    const [eventSource, setEventSource] = useState(null);

    const { isAuthenticated, user } = useAuth();

    const [formData, setFormData] = useState({
        title: '',
        content: '',
        author: user ? `${user.firstName} ${user.lastName}` : ''
    }, []);

    useEffect(() => {
        if (user) {
            setFormData(prevFormData => ({
                ...prevFormData,
                author: `${user.firstName} ${user.lastName}`
            }));
        }
    }, [user]);

    useEffect(() => {
        const fetchTags = async () => {
            const fetchedTags = await apiHandlers.fetchTagsFromServer();
            setTags(fetchedTags);
        };

        fetchTags();
    }, [selectedTag]);

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
        let tagName = tag.name;

        setSelectedTag(tagName);
        const success = await apiHandlers.fetchTaggedTweets(tagName, setTaggedTweets, setEventSource);

        if (success) {
            const eventSource = eventSourceHandlers.setUpFeedsTagEventSource(tagName, setTaggedTweets);
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
    }, [eventSource, selectedTag]);

    const handleInputChange = (e) => {
        const { name, value } = e.target;
        setFormData({
            ...formData,
            [name]: value
        });
    };

    return (
        <div className="container">
            <div className="main-panel">
                {isAuthenticated ? (
                    <>
                        <InputForm
                            user={user}
                            formData={formData}
                            handleInputChange={handleInputChange}
                            handleAddTweet={handleAddTweet}
                        />
                    </>
                ) : (
                    <div className="welcomeMessage">
                        <Link to="/account/login">Login</Link> to post tweets.
                    </div>
                )}
                {selectedTag && <TweetList taggedTweets={taggedTweets} selectedTag={selectedTag} handleTagClick={handleTagClick} />}
            </div>

            <div>
                <TagList tags={tags} handleTagClick={handleTagClick} />
            </div>
        </div>
    );
}

export default Main;

import React from 'react';
import './InputForm.css';

const InputForm = ({ user, formData, handleInputChange, handleAddTweet }) => {
    return (
        <form className="input-form">
            <div className="tweet-input-row">
                {/* Avatar */}
                <img src={user.picture} alt="User Avatar" className="avatar" />
                <textarea
                    className="tweet-input"
                    placeholder="What is happening?!"
                    name="content"
                    value={formData.content}
                    onChange={handleInputChange}
                ></textarea>
            </div>

            <div className="tweet-actions">
                {/* Post Button */}
                <button type="button" className="tweet-post-button" onClick={handleAddTweet}>
                    Post
                </button>
            </div>
        </form>
    );
};

export default InputForm;

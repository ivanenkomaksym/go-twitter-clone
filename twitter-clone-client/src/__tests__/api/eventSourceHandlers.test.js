jest.mock('../../common', () => require('../../__mocks__/configMock'));

import { setUpFeedsEventSource, setUpFeedsTagEventSource } from '../../api/eventSourceHandlers';
import { feedsUrl, getTaggedFeedsUrl } from '../../common';

describe('Event Source Handlers', () => {
    let mockSetDataCallback;
    let mockEventSource;

    beforeEach(() => {
        mockSetDataCallback = jest.fn();

        // Mock EventSource globally
        global.EventSource = jest.fn(() => {
            mockEventSource = {
                addEventListener: jest.fn(),
                removeEventListener: jest.fn(),
                close: jest.fn(),
            };
            return mockEventSource;
        });
    });

    afterEach(() => {
        jest.clearAllMocks();
        delete global.EventSource;
    });

    describe('setUpFeedsEventSource', () => {
        it('should initialize EventSource with the correct URL and credentials', () => {
            setUpFeedsEventSource(mockSetDataCallback);

            expect(EventSource).toHaveBeenCalledWith(
                feedsUrl,
                { withCredentials: true }
            );
        });

        it('should set up an event listener for "data" event and call setDataCallback with parsed data', () => {
            const expectedTags = [
                { name: 'tag1', nofTweets: 5 },
                { name: 'tag2', nofTweets: 10 },
            ];

            // Setup event source and listener
            setUpFeedsEventSource(mockSetDataCallback);

            // Simulate "data" event with expected JSON data
            const event = new Event('data');
            event.data = JSON.stringify({
                feeds: [
                    { name: 'tag1', tweets: 5 },
                    { name: 'tag2', tweets: 10 },
                ],
            });

            // Trigger the event
            mockEventSource.addEventListener.mock.calls[0][1](event);

            // Check if the callback was called with expectedTags
            expect(mockSetDataCallback).toHaveBeenCalledWith(expectedTags);
        });
    });

    describe('setUpFeedsTagEventSource', () => {
        const tag = 'sampleTag';

        it('should initialize EventSource with the correct URL and credentials for a specific tag', () => {
            setUpFeedsTagEventSource(tag, mockSetDataCallback);

            expect(EventSource).toHaveBeenCalledWith(
                getTaggedFeedsUrl(tag),
                { withCredentials: true }
            );
        });

        it('should set up an event listener for "data" event and call setDataCallback with parsed tweets data', () => {
            const data = ['tweet1', 'tweet2'];

            // Setup event source and listener
            setUpFeedsTagEventSource(tag, mockSetDataCallback);

            // Simulate "data" event with expected JSON data
            const event = new Event('data');
            event.data = JSON.stringify({
                tweets: data
            });

            // Trigger the event
            mockEventSource.addEventListener.mock.calls[0][1](event);

            // Check if the callback was called with expectedTags
            expect(mockSetDataCallback).toHaveBeenCalledWith(data);
        });
    });
});

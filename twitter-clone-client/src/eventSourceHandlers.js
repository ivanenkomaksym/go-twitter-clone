import clientConfig from './clientConfig';

// Function to set up EventSource for data changes
export const setUpFeedsEventSource = (setDataCallback) => {
    const eventSource = new EventSource(clientConfig.applicationUri + '/api/feeds');
  
    // Add event listener for the 'data' event
    eventSource.addEventListener('data', (event) => {
      try {
        const data = JSON.parse(event.data);
        const updatedTags = data.feeds.map(feed => feed.name);
        setDataCallback(updatedTags);
      } catch (error) {
        console.error('Error parsing event data:', error);
      }
    });
  
    // Return the EventSource instance to allow cleanup
    return eventSource;
  };

// Function to set up EventSource for data changes
export const setUpFeedsTagEventSource = (tag, setDataCallback) => {
    // Set up a new EventSource connection
    const eventSource = new EventSource(clientConfig.applicationUri + `/api/feeds/${tag}`);

    // Add event listener for the 'data' event
    eventSource.addEventListener('data', (event) => {
        try {
        const data = JSON.parse(event.data);
        setDataCallback(data.tweets || []);
      } catch (error) {
        console.error('Error parsing event data:', error);
      }
    });
  
    // Return the EventSource instance to allow cleanup
    return eventSource;
  };
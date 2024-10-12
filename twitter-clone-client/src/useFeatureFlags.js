import { useState, useEffect } from 'react';
import featureFlags from './config.js';

const useFeatureFlags = () => {
  const [features, setFeatures] = useState({});

  useEffect(() => {
    setFeatures(featureFlags);
  }, []);

  const isFeatureEnabled = (feature) => features[feature];

  return { isFeatureEnabled };
};

const IsAuthnEnabled = () => {
    const { isFeatureEnabled } = useFeatureFlags();
    return isFeatureEnabled('authn');
}

export default IsAuthnEnabled;
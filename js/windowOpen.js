import { prepareUrl } from './utils';

// Fix resource url in window.open method
const originalWindowOpen = window.open;

window.open = (url, windowName, windowFeatures) => {
  return originalWindowOpen(prepareUrl(url), windowName, windowFeatures);
};

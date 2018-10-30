import { prepareUrl } from './utils';

/**
 * History API CORS errors stubbing with window monkey patching
 */
if (window.history) {
  const { history } = window;
  const originalPushState = history.pushState.bind(history);
  const originalReplaceState = history.replaceState.bind(history);

  history.pushState = function(state, title, url) {
    originalPushState(state, title, prepareUrl(url))
  }
  history.replaceState = function(state, title, url) {
    originalReplaceState(state, title, prepareUrl(url))
  }
}

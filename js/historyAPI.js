import { prepareUrl } from './utils';

/**
 * History API CORS errors stubbing with window monkey patching
 */
if (window.history) {
  const { history } = window;
  const originalPushState = history.pushState.bind(history);
  const originalReplaceState = history.replaceState.bind(history);
  const changeStateFunc = (originalFunction, state, title, url) => {
    if (!url.includes('/http')) {
      originalFunction(state, title, prepareUrl(url));
    }
  };

  history.pushState = function(state, title, url) {
    changeStateFunc(originalPushState, state, title, url);
  }

  history.replaceState = function(state, title, url) {
    changeStateFunc(originalReplaceState, state, title, url);
  }
}

// Javascript automatically injected in every page as inline script

import './elementMutations';
import './historyAPI';
import './requests';
import './windowOpen';
import './postMessage';

import window from './window';
import location from './location';

window.location = location;

const saveVarsObject = new Proxy({}, {
  has() { return true; },
  get(target, property) { return property === 'window' ? window : window[property]; },
  set(target, property, value) {
    window[property] = value;
    return true;
  }
});

const loadScript = url =>
  fetch(url, {}, true)
    .then(response => response.text())
    .catch(error => {
      console.error(`Error loading script ${url}: ${error}`);
    });

const scriptPromises = [...targetScriptUrls, ...targetAsyncScriptUrls].map(url => loadScript(url));

Promise.all(scriptPromises).then(scripts => {
  scripts.forEach(script => {
    try {
      // Attempt to save global script variables in the window object
      const dynamicFunction = new Function('window', 'saveVarsObject', `with(saveVarsObject){${script}}`);
      dynamicFunction.call(window, window, saveVarsObject);
    } catch (error) {
      console.error('Error execute script:', error);
      const dynamicFunction = new Function('window', script);
      dynamicFunction.call(window, window);
    }
  })
});

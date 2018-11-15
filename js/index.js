// Javascript automatically injected in every page as inline script

import './elementMutations';
import './historyAPI';
import './requests';
import './windowOpen';
import './postMessage';

import window from './window';
import location from './location';

window.location = location;

const loadScript = url =>
  fetch(url)
    .then(response => response.text())
    .catch(error => {
      console.error(`Error loading script ${url}: ${error}`);
    });

const scriptPromises = window.targetScriptUrls.map(url => loadScript(url));

Promise.all(scriptPromises).then(scripts => {
  scripts.forEach(script => {
    const dinamicFunction = new Function('window', script);
    dinamicFunction.call(window, window);
  })
});

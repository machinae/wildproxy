// Javascript automatically injected in every page as inline script

import './elementMutations';
import './historyAPI';
import './requests';
import './windowOpen';
import './postMessage';

import window, { originalWindow } from './window';
import location from './location';

window.location = location;

function loadScript(url, window) {
  return fetch(url)
    .then(response => response.text())
    .then(response => eval(response))
    .catch(error => {
      console.warn(`Error loading script ${url}: ${error}`);
    });
};

window.targetScriptUrls.forEach(url => loadScript(url, window));

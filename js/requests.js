import { prepareUrl } from './utils';

// Proxying XHR requests
const originalOpenXHR = window.XMLHttpRequest.prototype.open;

window.XMLHttpRequest.prototype.open = function() {
  const args = [].slice.call(arguments);

  args[1] = prepareUrl(args[1]);

  return originalOpenXHR.apply(this, args);
};

// Proxying fetch requests
const originalFetch = window.fetch;

window.fetch = function(request, init = {}, dontPrepare = false) {
  const options = {
    method: init.method,
    headers: init.headers,
    body: init.body,
  };
  const url = dontPrepare ? request : prepareUrl(request.url || request);

  return originalFetch.call(this, url, options)
}

// Monkey patch jQuery.ajax if it exists
if (window.jQuery) {
  window.jQuery.ajaxPrefilter(function(options) {
    options.url = prepareUrl(options.url);
    if (options.crossDomain) {
      options.crossDomain = false;
    }
  });
}

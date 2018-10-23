// Javascript automatically injected in every page as inline script

import parseUrl from 'url-parse';

(function() {
  // Monkey patches XHR to proxy URLs
  // Source: https://github.com/Rob--W/cors-anywhere
  // Source: https://stackoverflow.com/questions/5202296/add-a-hook-to-all-ajax-requests-on-a-page
  var origin = window.location.origin || (window.location.protocol + '//' + window.location.host + (window.location.port ? ':' + window.location.port : ''));
  var targetURL = /^\/(https?:\/\/)?(w{3})?[a-z-\.]+/.exec(window.location.pathname)[0].replace(/^\//, '')

  var open = XMLHttpRequest.prototype.open;

  XMLHttpRequest.prototype.open = function() {
    var args = [].slice.call(arguments);

    args[1] = prepareUrl(args[1]);

    return open.apply(this, args);
  };

  // Proxying fetch requests
  const originalFetch = window.fetch;

  window.fetch = function(request, init = {}) {
    const options = {
      method: init.method,
      headers: init.headers,
      body: init.body,
    };
    const url = prepareUrl(request.url || request);

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

  function prepareUrl(url) {
    const parsedUrl = parseUrl(url);

    if (/^http/.test(url)) {
      return parsedUrl.origin === origin ? url : `${origin}/${url}`;
    } else {
      const withoutLeadingSlash = url[0] !== '/';
      let result = origin;

      result += withoutLeadingSlash ? parsedUrl.pathname : `/${targetURL}${url}`;

      return result;
    }
  }

  /**
   * Wraps function in try/catch for bypassing errors and application crashes
   * @param {Function} func Function for wrap
   * @returns {Function} Wrapper function
   */
  const silentWrapper = (func) => (...args) => {
    try {
      return func(...args)
    } catch (err) {
      console.warn(err)
    }
  }

  /**
   * History API CORS errors stubbing with window monkey patching
   */
  if (window.history) {
    window.history.pushState = silentWrapper(window.history.pushState)
    window.history.replaceState = silentWrapper(window.history.replaceState)
  }
})();

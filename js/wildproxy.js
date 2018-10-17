// Javascript automatically injected in every page as inline script

(function() {
  // Monkey patches XHR to proxy URLs
  // Source: https://github.com/Rob--W/cors-anywhere
  // Source: https://stackoverflow.com/questions/5202296/add-a-hook-to-all-ajax-requests-on-a-page
  // TODO support fetch API
  const origin = window.location.origin || (window.location.protocol + '//' + window.location.host + (window.location.port ? ':' + window.location.port : ''));
  const open = XMLHttpRequest.prototype.open;

  XMLHttpRequest.prototype.open = () => {
    const args = [].slice.call(arguments);
    args[1] = prependOrigin(args[1]);

    return open.apply(this, args);
  };

  // Monkey patch jQuery.ajax if it exists
  if (window.jQuery) {
    window.jQuery.ajaxPrefilter((options) => {
      options.url = prependOrigin(options.url)
      if (options.crossDomain) {
        options.crossDomain = false;
      }
    });
  }

  // prepend origin(proxy url) to the given URL if is a cross-domain URL
  const prependOrigin = (reqUrl) => {
    const targetOrigin = /^https?:\/\/([^\/]+)/i.exec(reqUrl);
    if (targetOrigin && targetOrigin.length && targetOrigin[0].toLowerCase() !== origin) {
      reqUrl = origin + '/' + reqUrl;
    }

    return reqUrl;
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
      console.error(err)
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

//Javascript automatically injected in every page

(function() {
  // Monkey patches XHR to proxy URLs
  // Source: https://github.com/Rob--W/cors-anywhere
  // Source: https://stackoverflow.com/questions/5202296/add-a-hook-to-all-ajax-requests-on-a-page
  // TODO look into relative urls in scripts like rel2abs
  // TODO support fetch API
  var origin = window.location.origin || (window.location.protocol + '//' + window.location.host + (window.location.port ? ':' + window.location.port : ''));

  var open = XMLHttpRequest.prototype.open;
  XMLHttpRequest.prototype.open = function() {
    var args = [].slice.call(arguments);
    args[1] = prependOrigin(args[1]);
    return open.apply(this, args);
  };

  // Monkey patch jQuery.ajax if it exists
  if (window.jQuery) {
    window.jQuery.ajaxPrefilter(function(options) {
      options.url = prependOrigin(options.url)
      if (options.crossDomain) {
        options.crossDomain = false;
      }
    });
  }

  // prepend origin(proxy url) to the given URL if is a cross-domain URL
  var prependOrigin = function(reqUrl) {
    var targetOrigin = /^https?:\/\/([^\/]+)/i.exec(reqUrl);
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
  var silentWrap = function(func) {
    return function () {
      try {
        func.apply(this, arguments)
      } catch (err) {
        console.error(err)
      }
    }
  }

  window.history.pushState = silentWrap(window.history.pushState)
  window.history.replaceState = silentWrap(window.history.replaceState)
})();

"use strict";

//Javascript automatically injected in every page
(function () {
  var _this = this,
      _arguments = arguments;

  // Monkey patches XHR to proxy URLs
  // Source: https://github.com/Rob--W/cors-anywhere
  // Source: https://stackoverflow.com/questions/5202296/add-a-hook-to-all-ajax-requests-on-a-page
  // TODO look into relative urls in scripts like rel2abs
  // TODO support fetch API
  var origin = window.location.origin || window.location.protocol + '//' + window.location.host + (window.location.port ? ':' + window.location.port : '');
  var open = XMLHttpRequest.prototype.open;

  XMLHttpRequest.prototype.open = function () {
    var args = [].slice.call(_arguments);
    args[1] = prependOrigin(args[1]);
    return open.apply(_this, args);
  }; // Monkey patch jQuery.ajax if it exists


  if (window.jQuery) {
    window.jQuery.ajaxPrefilter(function (options) {
      options.url = prependOrigin(options.url);

      if (options.crossDomain) {
        options.crossDomain = false;
      }
    });
  } // prepend origin(proxy url) to the given URL if is a cross-domain URL


  var prependOrigin = function prependOrigin(reqUrl) {
    var targetOrigin = /^https?:\/\/([^\/]+)/i.exec(reqUrl);

    if (targetOrigin && targetOrigin.length && targetOrigin[0].toLowerCase() !== origin) {
      reqUrl = origin + '/' + reqUrl;
    }

    return reqUrl;
  };
  /**
   * Wraps function in try/catch for bypassing errors and application crashes
   * @param {Function} func Function for wrap
   * @returns {Function} Wrapper function
   */


  var silentWrapper = function silentWrapper(func) {
    return function () {
      try {
        return func.apply(void 0, arguments);
      } catch (err) {
        console.error(err);
      }
    };
  };
  /**
   * History API CORS errors stubbing with window monkey patching
   */


  if (window.history) {
    window.history.pushState = silentWrapper(window.history.pushState);
    window.history.replaceState = silentWrapper(window.history.replaceState);
  }
})();
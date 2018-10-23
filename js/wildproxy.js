// Javascript automatically injected in every page as inline script

(function() {
  // Monkey patches XHR to proxy URLs
  // Source: https://github.com/Rob--W/cors-anywhere
  // Source: https://stackoverflow.com/questions/5202296/add-a-hook-to-all-ajax-requests-on-a-page
  // TODO support fetch API
  var origin = window.location.origin || (window.location.protocol + '//' + window.location.host + (window.location.port ? ':' + window.location.port : ''));
  var targetURL = /^\/(https?:\/\/)?(w{3})?[a-z-\.]+/.exec(window.location.pathname)[0].replace(/^\//, '')

  var open = XMLHttpRequest.prototype.open;

  XMLHttpRequest.prototype.open = function() {
    var args = [].slice.call(arguments);


    args[1] = prependOrigin(args[1]);
    args[1] = normalizeRelativePath(args[1])

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

  const normalizeRelativePath = (path) => {
    if (!/^https?/.test(path.replace(`${origin}/`, ''))) {
      return path.replace(origin, '')
    }

    return path
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

  const replace = (func) => (...args) => {
    try {
      return func(...args)
    } catch (err) {
      console.error(err)
    }
  }

  /**
   * history API CORS errors stubbing with window monkey patching
   */
  if (window.history) {
    const pushState = window.history.pushState.bind(window.history),
    replaceState = window.history.replaceState.bind(window.history);
    history.pushState = function(state, title, url) {
      const prependUrl = origin + '/' + targetURL + url;
      pushState(state, title, prependUrl)
    }
    history.replaceState = function(state, title, url) {
      const prependUrl = origin + '/' + targetURL + url;
      replaceState(state, title, prependUrl)
    }
  }
})();

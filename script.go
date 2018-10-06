package main

// Javascript to inject into the page
// Minified contents of injectScript.js manually added here
var injectScript = `!function(){var n=window.location.origin||window.location.protocol+"//"+window.location.host+(window.location.port?":"+window.location.port:""),t=XMLHttpRequest.prototype.open;XMLHttpRequest.prototype.open=function(){var o=[].slice.call(arguments);return o[1]=i(o[1]),t.apply(this,o)},window.jQuery&&window.jQuery.ajaxPrefilter(function(o){o.url=i(o.url),o.crossDomain&&(o.crossDomain=!1)});var i=function(o){var t=/^https?:\/\/([^\/]+)/i.exec(o);return t&&t.length&&t[0].toLowerCase()!==n&&(o=n+"/"+o),o}}();`

package main

// Javascript to inject into the page
// Minified contents of injectScript.js manually added here
var injectScript = `!function(){var e=window.location.protocol+"//"+window.location.host,p=XMLHttpRequest.prototype.open;XMLHttpRequest.prototype.open=function(){var t=[].slice.call(arguments),o=/^https?:\/\/([^\/]+)/i.exec(t[1]);return o&&o[0].toLowerCase()!==e&&(t[1]=e+"/"+t[1]),p.apply(this,t)}}();`

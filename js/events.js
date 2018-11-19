const { document } = window;
const originalAddEventListener = window.addEventListener;
const addEventListener = (type, listener, options) => {
  if (
    (document.readyState !== 'loading' && type === 'DOMContentLoaded') ||
    (document.readyState === 'complete' && type === 'load')
  ) {
    listener();
  } else {
    originalAddEventListener(type, listener, options);
  }
};

window.addEventListener = addEventListener;
document.addEventListener = addEventListener;

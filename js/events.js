const { document } = window;
const originalAddEventListener = window.addEventListener;
const originalRemoveEventListener = window.removeEventListener;

const addEventListener = (type, listener, options) => {
  if (
    (document.readyState !== 'loading' && type === 'DOMContentLoaded') ||
    (document.readyState === 'complete' && type === 'load')
  ) {
    listener();
  } else {
    originalAddEventListener.call(window, type, listener, options);
  }
};

const removeEventListener = () => {
  originalRemoveEventListener.apply(window, arguments);
}

window.addEventListener = document.addEventListener = addEventListener;
window.removeEventListener = removeEventListener;


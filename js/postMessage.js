import { TARGET_URL } from './constants';

const originalWindowPostMessage = window.postMessage;
const getData = message => ({
  message,
  origin: TARGET_URL,
});

window.postMessage = (message, targetOrigin, transfer) =>
  originalWindowPostMessage.call(window, getData(message), '*', transfer);

if (window.parent) {
  window.parent.postMessage = (message, targetOrigin, transfer) =>
    originalWindowPostMessage.call(window.parent, getData(message), '*', transfer);
}

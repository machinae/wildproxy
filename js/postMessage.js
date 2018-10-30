// Fix targetOrigin in window.postMessage method
const originalWindowPostMessage = window.postMessage;

window.postMessage = (message, targetOrigin, transfer) =>
  originalWindowPostMessage(message, '*', transfer);

if (window.parent) {
  window.parent.postMessage = (message, targetOrigin, transfer) =>
    originalWindowPostMessage(message, '*', transfer);
}

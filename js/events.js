const { document } = window;
const originalAddEventListener = window.addEventListener;
const originalRemoveEventListener = window.removeEventListener;

const handleMessageEvent = listener => event => {
  try {
    const { message, origin } = JSON.parse(event.data);
    const newEvent = new MessageEvent(event.type, {
      data: message,
      origin: origin,
      lastEventId: event.lastEventId,
      source: event.source,
      ports: event.ports,
    });

    listener(newEvent);
  } catch(error) {
    listener(event);
  }
};

const addEventListener = (type, listener, options) => {
  if (
    (document.readyState === 'complete' && type === 'load')
  ) {
    listener();
  } else if (type === 'message') {
    originalAddEventListener.call(window, type, handleMessageEvent(listener), options);
  } else {
    originalAddEventListener.call(window, type, listener, options);
  }
};

const removeEventListener = () => {
  originalRemoveEventListener.apply(window, arguments);
}

window.addEventListener = document.addEventListener = addEventListener;
window.removeEventListener = removeEventListener;

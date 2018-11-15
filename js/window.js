import { proxyObject } from './utils';

export const originalWindow = window;

export default new Proxy(proxyObject(originalWindow), {
  set(target, property, value) {
    target[property] = value;

    if (!value || !value.isProxyObject) {
      originalWindow[property] = value;
    }

    return true;
  },
});

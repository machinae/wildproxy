import parseUrl from 'url-parse';

import { TARGET_URL, ORIGIN_URL } from './constants';

export const prepareUrl = url => {
  if (!url) {
    return url;
  }

  const parsedUrl = parseUrl(url);

  if (/^http/.test(url)) {
    return parsedUrl.origin === origin ? url : `${origin}/${url}`;
  } else {
    const withoutLeadingSlash = url[0] !== '/';
    let result = ORIGIN_URL;

    result += withoutLeadingSlash ? parsedUrl.pathname : `/${TARGET_URL}${url}`;

    return result;
  }
};

// The standard proxy is not suitable because it has problems with proxying non-configurable properties
export const proxyObject = (originalObject, originalContext) => {
  const targetObjet = Object.create({});

  if (!originalContext) {
    originalContext = originalObject;
  }

  Object.getOwnPropertyNames(originalObject).forEach(propertyKey => {
    const prop = originalObject[propertyKey];

    if (prop instanceof Function) {
      const hasPrototype = !!originalObject[propertyKey].prototype;

      targetObjet.__proto__[propertyKey] = hasPrototype ? prop : prop.bind(originalContext);
    } else if (prop instanceof Object) {
      targetObjet.__proto__[propertyKey] = prop;
    } else {
      Object.defineProperty(targetObjet.__proto__, propertyKey, {
        get() {
          return originalObject[propertyKey];
        },
        set(value) {
          const descriptor = Object.getOwnPropertyDescriptor(originalObject, propertyKey);

          if (descriptor.set) {
            originalObject[propertyKey] = value;
          }
        }
      });
    }
  });

  if (originalObject.__proto__ instanceof Object) {
    targetObjet.__proto__.__proto__ = proxyObject(originalObject.__proto__, originalContext);
  }

  targetObjet.isProxyObject = true;

  return targetObjet;
}

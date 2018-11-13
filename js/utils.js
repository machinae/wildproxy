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
export const proxyObject = originalObject => {
  const targetObjet = Object.create({});

  for (let propertyKey in originalObject) {
    const prop = originalObject[propertyKey];

    if (prop instanceof Function) {
      targetObjet.__proto__[propertyKey] = prop.bind(originalObject);
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
  }

  return targetObjet;
}

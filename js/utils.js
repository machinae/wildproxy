import parseUrl from 'url-parse';

import { TARGET_URL, ORIGIN_URL } from './constants';

const parsedTargetUrl = parseUrl(TARGET_URL);

export const prepareUrl = url => {
  const parsedUrl = parseUrl(url);

  if (/^http/.test(url)) {
    if (parsedUrl.origin !== ORIGIN_URL && parsedUrl.hostname.includes(parsedTargetUrl.hostname)) {
      return `${ORIGIN_URL}/${url}`;
    } else {
      return url;
    }
  } else {
    const withoutLeadingSlash = url[0] !== '/';
    let result = ORIGIN_URL;

    result += withoutLeadingSlash ? parsedUrl.pathname : `/${TARGET_URL}${url}`;

    return result;
  }
};

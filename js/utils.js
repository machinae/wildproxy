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

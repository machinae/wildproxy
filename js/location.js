import parseUrl from 'url-parse';

import { originalWindow } from './window';
import { prepareUrl, proxyObject } from './utils';

const { location } = originalWindow;
const proxiedLocation = proxyObject(location);

const getParsedOriginalPathname = () => parseUrl(location.pathname.substring(1));

const updateUrlField = (field, value) => {
  const changed = getParsedOriginalPathname().set(field, value);

  originalWindow.location = prepareUrl(changed.toString());
}

Object.defineProperties(proxiedLocation, {
  hostname: {
    get() {
      return getParsedOriginalPathname().hostname;
    },
    set(value) {
      updateUrlField('hostname', value);
    },
  },
  host: {
    get() {
      return getParsedOriginalPathname().host;
    },
    set(value) {
      updateUrlField('host', value);
    },
  },
  href: {
    get() {
      return getParsedOriginalPathname().href;
    },
    set(value) {
     originalWindow.location = prepareUrl(value);
    },
  },
  origin: {
    get() {
      return getParsedOriginalPathname().origin;
    },
  },
  pathname: {
    get() {
      return getParsedOriginalPathname().pathname;
    },
    set(value) {
      updateUrlField('pathname', value);
    },
  },
  port: {
    get() {
      return getParsedOriginalPathname().port;
    },
    set(value) {
      updateUrlField('port', value);
    },
  },
  protocol: {
    get() {
      return getParsedOriginalPathname().protocol;
    },
    set(value) {
      updateUrlField('protocol', value);
    },
  },
  assign: {
    value: url => location.assign(prepareUrl(url)),
  },
  replace: {
    value: url => location.replace(prepareUrl(url)),
  },
});

export default proxiedLocation;

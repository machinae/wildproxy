const {
  location: {
    host,
    origin,
    pathname,
    port,
    protocol,
  }
} = window;

export const ORIGIN_URL = origin || (protocol + '//' + host + (port ? ':' + port : ''));
export const TARGET_URL = /^\/(https?:\/\/)?(w{3})?[a-z-\.]+/.exec(pathname)[0].replace(/^\//, '');
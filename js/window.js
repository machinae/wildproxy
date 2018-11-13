import { proxyObject } from './utils';

export const originalWindow = window;

export default proxyObject(originalWindow);

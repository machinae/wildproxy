try {
    window.parent.postMessage({
        type: 'proxy.beforeStart',
    }, '*');
} catch (e) {}

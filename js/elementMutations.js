import { prepareUrl } from './utils';

function updateNodeSourceAttribute(node, attributeName) {
  const oldUrl = node.getAttribute(attributeName);
  const newUrl = prepareUrl(oldUrl);

  if (newUrl !== oldUrl && node.parentNode) {
    // Fix this.removeEventListener in load event
    const clonedNode = node.cloneNode(true);

    clonedNode.setAttribute(attributeName, newUrl);
    node.parentNode.replaceChild(clonedNode, node);
  }
}

window.addEventListener('load', () => {
  const attributeFilter = ['src', 'href'];
  const sourceSelector = attributeFilter.reduce((accumulator, currentValue) => `[${accumulator}], [${currentValue}]`);
  const observer = new MutationObserver(mutations => {
    mutations.forEach(({ addedNodes, attributeName, target, type}) => {
      if (type === 'childList') {
        addedNodes.forEach(node => {
          const children = node.querySelectorAll ? node.querySelectorAll(sourceSelector) : [];

          [node, ...children].forEach(n => {
            const attr = attributeFilter.find(attribute => attribute in n);

            if (attr) {
              updateNodeSourceAttribute(n, attr);
            }
          });
        });
      } else if (type === 'attributes') {
        updateNodeSourceAttribute(target, attributeName);
      }
    });
  });

  observer.observe(document.body, {
    attributeFilter,
    attributes: true,
    childList: true,
    subtree: true
  });
});

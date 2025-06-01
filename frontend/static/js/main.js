import { initEventListeners } from './handlers.js';
import { dom } from './dom.js';

document.addEventListener("DOMContentLoaded", () => {
  dom.visualizationSection.style.display = "none";
  initEventListeners();
});
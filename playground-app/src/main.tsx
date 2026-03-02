import React from "react";
import ReactDOM from "react-dom/client";
import App from "./App.tsx";

ReactDOM.createRoot(document.getElementById("playground")!).render(
  <React.StrictMode>
    <App />
  </React.StrictMode>
);

// MkDocs "navigation.instant" swaps DOM without re-executing <script type="module">.
// Monaco's CSS is lost when this happens, so force a full reload on re-navigation.
let navigatedAway = false;
new MutationObserver(() => {
  if (!document.getElementById("playground")) {
    navigatedAway = true;
  } else if (navigatedAway) {
    window.location.reload();
  }
}).observe(document.body, { childList: true, subtree: true });

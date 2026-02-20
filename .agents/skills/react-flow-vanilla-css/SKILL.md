---
name: React Flow Vanilla CSS Styling
description: Best practices for styling React Flow with native CSS and avoiding canvas collapsing or rendering errors.
---
# React Flow Vanilla CSS Styling

## Key Findings & Caveats
1. **Container Height is Critical:** The React Flow canvas will disappear (compile but render completely blank) if the parent `.workflow-container` or the `.workflow-canvas` itself does not have an explicit `height` or `min-height` defined.
2. **Tailwind Class Conflicts:** If using custom Tailwind configurations or if Tailwind classes are not being aggressively applied/compiled to a specific component, it's safer to use native CSS classes (`.node-card`, `.workflow-container`) specifically on customized Flow components.
3. **Debugging Blank Canvas:** If the canvas is blank:
   - Wrap the `<ReactFlow>` component in a `try/catch` ErrorBoundary to catch silent JavaScript errors caused by bad custom node configurations or missing imports.
   - Add a loud structural border (e.g., `border: 2px dashed red`) to the wrapper to visually confirm whether it's an invisible canvas container or missing node elements.
   - Do a hard browser refresh if using a Docker container, as local browser caches might not pick up new `.css` outputs.

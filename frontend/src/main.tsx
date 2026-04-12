// 1. We import the core React library. 
import React from 'react'

// 2. We import the ReactDOM library. This is the specific tool that takes React 
// code and translates it into actual DOM nodes the browser can render.
import ReactDOM from 'react-dom/client'

// 3. We import our main App component (which we will write next).
import App from './App.tsx'

// 4. We find the empty <div id="root"></div> in our index.html file.
const rootElement = document.getElementById('root')!

// 5. We tell ReactDOM to take control of that element.
const root = ReactDOM.createRoot(rootElement)

// 6. We render our App inside of it. 
// StrictMode is a development tool that renders everything twice to catch bugs.
root.render(
  <React.StrictMode>
    <App />
  </React.StrictMode>,
)
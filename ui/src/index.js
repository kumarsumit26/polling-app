import React from 'react';
import ReactDOM from 'react-dom/client';
import './index.css';
import App from './App';
import axios from 'axios'; // Import axios

// Set the Axios default base URL from the environment variable
axios.defaults.baseURL = process.env.REACT_APP_BACKEND_URL || 'http://localhost:8080';

const root = ReactDOM.createRoot(document.getElementById('root'));
root.render(
  <React.StrictMode>
    <App />
  </React.StrictMode>
);

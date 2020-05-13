import React from 'react';
import ReactDOM from 'react-dom';
import './index.css';
import HeaderFloating from './home/header';
import App from './home/home';
import Logging from './usercenter/login';
import * as serviceWorker from './serviceWorker';


ReactDOM.render(<HeaderFloating />, document.getElementById('HeaderFloating'));

if (window.location.pathname === "/") {
  ReactDOM.render(<App />, document.getElementById('root'));
}

if (window.location.pathname === "/login") {
  ReactDOM.render(<Logging />, document.getElementById('Logging'));
}





// If you want your app to work offline and load faster, you can change
// unregister() to register() below. Note this comes with some pitfalls.
// Learn more about service workers: https://bit.ly/CRA-PWA
serviceWorker.unregister();

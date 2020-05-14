import React from 'react';
import ReactDOM from 'react-dom';
import './index.css';
import HeaderFloating from './home/header';
import App from './home/home';
import Logging from './usercenter/login';
import * as serviceWorker from './serviceWorker';
import Cookies from 'universal-cookie';

const cookies = new Cookies();


ReactDOM.render(<HeaderFloating />, document.getElementById('HeaderFloating'));

if (cookies.get("username") && cookies.get("usertoken")) {
  ReactDOM.render(<App />, document.getElementById('root'));
}

if (!cookies.get("username") || !cookies.get("usertoken")) {
  ReactDOM.render(<Logging />, document.getElementById('Logging'));
}





// If you want your app to work offline and load faster, you can change
// unregister() to register() below. Note this comes with some pitfalls.
// Learn more about service workers: https://bit.ly/CRA-PWA
serviceWorker.unregister();

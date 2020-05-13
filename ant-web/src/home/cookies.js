import Cookies from 'universal-cookie';

const cookies = new Cookies();

// 清理token
export function CleanToken() {
  cookies.remove('username');
  cookies.remove('usertoken');
  window.location.href="/login";
}

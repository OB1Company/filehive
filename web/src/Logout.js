import React from 'react';
import { useCookies } from 'react-cookie';

export default function Logout() {
  const [cookies, setCookie, removeCookie] = useCookies(['token']);
  removeCookie("token");

  return null;
}


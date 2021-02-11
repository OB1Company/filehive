import React from 'react';
import { useCookies } from 'react-cookie';
import {getAxiosInstance} from "./components/Auth";

export default function Logout() {
  const [cookies, setCookie, removeCookie] = useCookies(['token']);
  removeCookie("token");

  const instance = getAxiosInstance();
  instance.post('/api/v1/logout')
      .then((result) => {
          localStorage.removeItem("name");
          localStorage.removeItem("email");
          localStorage.removeItem("admin");
      })

  return null;
}


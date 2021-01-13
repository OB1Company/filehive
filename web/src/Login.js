import React, {useState} from 'react'
import axios from "axios";
import { SET_TOKEN } from "./components/Store";
import { useHistory } from "react-router-dom";
import {Link} from "react-router-dom";
import ErrorBox from './components/ErrorBox'
import { config } from 'dotenv'

function Login() {

  const history = useHistory();

  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [isError, setIsError] = useState(false);
  const [error, setError] = useState(false)

  const HandleFormSubmit = async (e) => {
    const env = config();

    e.preventDefault();


    try {
      const data = { email, password };

      const getCsrfToken = async () => {
        try {

          const csrftoken = localStorage.getItem('token');
          const instance = axios.create({
            baseURL: "",
            headers: { "x-csrf-token": csrftoken }
          })

          const loginUrl = "/api/v1/login";
          const apiReq = await instance.post(
              loginUrl,
              data
          );

          // Successful login
          console.log(apiReq);
          history.push("/");

        } catch(err) {
          console.log(err);
        }
      };
      await getCsrfToken();


      //history.push("/");
    } catch (error) {

      // Check for csrf issue
      console.log(error.response);
      if(error.response.data === "Forbidden - CSRF token invalid\n") {
        console.log("NO CSRF");

      }

      setIsError(true);
      setError(error.response.data.message);
    }
    return false;
  }

  return (
    <div class="Login form-540">
      <h2>Log in</h2>
      <form onSubmit={HandleFormSubmit}>
        <label>
          Email address
          <input type="text" name="email" placeholder="Enter email" onChange={e => setEmail(e.target.value)}/>
        </label>
        <label>
          Password
          <input type="password" name="password" placeholder="Password"  onChange={e => setPassword(e.target.value )} />
        </label>
        <div>
          <input type="submit" value="Log in" class="orange-button" />
          <Link to='/passwordreset'>Forgot password?</Link>
        </div>
        
        {error &&
          <ErrorBox message="Incorrect email/password. Try again."/>
        }
      </form>
    </div>
  )
}

export default Login
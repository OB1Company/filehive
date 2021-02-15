import React, {useEffect, useState} from 'react'
import axios from "axios";
import {useHistory, useLocation} from "react-router-dom";
import {Link} from "react-router-dom";
import ErrorBox, {SuccessBox} from "./components/ErrorBox";
import {Helmet} from "react-helmet";
import spinner from "./images/spinner.gif";

function Login() {

  const history = useHistory();
  const location = new URLSearchParams(useLocation().search)

  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [isError, setIsError] = useState(false);
  const [error, setError] = useState("");
  const [confirmationMessage, setConfirmationMessage] = useState("");
  const [isLoggingIn, setIsLoggingIn] = useState(false);

  if(location.get("confirmed") === "1") {
    setConfirmationMessage("Your account has been confirmed.");
  }

  const HandleFormSubmit = async (e) => {
    e.preventDefault();

    setError("");
    setIsLoggingIn(true);

    const data = { email, password };

    if(data.email === "") {
      setError("Enter an email address");
      setIsLoggingIn(false);
      return false;
    }

    if(data.password === "") {
      setError("Enter a password");
      setIsLoggingIn(false);
      return false;
    }

    const login = async () => {

      const csrftoken = localStorage.getItem('csrf_token');
      const instance = axios.create({
        baseURL: "",
        headers: { "x-csrf-token": csrftoken }
      })
      const loginUrl = "/api/v1/login";

      try {
        await instance.post(
            loginUrl,
            data
        ).then((data)=>{
          localStorage.setItem("email", email);

          instance.get("/api/v1/user/" + email)
              .then((data) => {

                if(!data.data.disabled) {
                  localStorage.setItem("name", data.data.name);
                  localStorage.setItem("userID", data.data.UserID);
                  localStorage.setItem("admin", data.data.admin);
                  history.push("/dashboard");
                } else {
                  localStorage.removeItem("name");
                  localStorage.removeItem("email");
                  localStorage.removeItem("admin");
                  localStorage.removeItem("userID");
                  setIsLoggingIn(false);
                  setError("This account has been disabled");
                  history.push("/login");
                }

              })

        }).catch(error => {
          console.log("Login Failure", error.response);
          setIsError(true);
          if(error.response.status === 403) {
            localStorage.removeItem('csrf_token');
            window.location.reload();
          } else {
            const errorMessage = error.response.data.error;
            setError(errorMessage);
          }
          setIsLoggingIn(false);
        });

        return false;

      } catch(err) {
        if(err.response.data === "Forbidden - CSRF token invalid\n") {
          localStorage.removeItem('csrf_token');
          history.push('/login');
        }
      }
    };

    try {
      await login();
    } catch (error) {
      console.log(error);
    }
    return false;
  }

  const LoginButton = () => {
    if (!isLoggingIn) {
      return  <input type="submit" value="Login" className="raise orange-button" />
    } else {
      return <span className="spinner-span">
        <img src={spinner} width="20" height="20" alt="spinner" className="noblock"/> Logging in...
        </span>
    }
  }

  return (

    <div className="Login form-540">

      <Helmet>
        <title>Filehive | Login</title>
      </Helmet>

      {confirmationMessage &&
          <div className="marginbottom-15">
            <SuccessBox message={confirmationMessage}/>
          </div>
      }

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
          <LoginButton/>
          <Link to='/password_reset'>Forgot password?</Link>
        </div>
        
        {error &&
          <ErrorBox message={error}/>
        }
      </form>
    </div>
  )
}

export default Login
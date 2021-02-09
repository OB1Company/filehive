import React, {useEffect, useState} from 'react'
import axios from "axios";
import {useHistory, useLocation} from "react-router-dom";
import {Link} from "react-router-dom";
import ErrorBox, {SuccessBox} from "./components/ErrorBox";
import {Helmet} from "react-helmet";

function Login() {

  const history = useHistory();
  const location = new URLSearchParams(useLocation().search)

  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [isError, setIsError] = useState(false);
  const [error, setError] = useState("");
  const [confirmationMessage, setConfirmationMessage] = useState("");

  useEffect(()=>{
    if(location.get("confirmed") === "1") {
      setConfirmationMessage("Your account has been confirmed.");
    }
  })

  const HandleFormSubmit = async (e) => {
    e.preventDefault();

    const data = { email, password };

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
                localStorage.setItem("name", data.data.Name);
                history.push("/dashboard");
              })


        }).catch(error => {
          console.log("Login Failure", error.response);
          setIsError(true);
          if(error.response.status === 403) {
            window.location.reload(false);
          } else {
            const errorMessage = error.response.data.error;
            setError(errorMessage);
          }
          localStorage.removeItem('csrf_token');
          history.push('/login');
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
      // Check for csrf issue
      console.log(error);
      // if(error.response.data === "Forbidden - CSRF token invalid\n") {
      //   console.log("NO CSRF");
      //
      // }
      //
      // setIsError(true);
      // setError(error.response.data.message);
      // console.log(error.response);
    }
    return false;
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
          <input type="submit" value="Log in" className="raise orange-button" />
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
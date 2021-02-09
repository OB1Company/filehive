import React, {useState} from 'react'
import {Link} from "react-router-dom";
import ErrorBox from './components/ErrorBox'
import { useHistory } from "react-router-dom";
import Select from 'react-select'
import { Countries } from './constants/Countries'
import { getAxiosInstance } from "./components/Auth";
import { Helmet } from 'react-helmet';

function Signup() {

  const history = useHistory();

  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [name, setName] = useState("");
  const [country, setCountry] = useState("");
  const [error, setError] = useState("")

  const HandleFormSubmit = async (e) => {
    e.preventDefault();

    const data = { email, password, country, name };

    if(name === "") {
      setError("Name is required");
      return false;
    }
    if(country === "") {
      setError("Country is required")
      return false;
    }

    const instance = getAxiosInstance();

    const createUserUrl = "/api/v1/user";
    await instance.post(
        createUserUrl,
        data
    ).then((data) => {
      // Successful login
      localStorage.setItem("email", email);
      localStorage.setItem("name", name);

      history.push("/dashboard");
    }).catch((error) => {
      if(error.response.status === 409) {
        setError("This email address has already been used");
        return false;
      }
      if(error.response.status === 400) {
        setError("There was an error with your registration: "+error.response.data.error);
        return false;
      }
      return false;
    });

  }

  const handleCountry = (e) => {
    setCountry(e.value);
  }

  return (
    <div className="Signup form-540 form-center">
      <Helmet>
        <title>Filehive | Signup</title>
      </Helmet>

      <h2>Sign up</h2>
      <form onSubmit={HandleFormSubmit}>
        <label>
          Email address*
          <input type="text" name="email" placeholder="Enter email" onChange={e => setEmail(e.target.value)}/>
        </label>
        <label>
          Password*
          <input type="password" name="password" placeholder="Password" onChange={e => setPassword(e.target.value)}/>
        </label>
        <label>
          Name*
          <input type="text" name="name" placeholder="Your name (shown publicly)" onChange={e => setName(e.target.value)}/>
        </label>
        <label>
          Country*
            <Select name="country" options={Countries} placeholder="--" onChange={handleCountry}/>
        </label>
        <div>
          <input type="submit" value="Sign up" className="raise orange-button" />
          <Link to='/login'>Already registered?</Link>
        </div>

        {error &&
          <ErrorBox message={error}/>
        }
      </form>
    </div>
  )
}

export default Signup
import React, {useState} from 'react'
import {Link} from "react-router-dom";
import ErrorBox from './components/ErrorBox'
import { useHistory } from "react-router-dom";
import Select from 'react-select'
import { Countries } from './constants/Countries'
import axios from "axios";

function Signup() {

  const history = useHistory();

  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [name, setName] = useState("");
  const [country, setCountry] = useState("");
  const [error, setError] = useState(false)

  const HandleFormSubmit = async (e) => {
    e.preventDefault();

    const data = { email, password, country, name };

    const csrftoken = localStorage.getItem('csrf_token');
    const instance = axios.create({
      baseURL: "",
      headers: { "x-csrf-token": csrftoken }
    })

    const loginUrl = "/api/v1/user";
    const apiReq = await instance.post(
        loginUrl,
        data
    );

    // Successful login
    console.log(apiReq);
    localStorage.setItem("username", name);
    localStorage.setItem("email", email);

    history.push("/");

    return false;
  }

  const handleCountry = (e) => {
    setCountry(e.value);
  }

  return (
    <div class="Signup form-540">
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
          <input type="submit" value="Sign up" class="orange-button" />
          <Link to='/login'>Already registered?</Link>
        </div>

        {error &&
          <ErrorBox message="An account is already registered with that email address."/>
        }
      </form>
    </div>
  )
}

export default Signup
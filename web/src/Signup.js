import React, {useState} from 'react'
import {Link} from "react-router-dom";
import ErrorBox from './components/ErrorBox'
import Select from 'react-select'
import { Countries } from './constants/Countries'
import axios from "axios";

function Signup() {

  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [country, setCountry] = useState("");
  const [error, setError] = useState(false)

  const HandleFormSubmit = async (e) => {
    e.preventDefault();

    const data = { email, password, country };
    console.log(data);

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
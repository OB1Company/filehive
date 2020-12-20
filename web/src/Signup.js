import React from 'react'
import {Link} from "react-router-dom";
import ErrorBox from './components/ErrorBox'
import Select from 'react-select'


function Signup() {

    const options = [
        { value: 'chocolate', label: 'Chocolate' },
        { value: 'strawberry', label: 'Strawberry' },
        { value: 'vanilla', label: 'Vanilla' }
    ]

  return (
    <div class="Signup form-540">
      <h2>Sign up</h2>
      <form>
        <label>
          Email address*
          <input type="text" name="email" placeholder="Enter email" />
        </label>
        <label>
          Password*
          <input type="text" name="password" placeholder="Password" />
        </label>
        <label>
          Country*
            <Select options={options} placeholder="--"/>
        </label>
        <div>
          <input type="submit" value="Sign up" class="orange-button" />
          <Link to='/login'>Already registered?</Link>
        </div>
        <ErrorBox message="An account is already registered with that email address."/>
      </form>
    </div>
  )
}

export default Signup
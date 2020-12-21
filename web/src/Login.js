import React, {useState} from 'react'
import {Link} from "react-router-dom";
import ErrorBox from './components/ErrorBox'

function Login() {
  const [error, setError] = useState(false)

  return (
    <div class="Login form-540">
      <h2>Log in</h2>
      <form>
        <label>
          Email address
          <input type="text" name="email" placeholder="Enter email" />
        </label>
        <label>
          Password
          <input type="password" name="password" placeholder="Password" />
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
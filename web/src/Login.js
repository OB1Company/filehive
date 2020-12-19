import React from 'react'
import {Link} from "react-router-dom";

function Login() {
  return (
    <div class="Login">
      <h2>Log in</h2>
      <form>
        <label>
          Email address
          <input type="text" name="email" placeholder="Enter email" />
        </label>
        <label>
          Password
          <input type="text" name="password" placeholder="Password" />
        </label>
        <div>
          <input type="submit" value="Log in" class="orange-button" />
          <Link to='/passwordreset'>Forgot password?</Link>
        </div>

      </form>
    </div>
  )
}

export default Login
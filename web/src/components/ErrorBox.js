import React from 'react'

function ErrorBox(props) {
  return (
    <div class="ErrorBox">
      {props.message}
    </div>
  )
}

export default ErrorBox
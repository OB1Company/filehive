import React from 'react'

function ErrorBox(props) {
  return (
    <div className="ErrorBox">
      {props.message}
    </div>
  )
}

export function SuccessBox(props) {
  return (
      <div className="SuccessBox">
        {props.message}
      </div>
  )
}

export default ErrorBox